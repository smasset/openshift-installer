package flexibleengine

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/core"
	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/types/flexibleengine"
	flexibleengineValidation "github.com/openshift/installer/pkg/types/flexibleengine/validation"
)

// Platform collects FE-specific configuration.
func Platform() (*flexibleengine.Platform, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	domain, err := getDomain(ctx)
	if err != nil {
		return nil, err
	}

	region, err := selectRegion(ctx)
	if err != nil {
		return nil, err
	}

	project, err := getProject(ctx)
	if err != nil {
		return nil, err
	}

	accessKey, secretKey, err := getAuthentication(ctx)
	if err != nil {
		return nil, err
	}

	return &flexibleengine.Platform{
		AccessKey:   accessKey,
		SecretKey:   secretKey,
		DomainName:  domain,
		Region:      region,
		ProjectName: project,
	}, nil
}

func getDomain(ctx context.Context) (string, error) {
	var domain string

	if err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Domain",
				Help:    "The FE domain where the cluster will be provisioned.",
			},
			Validate: survey.Required,
		},
	}, &domain); err != nil {
		return domain, errors.Wrap(err, "failed UserInput")
	}

	return domain, nil
}

func selectRegion(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	validRegions := flexibleengineValidation.Regions

	defaultRegion := "eu-west-0"
	defaultRegionName := ""
	longRegions := make([]string, 0, len(validRegions))
	shortRegions := make([]string, 0, len(validRegions))
	for key, value := range validRegions {
		shortRegions = append(shortRegions, key)
		regionDesc := fmt.Sprintf("%s (%s)", key, value)
		longRegions = append(longRegions, regionDesc)

		if defaultRegionName == "" && key == defaultRegion {
			defaultRegionName = regionDesc
		}
	}

	var regionTransform survey.Transformer = func(ans interface{}) interface{} {
		switch v := ans.(type) {
		case core.OptionAnswer:
			return core.OptionAnswer{Value: strings.SplitN(v.Value, " ", 2)[0], Index: v.Index}
		case string:
			return strings.SplitN(v, " ", 2)[0]
		}
		return ""
	}

	sort.Strings(longRegions)
	sort.Strings(shortRegions)

	if defaultRegionName == "" && len(longRegions) > 0 {
		defaultRegionName = longRegions[0]
	}

	var selectedRegion string
	err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Select{
				Message: "Region",
				Help:    "The FE region to be used for installation.",
				Default: defaultRegionName,
				Options: longRegions,
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				choice := regionTransform(ans).(core.OptionAnswer).Value
				i := sort.SearchStrings(shortRegions, choice)
				if i == len(shortRegions) || shortRegions[i] != choice {
					return errors.Errorf("invalid region %q", choice)
				}
				return nil
			}),
			Transform: regionTransform,
		},
	}, &selectedRegion)
	if err != nil {
		return "", err
	}

	return selectedRegion, nil
}

func getProject(ctx context.Context) (string, error) {
	var project string

	if err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Project",
				Help:    "The FE project where the cluster will be provisioned.",
			},
		},
	}, &project); err != nil {
		return project, errors.Wrap(err, "failed UserInput")
	}

	return project, nil
}

func getAuthentication(ctx context.Context) (string, string, error) {
	var accessKey, secretKey string

	if err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Access Key",
				Help:    "The access key to login to Flexible Engine.",
			},
			Validate: survey.Required,
		},
	}, &accessKey); err != nil {
		return accessKey, secretKey, errors.Wrap(err, "failed UserInput")
	}

	if err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Password{
				Message: "Secret Key",
				Help:    "The secret key to login to Flexible Engine.",
			},
			Validate: survey.Required,
		},
	}, &secretKey); err != nil {
		return accessKey, secretKey, errors.Wrap(err, "failed UserInput")
	}

	return accessKey, secretKey, nil
}
