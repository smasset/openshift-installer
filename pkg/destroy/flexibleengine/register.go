package flexibleengine

import "github.com/openshift/installer/pkg/destroy/providers"

func init() {
	providers.Registry["flexibleengine"] = New
}
