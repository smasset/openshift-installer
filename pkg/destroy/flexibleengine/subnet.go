package flexibleengine

import (
	"strings"

	hwtags "github.com/chnsz/golangsdk/openstack/common/tags"
	hwsubnets "github.com/chnsz/golangsdk/openstack/networking/v2/subnets"
	hwpagination "github.com/chnsz/golangsdk/pagination"

	"github.com/pkg/errors"
)

func (opts ListByTagsOpts) ToSubnetListQuery() (string, error) {
	return opts.ToListQuery()
}

func (o *ClusterUninstaller) destroySubnets() error {
	o.Logger.Debug("Starting to destroy subnets")

	err := o.deleteSubnets()
	if err != nil {
		errors.Wrap(err, "Failed to delete subnets")
	}

	return err
}

func (o *ClusterUninstaller) deleteSubnets() error {
	// Look for tagged subnets
	listOpts := ListByTagsOpts{Tags: strings.Join(GetTagNames(o.Tags), ",")}
	pager := hwsubnets.List(o.NetworkClientV2, listOpts)

	// Define an anonymous function to be executed on each page's iteration
	err := pager.EachPage(func(page hwpagination.Page) (bool, error) {
		subnetList, err := hwsubnets.ExtractSubnets(page)
		if err != nil {
			return false, err
		}

		for _, s := range subnetList {
			subnetTags, err := hwtags.Get(o.NetworkClientV2, "subnets", s.NetworkID).Extract()
			if err != nil {
				return false, err
			}

			// Check if tag values match
			subnetTagsMap := TagsToMap(subnetTags.Tags)
			if IsMapSubset(subnetTagsMap, o.Tags) {
				// TODO delete subnet
				o.Logger.Debugf("Deleting subnet %s", s.ID)
			}
		}

		return true, nil
	})

	return err
}
