package nhncloud

import (
	"context"
	"fmt"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/nhn-cloud/nhncloud.gophercloud/nhncloud/nas/v1/volumes"
)

func resourceNasStorageVolumeV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNasStorageVolumeV1Create,
		ReadContext:   resourceNasStorageVolumeV1Read,
		UpdateContext: resourceNasStorageVolumeV1Update,
		DeleteContext: resourceNasStorageVolumeV1Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"size_gb": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"acl": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				DefaultFunc: func() (any, error) {
					return []string{}, nil
				},
			},

			"encryption": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
						"keys": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"key_store_id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
					},
				},
				DiffSuppressFunc: nasStorageVolumeV1ListDiffSuppressFunc("encryption.#"),
			},

			"mount_protocol": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"cifs_auth_ids": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
					},
				},
			},

			"snapshot_policy": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max_scheduled_count": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"reserve_percent": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  5,
						},
						"schedule": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"time": {
										Type:     schema.TypeString,
										Required: true,
									},
									"time_offset": {
										Type:     schema.TypeString,
										Required: true,
									},
									"weekdays": {
										Type:     schema.TypeSet,
										Required: true,
										Elem:     &schema.Schema{Type: schema.TypeInt},
										Set:      schema.HashInt,
									},
								},
							},
						},
					},
				},
				DiffSuppressFunc: nasStorageVolumeV1ListDiffSuppressFunc("snapshot_policy.#"),
			},

			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNasStorageVolumeV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	nasStorageClient, err := config.NasStorageV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud NAS storage client: %s", err)
	}

	createOpts := &volumes.CreateOpts{
		Name:           d.Get("name").(string),
		Description:    d.Get("description").(string),
		SizeGb:         d.Get("size_gb").(int),
		ACL:            resourceToNasStorageVolumeACL(d),
		Encryption:     resourceToNasStorageVolumeEncryption(d),
		MountProtocol:  resourceToNasStorageVolumeMountProtocol(d),
		SnapshotPolicy: resourceToNasStorageVolumeSnapshotPolicy(d),
	}

	v, err := volumes.Create(nasStorageClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud NAS storage volume: %s", err)
	}
	d.SetId(v.ID)

	err = waitForCreateNasStorageVolume(ctx, d, nasStorageClient, v.ID)
	if err != nil {
		return diag.Errorf("Error waiting for NHN Cloud NAS storage volume %s to become ready: %s", v.ID, err)
	}

	return resourceNasStorageVolumeV1Read(ctx, d, meta)
}

func resourceNasStorageVolumeV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	nasStorageClient, err := config.NasStorageV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud NAS storage client: %s", err)
	}

	v, err := volumes.Get(nasStorageClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving nhncloud_nas_storage_volume_v1"))
	}

	d.Set("region", GetRegion(d, config))
	d.Set("name", v.Name)
	d.Set("description", v.Description)
	d.Set("size_gb", v.SizeGb)
	d.Set("project_id", v.ProjectID)
	d.Set("tenant_id", v.TenantID)
	d.Set("acl", v.ACL)
	d.Set("encryption", flattenNasStorageVolumeEncryption(v.Encryption))
	d.Set("mount_protocol", flattenNasStorageVolumeMountProtocol(v.MountProtocol))
	d.Set("snapshot_policy", flattenNasStorageVolumeSnapshotPolicy(v.SnapshotPolicy))
	d.Set("created_at", v.CreatedAt)
	d.Set("updated_at", v.UpdatedAt)

	return nil
}

func resourceNasStorageVolumeV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	nasStorageClient, err := config.NasStorageV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud NAS storage client: %s", err)
	}

	updateOpts := &volumes.UpdateOpts{}

	if d.HasChange("description") {
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}
	if d.HasChange("size_gb") {
		sizeGb := d.Get("size_gb").(int)
		updateOpts.SizeGb = &sizeGb
	}
	if d.HasChange("acl") {
		acl := resourceToNasStorageVolumeACL(d)
		updateOpts.ACL = &acl
	}
	if d.HasChange("mount_protocol") {
		mountProtocol := resourceToNasStorageVolumeMountProtocol(d)
		updateOpts.MountProtocol = mountProtocol
	}
	if d.HasChange("snapshot_policy") {
		snapshotPolicy := resourceToNasStorageVolumeSnapshotPolicy(d)
		updateOpts.SnapshotPolicy = snapshotPolicy
	}
	_, err = volumes.Update(nasStorageClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return diag.Errorf("Error updating NHN Cloud NAS storage volume: %s", err)
	}

	err = waitForUpdateNasStorageVolume(ctx, d, nasStorageClient, d.Id())
	if err != nil {
		return diag.Errorf("Error waiting for NHN Cloud NAS storage volume %s to update: %s", d.Id(), err)
	}

	return resourceNasStorageVolumeV1Read(ctx, d, meta)
}

func resourceNasStorageVolumeV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	nasStorageClient, err := config.NasStorageV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud NAS storage client: %s", err)
	}

	if err := volumes.Delete(nasStorageClient, d.Id()).ExtractErr(); err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting nhncloud_nas_storage_volume_v1"))
	}

	err = waitForDeleteNasStorageVolume(ctx, d, nasStorageClient, d.Id())
	if err != nil {
		return diag.Errorf(
			"Error waiting for nhncloud_nas_storage_volume_v1 %s to delete: %s", d.Id(), err)
	}

	return nil
}

func resourceToNasStorageVolumeACL(d *schema.ResourceData) []string {
	rawACLs := d.Get("acl").(*schema.Set).List()
	acls := make([]string, len(rawACLs))
	for i, raw := range rawACLs {
		acls[i] = raw.(string)
	}
	return acls
}

func resourceToNasStorageVolumeEncryption(d *schema.ResourceData) *volumes.EncryptionOpts {
	encryptionList := d.Get("encryption").([]interface{})
	if len(encryptionList) == 0 {
		return nil
	}

	encryption := encryptionList[0].(map[string]any)
	return &volumes.EncryptionOpts{
		Enabled: encryption["enabled"].(bool),
	}
}

func resourceToNasStorageVolumeMountProtocol(d *schema.ResourceData) *volumes.MountProtocolOpts {
	rawMountProtocolList := d.Get("mount_protocol").([]interface{})
	if len(rawMountProtocolList) == 0 {
		return nil
	}

	rawMountProtocol := rawMountProtocolList[0].(map[string]any)
	rawCifsAuthIDs := rawMountProtocol["cifs_auth_ids"].(*schema.Set).List()
	cifsAuthIDs := make([]string, len(rawCifsAuthIDs))
	for i, raw := range rawCifsAuthIDs {
		cifsAuthIDs[i] = raw.(string)
	}

	return &volumes.MountProtocolOpts{
		CifsAuthIDs: cifsAuthIDs,
		Protocol:    rawMountProtocol["protocol"].(string),
	}
}

func resourceToNasStorageVolumeSnapshotPolicy(d *schema.ResourceData) *volumes.SnapshotPolicyOpts {
	snapshotPolicyList := d.Get("snapshot_policy").([]interface{})
	if len(snapshotPolicyList) == 0 {
		return nil
	}
	snapshotPolicy := snapshotPolicyList[0].(map[string]any)

	opts := &volumes.SnapshotPolicyOpts{}

	maxScheduledCount := snapshotPolicy["max_scheduled_count"].(int)
	if maxScheduledCount > 0 {
		opts.MaxScheduledCount = &maxScheduledCount
	}
	opts.ReservePercent = snapshotPolicy["reserve_percent"].(int)
	opts.Schedule = resourceToNasStorageVolumeSnapshotPolicySchedule(snapshotPolicy)

	return opts
}

func resourceToNasStorageVolumeSnapshotPolicySchedule(snapshotPolicy map[string]any) *volumes.ScheduleOpts {
	scheduleList := snapshotPolicy["schedule"].([]interface{})
	if len(scheduleList) == 0 {
		return nil
	}

	schedule := scheduleList[0].(map[string]any)

	rawWeekdays := schedule["weekdays"].(*schema.Set).List()
	weekdays := make([]int, len(rawWeekdays))
	for i, raw := range rawWeekdays {
		weekdays[i] = raw.(int)
	}

	return &volumes.ScheduleOpts{
		Time:       schedule["time"].(string),
		TimeOffset: schedule["time_offset"].(string),
		Weekdays:   weekdays,
	}
}

func flattenNasStorageVolumeEncryption(encryption *volumes.Encryption) []any {
	if encryption == nil {
		return nil
	}

	flattenEncryption := map[string]any{}
	flattenEncryption["enabled"] = encryption.Enabled
	if len(encryption.Keys) > 0 {
		var keys []map[string]any = nil
		for _, rawKey := range encryption.Keys {
			keys = append(keys, map[string]any{
				"key_id":       rawKey.KeyID,
				"key_store_id": rawKey.KeyStoreID,
			})
		}
		flattenEncryption["keys"] = keys
	}

	return []any{flattenEncryption}
}

func flattenNasStorageVolumeMountProtocol(mountProtocol *volumes.MountProtocol) []any {
	if mountProtocol == nil {
		return nil
	}

	return []any{map[string]any{
		"protocol":      mountProtocol.Protocol,
		"cifs_auth_ids": mountProtocol.CifsAuthIDs,
	}}
}

func flattenNasStorageVolumeSnapshotPolicy(snapshotPolicy *volumes.SnapshotPolicy) []any {
	if snapshotPolicy == nil {
		return nil
	}

	return []any{map[string]any{
		"max_scheduled_count": snapshotPolicy.MaxScheduledCount,
		"reserve_percent":     snapshotPolicy.ReservePercent,
		"schedule":            flattenNasStorageVolumeSnapshotPolicySchedule(snapshotPolicy.Schedule),
	}}
}

func flattenNasStorageVolumeSnapshotPolicySchedule(schedule *volumes.Schedule) []any {
	if schedule == nil {
		return nil
	}

	return []any{map[string]any{
		"time":        schedule.Time,
		"time_offset": schedule.TimeOffset,
		"weekdays":    schedule.Weekdays,
	}}
}

func waitForCreateNasStorageVolume(ctx context.Context, d *schema.ResourceData, nasStorageClient *gophercloud.ServiceClient, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"BUILDING", "INITIALIZING"},
		Target:     []string{"ACTIVE"},
		Refresh:    nasStorageVolumeV1RefreshFunc(nasStorageClient, id),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func waitForUpdateNasStorageVolume(ctx context.Context, d *schema.ResourceData, nasStorageClient *gophercloud.ServiceClient, id string) error {
	stateConf := &retry.StateChangeConf{
		Pending:    []string{"UPDATING"},
		Target:     []string{"ACTIVE"},
		Refresh:    nasStorageVolumeV1RefreshFunc(nasStorageClient, id),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func waitForDeleteNasStorageVolume(ctx context.Context, d *schema.ResourceData, nasStorageClient *gophercloud.ServiceClient, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE", "DELETING"},
		Target:     []string{"DELETED"},
		Refresh:    nasStorageVolumeV1RefreshFunc(nasStorageClient, id),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func nasStorageVolumeV1RefreshFunc(client *gophercloud.ServiceClient, id string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		v, err := volumes.Get(client, id).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return v, "DELETED", nil
			}
			return nil, "", err
		}

		if v.Status == "error" {
			return v, v.Status, fmt.Errorf("The volume is in error status: " +
				"Please check with your cloud admin or check the Nas Storage " +
				"API logs to see why this error occurred")
		}

		return v, v.Status, nil
	}
}

func nasStorageVolumeV1ListDiffSuppressFunc(target string) func(string, string, string, *schema.ResourceData) bool {
	return func(k, old, new string, d *schema.ResourceData) bool {
		if k == target {
			if old == "1" && new == "0" {
				return true
			}
		}
		return old == new
	}
}
