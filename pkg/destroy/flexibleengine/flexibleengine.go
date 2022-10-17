package flexibleengine

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"

	"github.com/openshift/installer/pkg/destroy/providers"
	"github.com/openshift/installer/pkg/types"

	hwsdk "github.com/chnsz/golangsdk"
	hwclient "github.com/chnsz/golangsdk/openstack"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/wait"
)

var (
	defaultTimeout = 2 * time.Minute
)

type ClusterUninstaller struct {
	Logger  logrus.FieldLogger
	Context context.Context

	AccessKey   string
	SecretKey   string
	DomainName  string
	Region      string
	ProjectName string
	Tags        map[string]string

	NetworkClientV1 *hwsdk.ServiceClient
	NetworkClientV2 *hwsdk.ServiceClient
}

func New(logger logrus.FieldLogger, metadata *types.ClusterMetadata) (providers.Destroyer, error) {
	return &ClusterUninstaller{
		Logger:  logger,
		Context: context.Background(),

		AccessKey:   metadata.FE.AccessKey,
		SecretKey:   metadata.FE.SecretKey,
		DomainName:  metadata.FE.DomainName,
		Region:      metadata.FE.Region,
		ProjectName: metadata.FE.ProjectName,
		Tags:        map[string]string{"openshift_cluster": metadata.ClusterID},
	}, nil
}

func (o *ClusterUninstaller) GenClient(ao hwsdk.AuthOptionsProvider) (*hwsdk.ProviderClient, error) {
	client, err := hwclient.NewClient(ao.GetIdentityEndpoint())
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{Proxy: http.ProxyFromEnvironment, TLSClientConfig: &tls.Config{}}

	client.HTTPClient = http.Client{
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if client.AKSKAuthOptions.AccessKey != "" {
				hwsdk.ReSign(req, hwsdk.SignOptions{
					AccessKey: client.AKSKAuthOptions.AccessKey,
					SecretKey: client.AKSKAuthOptions.SecretKey,
				})
			}
			return nil
		},
	}

	// Validate authentication normally.
	err = hwclient.Authenticate(client, ao)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (o *ClusterUninstaller) GetProviderClient() (*hwsdk.ProviderClient, error) {
	var pao, dao hwsdk.AKSKAuthOptions

	pao = hwsdk.AKSKAuthOptions{
		ProjectName: o.ProjectName,
	}

	dao = hwsdk.AKSKAuthOptions{
		Domain: o.DomainName,
	}

	for _, ao := range []*hwsdk.AKSKAuthOptions{&pao, &dao} {
		// TODO check auth_url
		ao.IdentityEndpoint = "https://iam.eu-west-0.prod-cloud-ocb.orange-business.com/v3"
		ao.AccessKey = o.AccessKey
		ao.SecretKey = o.SecretKey
	}

	return o.GenClient(pao)
}

// Run is the entrypoint to start the uninstall process
func (o *ClusterUninstaller) Run() (*types.ClusterQuota, error) {
	_, cancel := context.WithTimeout(o.Context, defaultTimeout)
	defer cancel()

	var err error
	providerClient, err := o.GetProviderClient()
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize provider client")
	}

	// opts := openstackdefaults.DefaultClientOpts(o.Cloud)
	o.NetworkClientV1, err = hwclient.NewNetworkV1(providerClient, hwsdk.EndpointOpts{Region: o.Region})
	o.NetworkClientV2, err = hwclient.NewNetworkV2(providerClient, hwsdk.EndpointOpts{Region: o.Region})
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize network client")
	}

	err = wait.PollImmediateInfinite(
		time.Second*10,
		o.destroyCluster,
	)

	if err != nil {
		return nil, errors.Wrap(err, "failed to destroy cluster")
	}

	return nil, nil
}

func (o *ClusterUninstaller) destroyCluster() (bool, error) {
	stagedFuncs := [][]struct {
		name    string
		execute func() error
	}{{
		{name: "Subnets", execute: o.destroySubnets},
		{name: "VPCs", execute: o.destroyVPCs},
	}}

	done := true
	for _, stage := range stagedFuncs {
		if done {
			for _, f := range stage {
				err := f.execute()
				if err != nil {
					o.Logger.Debugf("%s: %v", f.name, err)
					done = false
				}
			}
		}
	}
	return done, nil
}
