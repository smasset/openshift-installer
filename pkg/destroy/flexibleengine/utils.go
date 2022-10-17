package flexibleengine

import (
	hwsdk "github.com/chnsz/golangsdk"
	hwtags "github.com/chnsz/golangsdk/openstack/common/tags"
	hwvpcs "github.com/chnsz/golangsdk/openstack/networking/v1/vpcs"
	hwpagination "github.com/chnsz/golangsdk/pagination"
)

type ListByTagsOptsBuilder interface {
	ToListQuery() (string, error)
}

type ListByTagsOpts struct {
	Tags string `q:"tags"`
}

func (opts ListByTagsOpts) ToListQuery() (string, error) {
	q, err := hwsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), nil
}

func (o *ClusterUninstaller) ListVPCs(c *hwsdk.ServiceClient, opts ListByTagsOptsBuilder) hwpagination.Pager {
	url := c.ServiceURL(c.ProjectID, "vpcs")
	if opts != nil {
		query, err := opts.ToListQuery()
		if err != nil {
			return hwpagination.Pager{Err: err}
		}
		url += query
	}
	return hwpagination.NewPager(c, url, func(r hwpagination.PageResult) hwpagination.Page {
		return hwvpcs.VpcPage{hwpagination.LinkedPageBase{PageResult: r}}
	})
}

func GetTagNames(tagMap map[string]string) [](string) {
	var tags [](string)
	for k := range tagMap {
		tags = append(tags, k)
	}
	return tags
}

func TagsToMap(tags []hwtags.ResourceTag) map[string]string {
	result := make(map[string]string)
	for _, val := range tags {
		result[val.Key] = val.Value
	}

	return result
}

func IsMapSubset[K, V comparable](m, sub map[K]V) bool {
	if len(sub) > len(m) {
		return false
	}
	for k, vsub := range sub {
		if vm, found := m[k]; !found || vm != vsub {
			return false
		}
	}
	return true
}
