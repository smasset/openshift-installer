package flexibleengine

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
)

// Validate executes platform-specific validation.
func Validate(config *types.InstallConfig) error {
	allErrs := field.ErrorList{}

	return allErrs.ToAggregate()
}
