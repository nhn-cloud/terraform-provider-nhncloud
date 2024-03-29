package nhncloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/gophercloud/gophercloud"
	"github.com/nhn-cloud/nhncloud.gophercloud/nhncloud/networking/v2/vpcsubnets"
)

func networkingVPCSubnetV2StateRefreshFunc(client *gophercloud.ServiceClient, subnetID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		subnet, err := vpcsubnets.Get(client, subnetID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return subnet, "DELETED", nil
			}
			return nil, "", err
		}
		return subnet, "ACTIVE", nil
	}
}

func networkingVPCSubnetV2StateRefreshFuncDelete(networkingClient *gophercloud.ServiceClient, subnetID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to delete nhncloud_networking_vpcsubnet_v2 %s", subnetID)
		s, err := vpcsubnets.Get(networkingClient, subnetID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted nhncloud_networking_vpcsubnet_v2 %s", subnetID)
				return s, "DELETED", nil
			}
			return s, "ACTIVE", err
		}

		err = vpcsubnets.Delete(networkingClient, subnetID).ExtractErr()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted nhncloud_networking_vpcsubnet_v2 %s", subnetID)
				return s, "DELETED", nil
			}
			// Subnet is still in use - we can retry.
			if _, ok := err.(gophercloud.ErrDefault409); ok {
				return s, "ACTIVE", nil
			}
			return s, "ACTIVE", err
		}

		log.Printf("[DEBUG] nhncloud_networking_vpcsubnet_v2 %s is still active", subnetID)
		return s, "ACTIVE", nil
	}
}
