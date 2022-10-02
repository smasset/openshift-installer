package flexibleengine

import (
	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/terraform/providers"
	"github.com/openshift/installer/pkg/terraform/stages"
)

// PlatformStages are the stages to run to provision the infrastructure in Flexible Engine.
var PlatformStages = []terraform.Stage{
	stages.NewStage(
		"flexibleengine",
		"cluster",
		[]providers.Provider{providers.FE},
	),
	stages.NewStage(
		"flexibleengine",
		"bootstrap",
		[]providers.Provider{providers.FE},
		stages.WithNormalBootstrapDestroy(),
	),
}
