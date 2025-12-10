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

func resourceNhncloudNasStorageVolumeMirrorV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNhncloudNasStorageVolumeMirrorV1Create,
		ReadContext:   resourceNhncloudNasStorageVolumeMirrorV1Read,
		UpdateContext: resourceNhncloudNasStorageVolumeMirrorV1Update,
		DeleteContext: resourceNhncloudNasStorageVolumeMirrorV1Delete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"src_region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"src_volume_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"dst_region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"dst_tenant_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"dst_volume": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
								},
							},
							DiffSuppressFunc: nasStorageVolumeV1ListDiffSuppressFunc("dst_volume.0.encryption.#"),
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
							DiffSuppressFunc: nasStorageVolumeV1ListDiffSuppressFunc("dst_volume.0.snapshot_policy.#"),
						},
					},
				},
			},

			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"role": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"direction": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"direction_changed_at": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"dst_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"dst_volume_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"dst_volume_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"src_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"src_tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"src_volume_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNhncloudNasStorageVolumeMirrorV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	srcRegion := getNasStorageVolumeMirrorV1SrcRegion(d, config)
	srcVolumeID := d.Get("src_volume_id").(string)
	config.MutexKV.Lock(srcVolumeID)
	defer config.MutexKV.Unlock(srcVolumeID)

	nasStorageClient, err := config.NasStorageV1Client(srcRegion)
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud NAS storage client: %s", err)
	}

	dstVolume := d.Get("dst_volume").([]any)[0].(map[string]any)
	createOpts := &volumes.SetReplicationOpts{
		DstRegion:   d.Get("dst_region").(string),
		DstTenantID: d.Get("dst_tenant_id").(string),
		DstVolume: &volumes.DstVolumeOpts{
			Name:           dstVolume["name"].(string),
			Description:    dstVolume["description"].(string),
			SizeGb:         dstVolume["size_gb"].(int),
			ACL:            resourceNhncloudNasStorageVolumeMirrorACL(dstVolume),
			Encryption:     resourceNhncloudNasStorageVolumeMirrorEncryption(dstVolume),
			MountProtocol:  resourceNhncloudNasStorageVolumeMirrorMountProtocol(dstVolume),
			SnapshotPolicy: resourceNhncloudNasStorageVolumeMirrorSnapshotPolicy(dstVolume),
		},
	}

	v, err := volumes.SetReplication(nasStorageClient, srcVolumeID, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud NAS storage volume mirror: %s", err)
	}

	d.SetId(v.ID)

	err = waitForCreateNhncloudNasStorageVolumeMirror(ctx, d, nasStorageClient, srcVolumeID, v.ID)
	if err != nil {
		return diag.Errorf("Error waiting for NHN Cloud NAS storage volume mirror %s to become ready: %s", d.Id(), err)
	}

	return resourceNhncloudNasStorageVolumeMirrorV1Read(ctx, d, meta)
}

func resourceNhncloudNasStorageVolumeMirrorV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	mirror, err := getNhncloudNasStorageVolumeMirrorV1(d, config)
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, fmt.Sprintf("Error getting NHN Cloud NAS storage volume %s mirror", d.Get("src_volume_id").(string))))
	}
	if mirror == nil {
		d.SetId("")
		return nil
	}

	dstVolume, err := getNhncloudNasStorageVolumeMirrorV1DstVolume(config, mirror.DstRegion, mirror.DstVolumeID)
	if err != nil {
		return diag.Errorf("Error getting NHN Cloud NAS storage volume mirror %s: %s", d.Id(), err)
	}

	d.Set("role", mirror.Role)
	d.Set("direction", mirror.Direction)
	d.Set("direction_changed_at", mirror.DirectionChangedAt)
	d.Set("dst_project_id", mirror.DstProjectID)
	d.Set("dst_region", mirror.DstRegion)
	d.Set("dst_tenant_id", mirror.DstTenantID)
	d.Set("dst_volume_id", mirror.DstVolumeID)
	d.Set("dst_volume_name", mirror.DstVolumeName)
	d.Set("src_project_id", mirror.SrcProjectID)
	d.Set("src_region", mirror.SrcRegion)
	d.Set("src_tenant_id", mirror.SrcTenantID)
	d.Set("src_volume_id", mirror.SrcVolumeID)
	d.Set("src_volume_name", mirror.SrcVolumeName)
	d.Set("created_at", mirror.CreatedAt)
	d.Set("dst_volume", flattenNhncloudNasStorageVolumeMirrorDstVolume(dstVolume))

	return nil
}

func resourceNhncloudNasStorageVolumeMirrorV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	dstRegionID := d.Get("dst_region").(string)
	nasStorageClient, err := config.NasStorageV1Client(dstRegionID)
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud NAS storage client: %s", err)
	}

	updateOpts := &volumes.UpdateOpts{}

	if d.HasChange("dst_volume.0.description") {
		description := d.Get("dst_volume.0.description").(string)
		updateOpts.Description = &description
	}

	dstVolumeID := d.Get("dst_volume_id").(string)
	_, err = volumes.Update(nasStorageClient, dstVolumeID, updateOpts).Extract()
	if err != nil {
		return diag.Errorf("Error updating NHN Cloud NAS storage volume mirror: %s", err)
	}

	err = waitForUpdateNhncloudNasStorageVolume(ctx, d, nasStorageClient, dstVolumeID)
	if err != nil {
		return diag.Errorf("Error waiting for NHN Cloud NAS storage volume %s to update: %s", d.Id(), err)
	}

	return resourceNhncloudNasStorageVolumeMirrorV1Read(ctx, d, meta)
}

func resourceNhncloudNasStorageVolumeMirrorV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	srcRegion := getNasStorageVolumeMirrorV1SrcRegion(d, config)
	srcVolumeID := d.Get("src_volume_id").(string)
	config.Lock(srcVolumeID)
	defer config.Unlock(srcVolumeID)

	nasStorageClient, err := config.NasStorageV1Client(srcRegion)
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud NAS storage client: %s", err)
	}

	if err := volumes.DisableReplication(nasStorageClient, srcVolumeID, d.Id()).ExtractErr(); err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting nhncloud_nas_storage_volume_mirror_v1"))
	}

	err = waitForDeleteNhncloudNasStorageVolumeMirror(ctx, d, nasStorageClient, srcVolumeID, d.Id())
	if err != nil {
		return diag.Errorf("Error waiting for NHN Cloud NAS storage volume mirror %s to be deleted: %s", d.Id(), err)
	}

	return nil
}

func resourceNhncloudNasStorageVolumeMirrorACL(dstVolume map[string]any) []string {
	rawACLs := dstVolume["acl"].(*schema.Set).List()
	acls := make([]string, len(rawACLs))
	for i, raw := range rawACLs {
		acls[i] = raw.(string)
	}
	return acls
}

func resourceNhncloudNasStorageVolumeMirrorEncryption(dstVolume map[string]any) *volumes.EncryptionOpts {
	rawEncryptionList := dstVolume["encryption"].([]any)
	if len(rawEncryptionList) == 0 {
		return nil
	}

	rawEncryption := rawEncryptionList[0].(map[string]any)
	encryption := &volumes.EncryptionOpts{
		Enabled: rawEncryption["enabled"].(bool),
	}
	return encryption
}

func resourceNhncloudNasStorageVolumeMirrorMountProtocol(dstVolume map[string]any) *volumes.MountProtocolOpts {
	rawMountProtocolList := dstVolume["mount_protocol"].([]any)
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
		Protocol:    rawMountProtocol["protocol"].(string),
		CifsAuthIDs: cifsAuthIDs,
	}
}

func resourceNhncloudNasStorageVolumeMirrorSnapshotPolicy(dstVolume map[string]any) *volumes.SnapshotPolicyOpts {
	rawSnapshotPolicyList := dstVolume["snapshot_policy"].([]any)
	if len(rawSnapshotPolicyList) == 0 {
		return nil
	}
	rawSnapshotPolicy := rawSnapshotPolicyList[0].(map[string]any)

	opts := &volumes.SnapshotPolicyOpts{}
	maxScheduledCount := rawSnapshotPolicy["max_scheduled_count"].(int)
	if maxScheduledCount > 0 {
		opts.MaxScheduledCount = &maxScheduledCount
	}
	opts.ReservePercent = rawSnapshotPolicy["reserve_percent"].(int)
	opts.Schedule = resourceNhncloudNasStorageVolumeMirrorSnapshotPolicySchedule(rawSnapshotPolicy)

	return opts
}

func resourceNhncloudNasStorageVolumeMirrorSnapshotPolicySchedule(snapshotPolicy map[string]any) *volumes.ScheduleOpts {
	scheduleList := snapshotPolicy["schedule"].([]any)
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

func waitForCreateNhncloudNasStorageVolumeMirror(ctx context.Context, d *schema.ResourceData, nasStorageClient *gophercloud.ServiceClient, volumeID, mirrorID string) error {
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING"},
		Target:     []string{"ACTIVE"},
		Refresh:    nasStorageVolumeMirrorV1RefreshFunc(nasStorageClient, volumeID, mirrorID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func waitForDeleteNhncloudNasStorageVolumeMirror(ctx context.Context, d *schema.ResourceData, nasStorageClient *gophercloud.ServiceClient, volumeID, mirrorID string) error {
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE", "DELETING"},
		Target:     []string{"DELETED"},
		Refresh:    nasStorageVolumeMirrorV1RefreshFunc(nasStorageClient, volumeID, mirrorID),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func nasStorageVolumeMirrorV1RefreshFunc(client *gophercloud.ServiceClient, volumeID, mirrorID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		m, err := volumes.GetReplicationStat(client, volumeID, mirrorID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return m, "DELETED", nil
			}
			return nil, "", err
		}

		if m.Status == "error" {
			return m, m.Status, fmt.Errorf("The volume mirror is in error status. " +
				"Please check with your cloud admin or check the NAS Storage " +
				"API logs to see why this error occurred.")
		}

		return m, m.Status, nil
	}
}

func getNhncloudNasStorageVolumeMirrorV1(d *schema.ResourceData, config *Config) (*volumes.Mirror, error) {
	srcRegion := getNasStorageVolumeMirrorV1SrcRegion(d, config)
	nasStorageClient, err := config.NasStorageV1Client(srcRegion)
	if err != nil {
		return nil, err
	}

	srcVolumeID := d.Get("src_volume_id").(string)
	volume, err := volumes.Get(nasStorageClient, srcVolumeID).Extract()
	if err != nil {
		return nil, err
	}

	index := slices.IndexFunc(volume.Mirrors, func(i *volumes.Mirror) bool {
		return i.ID == d.Id()
	})
	if index == -1 {
		return nil, nil
	}

	return volume.Mirrors[index], nil
}

func getNhncloudNasStorageVolumeMirrorV1DstVolume(config *Config, dstRegionID, dstVolumeID string) (*volumes.Volume, error) {
	nasStorageClient, err := config.NasStorageV1Client(dstRegionID)
	if err != nil {
		return nil, err
	}

	volume, err := volumes.Get(nasStorageClient, dstVolumeID).Extract()
	if err != nil {
		return nil, err
	}

	return volume, nil
}

func flattenNhncloudNasStorageVolumeMirrorDstVolume(dstVolume *volumes.Volume) []any {
	return []any{map[string]any{
		"name":        dstVolume.Name,
		"description": dstVolume.Description,
		"size_gb":     dstVolume.SizeGb,
		"acl":         dstVolume.ACL,
		"encryption": []any{map[string]any{
			"enabled": dstVolume.Encryption.Enabled,
		}},
		"mount_protocol": []any{map[string]any{
			"protocol":      dstVolume.MountProtocol.Protocol,
			"cifs_auth_ids": dstVolume.MountProtocol.CifsAuthIDs,
		}},
		"snapshot_policy": []any{map[string]any{
			"max_scheduled_count": dstVolume.SnapshotPolicy.MaxScheduledCount,
			"reserve_percent":     dstVolume.SnapshotPolicy.ReservePercent,
			"schedule":            flattenNhncloudNasStorageVolumeMirrorDstVolumeSnapshotPolicySchedule(dstVolume.SnapshotPolicy.Schedule),
		}},
	}}
}

func flattenNhncloudNasStorageVolumeMirrorDstVolumeSnapshotPolicySchedule(schedule *volumes.Schedule) []any {
	if schedule == nil {
		return nil
	}

	return []any{map[string]any{
		"time":        schedule.Time,
		"time_offset": schedule.TimeOffset,
		"weekdays":    schedule.Weekdays,
	}}
}

func getNasStorageVolumeMirrorV1SrcRegion(d *schema.ResourceData, config *Config) string {
	srcRegion := d.Get("src_region").(string)
	if srcRegion != "" {
		return srcRegion
	}
	return config.Region
}
