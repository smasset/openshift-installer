package flexibleengine

import (
	"fmt"

	machineapi "github.com/openshift/api/machine/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	openstackprovider "sigs.k8s.io/cluster-api-provider-openstack/pkg/apis/openstackproviderconfig/v1alpha1"

	"github.com/openshift/installer/pkg/types"
)

// Machines returns a list of machines for a machinepool.
func Machines(clusterID string, config *types.InstallConfig, pool *types.MachinePool, osImage, role, userDataSecret string) ([]machineapi.Machine, error) {
	var machines []machineapi.Machine

	provider, err := provider()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create provider")
	}

	machine := machineapi.Machine{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "machine.openshift.io/v1beta1",
			Kind:       "Machine",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "openshift-machine-api",
			Name:      fmt.Sprintf("%s-%s-%d", clusterID, pool.Name, 1),
			Labels: map[string]string{
				"machine.openshift.io/cluster-api-cluster":      clusterID,
				"machine.openshift.io/cluster-api-machine-role": role,
				"machine.openshift.io/cluster-api-machine-type": role,
			},
		},
		Spec: machineapi.MachineSpec{
			ProviderSpec: machineapi.ProviderSpec{
				Value: &runtime.RawExtension{Object: provider},
			},
		},
	}
	machines = append(machines, machine)

	return machines, nil
}

func provider() (*openstackprovider.OpenstackProviderSpec, error) {
	return &openstackprovider.OpenstackProviderSpec{
		TypeMeta: metav1.TypeMeta{
			APIVersion: openstackprovider.SchemeGroupVersion.String(),
			Kind:       "OpenstackProviderSpec",
		},
	}, nil
}
