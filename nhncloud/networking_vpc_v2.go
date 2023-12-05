package nhncloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/gophercloud/gophercloud"
	"github.com/nhn/nhncloud.gophercloud/nhncloud/networking/v2/vpcs"
)

type vpcExtended struct {
	vpcs.Network
}

type vpcExtendedDetail struct {
	vpcs.NetworkDetail
}

func resourceNetworkingVPCV2StateRefreshFunc(client *gophercloud.ServiceClient, networkID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := vpcs.Get(client, networkID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return n, "DELETED", nil
			}
			if _, ok := err.(gophercloud.ErrDefault409); ok {
				return n, "ACTIVE", nil
			}

			return n, "", err
		}

		return n, n.State, nil
	}
}
