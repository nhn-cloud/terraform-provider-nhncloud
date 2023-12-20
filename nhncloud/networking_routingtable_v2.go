package nhncloud

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/routers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/nhn/nhncloud.gophercloud/nhncloud/networking/v2/routingtables"
)

type routingtableExtended struct {
	routingtables.Routingtable
}

type routingtableExtendedDetail struct {
	routingtables.RoutingtableDetail
}

func resourceNetworkingRoutingtableV2StateRefreshFunc(client *gophercloud.ServiceClient, routingtableID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := routers.Get(client, routingtableID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return n, "DELETED", nil
			}

			return n, "", err
		}

		return n, n.Status, nil
	}
}
