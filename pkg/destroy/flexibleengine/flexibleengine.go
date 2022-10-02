package flexibleengine

import (
	"context"
	"time"

	"github.com/openshift/installer/pkg/destroy/providers"
	"github.com/openshift/installer/pkg/types"

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

// Run is the entrypoint to start the uninstall process
func (o *ClusterUninstaller) Run() (*types.ClusterQuota, error) {
	_, cancel := context.WithTimeout(o.Context, defaultTimeout)
	defer cancel()

	err := wait.PollImmediateInfinite(
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
