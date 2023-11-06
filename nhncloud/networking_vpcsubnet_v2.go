package nhncloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	gophercloud "github.com/nhn/nhncloud.gophercloud"
	"github.com/nhn/nhncloud.gophercloud/openstack/networking/v2/vpcsubnets"
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
		log.Printf("[DEBUG] Attempting to delete openstack_networking_vpcsubnet_v2 %s", subnetID)
		s, err := vpcsubnets.Get(networkingClient, subnetID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted openstack_networking_vpcsubnet_v2 %s", subnetID)
				return s, "DELETED", nil
			}
			return s, "ACTIVE", err
		}

		err = vpcsubnets.Delete(networkingClient, subnetID).ExtractErr()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted openstack_networking_vpcsubnet_v2 %s", subnetID)
				return s, "DELETED", nil
			}
			// Subnet is still in use - we can retry.
			if _, ok := err.(gophercloud.ErrDefault409); ok {
				return s, "ACTIVE", nil
			}
			return s, "ACTIVE", err
		}

		log.Printf("[DEBUG] openstack_networking_vpcsubnet_v2 %s is still active", subnetID)
		return s, "ACTIVE", nil
	}
}
