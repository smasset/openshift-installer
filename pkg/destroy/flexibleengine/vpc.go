package flexibleengine

import (
	"strings"

	hwtags "github.com/chnsz/golangsdk/openstack/common/tags"
	hwvpcs "github.com/chnsz/golangsdk/openstack/networking/v1/vpcs"
	hwpagination "github.com/chnsz/golangsdk/pagination"

	"github.com/pkg/errors"
)

func (opts ListByTagsOpts) ToVpcListQuery() (string, error) {
	return opts.ToListQuery()
}

func (o *ClusterUninstaller) destroyVPCs() error {
	o.Logger.Debug("Starting to destroy VPCs")

	err := o.deleteVPCs()
	if err != nil {
		errors.Wrap(err, "Failed to delete VPCs")
	}

	return err
}

func (o *ClusterUninstaller) deleteVPCs() error {
	// Look for tagged VPCs
	listOpts := ListByTagsOpts{Tags: strings.Join(GetTagNames(o.Tags), ",")}
	pager := o.ListVPCs(o.NetworkClientV1, listOpts)

	// Define an anonymous function to be executed on each page's iteration
	err := pager.EachPage(func(page hwpagination.Page) (bool, error) {
		vpcList, err := hwvpcs.ExtractVpcs(page)
		if err != nil {
			return false, err
		}

		for _, v := range vpcList {
			vpcTags, err := hwtags.Get(o.NetworkClientV2, "vpcs", v.ID).Extract()
			if err != nil {
				return false, err
			}

			// Check if tag values match
			vpcTagsMap := TagsToMap(vpcTags.Tags)
			if IsMapSubset(vpcTagsMap, o.Tags) {
				// TODO delete VPC
				o.Logger.Debugf("Deleting VPC %s", v.ID)
			}
		}

		return true, nil
	})

	return err
}
