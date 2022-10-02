package flexibleengine

import (
	"github.com/openshift/installer/pkg/types"
	flexibleenginetypes "github.com/openshift/installer/pkg/types/flexibleengine"
)

// Metadata converts an install configuration to FE metadata.
func Metadata(config *types.InstallConfig) *flexibleenginetypes.Metadata {
	return &flexibleenginetypes.Metadata{
		AccessKey:   config.Platform.FE.AccessKey,
		SecretKey:   config.Platform.FE.SecretKey,
		DomainName:  config.Platform.FE.DomainName,
		Region:      config.Platform.FE.Region,
		ProjectName: config.Platform.FE.ProjectName,
	}
}
