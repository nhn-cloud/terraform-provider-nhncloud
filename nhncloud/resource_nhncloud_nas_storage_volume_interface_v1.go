package nhncloud

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/nhn-cloud/nhncloud.gophercloud/nhncloud/nas/v1/volumes"
)

func resourceNhncloudNasStorageVolumeInterfaceV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNhncloudNasStorageVolumeInterfaceV1Create,
		ReadContext:   resourceNhncloudNasStorageVolumeInterfaceV1Read,
		DeleteContext: resourceNhncloudNasStorageVolumeInterfaceV1Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"volume_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNhncloudNasStorageVolumeInterfaceV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	volumeID := d.Get("volume_id").(string)
	config.Lock(volumeID)
	defer config.Unlock(volumeID)

	nasStorageClient, err := config.NasStorageV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud NAS storage client: %s", err)
	}

	createOpts := &volumes.ConnectInterfaceOpts{
		SubnetID: d.Get("subnet_id").(string),
	}

	vInterface, err := volumes.ConnectInterface(nasStorageClient, volumeID, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud NAS storage volume interface: %s", err)
	}

	d.SetId(vInterface.ID)

	err = waitForCreateNhncloudNasStorageVolumeInterface(ctx, d, nasStorageClient, volumeID, vInterface.ID)
	if err != nil {
		return diag.Errorf("Error waiting for NHN Cloud NAS storage volume %s interface %s to become ready: %s", volumeID, vInterface.ID, err)
	}

	return resourceNhncloudNasStorageVolumeInterfaceV1Read(ctx, d, meta)
}

func resourceNhncloudNasStorageVolumeInterfaceV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	nasStorageClient, err := config.NasStorageV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud NAS storage client: %s", err)
	}

	volumeID := d.Get("volume_id").(string)
	vInterface, err := getNhncloudNasStorageVolumeInterfaceV1(nasStorageClient, volumeID, d.Id())
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, fmt.Sprintf("Error getting NHN Cloud NAS storage volume %s interface", volumeID)))
	}
	if vInterface == nil {
		d.SetId("")
		return nil
	}

	d.Set("subnet_id", vInterface.SubnetID)
	d.Set("path", vInterface.Path)
	d.Set("tenant_id", vInterface.TenantID)

	return nil
}

func resourceNhncloudNasStorageVolumeInterfaceV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	volumeID := d.Get("volume_id").(string)
	config.MutexKV.Lock(volumeID)
	defer config.MutexKV.Unlock(volumeID)

	nasStorageClient, err := config.NasStorageV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud NAS storage volume interface client: %s", err)
	}

	if err := volumes.DeleteInterface(nasStorageClient, volumeID, d.Id()).ExtractErr(); err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting nhncloud_nas_storage_volume_interface_v1"))
	}

	err = waitForDeleteNhncloudNasStorageVolumeInterface(ctx, d, nasStorageClient, volumeID, d.Id())
	if err != nil {
		return diag.Errorf("Error waiting for NHN Cloud NAS storage volume %s interface %s to delete: %s", volumeID, d.Id(), err)
	}

	return nil
}

func waitForCreateNhncloudNasStorageVolumeInterface(ctx context.Context, d *schema.ResourceData, nasStorageClient *gophercloud.ServiceClient, volumeID, interfaceID string) error {
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"BUILDING"},
		Target:     []string{"ACTIVE"},
		Refresh:    nasStorageVolumeInterfaceV1RefreshFunc(nasStorageClient, volumeID, interfaceID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func waitForDeleteNhncloudNasStorageVolumeInterface(ctx context.Context, d *schema.ResourceData, nasStorageClient *gophercloud.ServiceClient, volumeID, interfaceID string) error {
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE", "DELETING"},
		Target:     []string{"DELETED"},
		Refresh:    nasStorageVolumeInterfaceV1RefreshFunc(nasStorageClient, volumeID, interfaceID),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func nasStorageVolumeInterfaceV1RefreshFunc(client *gophercloud.ServiceClient, volumeID, interfaceID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		vInterface, err := getNhncloudNasStorageVolumeInterfaceV1(client, volumeID, interfaceID)
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return &volumes.Interface{}, "DELETED", nil
			}
			return &volumes.Interface{}, "", err
		}

		if vInterface == nil {
			return &volumes.Interface{}, "DELETED", nil
		}

		status := vInterface.Status
		if status == "error" {
			return vInterface, status, fmt.Errorf("The volume interface is in error status. " +
				"Please check with your cloud admin or check the NAS Storage " +
				"API logs to see why this error occurred.")
		}

		return vInterface, status, nil
	}
}

func getNhncloudNasStorageVolumeInterfaceV1(client *gophercloud.ServiceClient, volumeID, interfaceID string) (*volumes.Interface, error) {
	volume, err := volumes.Get(client, volumeID).Extract()
	if err != nil {
		return nil, err
	}

	index := slices.IndexFunc(volume.Interfaces, func(i *volumes.Interface) bool {
		return i.ID == interfaceID
	})
	if index == -1 {
		return nil, nil
	}

	return volume.Interfaces[index], nil
}
