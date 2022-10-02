package validation

import (
	"sort"

	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/flexibleengine"
)

var (
	// Regions is a map of known FE regions. The key of the map is
	// the short name of the region. The value of the map is the long
	// name of the region.
	Regions = map[string]string{
		"eu-west-0": "Paris, France",
		"eu-west-1": "Amsterdam, Netherlands",
	}
	validRegionValues = func() []string {
		validValues := make([]string, len(Regions))
		i := 0
		for r := range Regions {
			validValues[i] = r
			i++
		}
		sort.Strings(validValues)
		return validValues
	}()
)

// ValidatePlatform checks that the specified platform is valid.
func ValidatePlatform(p *flexibleengine.Platform, fldPath *field.Path, ic *types.InstallConfig) field.ErrorList {
	allErrs := field.ErrorList{}

	if p.AccessKey == "" {
		allErrs = append(allErrs, field.Required(fldPath.Child("accessKey"), "must provide an access key"))
	}

	if p.SecretKey == "" {
		allErrs = append(allErrs, field.Required(fldPath.Child("secretKey"), "must provide a secret key"))
	}

	if p.DomainName == "" {
		allErrs = append(allErrs, field.Required(fldPath.Child("domain"), "must provide a domain name"))
	}

	if p.Region == "" {
		allErrs = append(allErrs, field.Required(fldPath.Child("region"), "must provide a region"))
	}

	return allErrs
}
