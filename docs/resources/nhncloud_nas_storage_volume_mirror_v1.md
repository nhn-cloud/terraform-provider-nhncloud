# Resource: nhncloud_nas_storage_volume_interface_v1

## Exmpale Usage

```hcl
# Create NAS Storage Volume
resource "nhncloud_nas_storage_volume_v1" "tf_nas_volume_01" {
  ...
}

# Create NAS Storage Mirror
resource "nhncloud_nas_storage_volume_mirror_v1" "tf_nas_mirror_01" {
  src_region    = "KR1"
  src_volume_id = nhncloud_nas_storage_volume_v1.tf_nas_volume_01.id
  dst_region    = "KR2"
  dst_tenant_id = "ba3be1254ab141bcaef674e74630a31f"

  dst_volume {
    name        = "tf_nas_mirror_01"
    description = "create nas mirror by terraform"
    size_gb     = 300

    acl = ["10.10.10.0/24"]

    mount_protocol {
      protocol = "nfs"
    }
  }
}
```

## Argument Reference

* `src_region` - (Optional) The region of the source NAS storage.<br>The default is the region configured in the provider.
* `src_volume_id` - (Required) Source NAS storage ID
* `dst_region` - (Required) - The region of the replication target storage.
* `dst_tenant_id` - (Required) - The tenant ID of the replication target storage.
* `dst_volume` - (Required) - The replication target storage object.
* `dst_volume.name` - (Required) The name of the NAS storage to create.
* `dst_volume.description` - (Optional) The description of the NAS storage.
* `dst_volume.size_gb` - (Required) The size(GB) of the NAS storage to create. NAS storage can be set from a minimum of 300GB to a maximum of 10,000GB, in 100GB increments.
* `dst_volume.acl` - (Optional) The list of the IPs or CIDR blocks that allow read and write permissions.
* `dst_volume.encryption` - (Optional) Encryption settings object when creating the NAS storage.
* `dst_volume.encryption.enabled` - (Optional) Whether to enable encryption settings. After the encryption keystore is set up, setting its field to `true` enables encryption.
* `dst_volume.mount_protocol` - (Required) Protocol settings object when creating the NAS storage.
* `dst_volume.mount_protocol.protocol` - (Required) Specifying protocols when mounting NAS storage. One among `nfs` and `cifs`
* `dst_volume.mount_protocol.cifs_auth_ids` - (Optional) The list of CIFS Authentication IDS. No input required for NFS protocol selection.
* `dst_volume.snapshot_policy` - (Optional) Snapshot Settings object when creating the NAS storage.
* `dst_volume.snapshot_policy.max_scheduled_count` - (Optional) The maximum number of snapshots that can be saved. You can set a maximum of 30, and the first automatically created snaphot will be deleted when the maximum number of saves is reached,
* `dst_volume.snapshot_policy.reserve_percent` - (Optional) Snapshot capacity ratio. The default is 5.
* `dst_volume.snapshot_policy.schedule` - (Optional) Snapshot auto-create objects. If `null`, snapshot auto-creation will not be configured.
* `dst_volume.snapshot_policy.schedule.time` - (Required) Snapshot auto-create time.
* `dst_volume.snapshot_policy.schedule.time_offset` - (Required) Time zone for snaphost auto-create.
* `dst_volume.snapshot_policy.schedule.weekdays` - (Required) Days of the week that snapshots are automatically created.
An empty list means every day, and the days of the week are specified as a list of numbers from 0 (Sunday) to 6 (Saturday).

## Attribute Reference

The following attributes are exported:

* `id` - The unique ID for the replication settings.
* `role` - The replication role. One among `SOURCE` and `DESTINATION`.
* `direction` - The replication direction. One among `FORWARD`(source storage -> replica storage) and `REVERSE`(replica stroage -> source storage).
* `direction_changed_at` - The replication direction changed time.
* `dst_project_id` - The project ID of the replication target storage.
* `dst_region_id` - See Argument Reference above.
* `dst_tenant_id` - See Argument Reference above.
* `dst_volume_id` - The NAS storage ID of the replication target storage.
* `dst_volume_name` - The NAS storage name of the replication target storage.
* `src_project_id` - The project ID of the source storage.
* `src_region` - See Argument Reference above.
* `src_tenant_id` - The tenant ID of the source storage.
* `src_volume_id` - See Argument Reference above.
* `src_volume_name` - The NAS storage name of the source storage.
* `created_at` - The replication created time.
* `dst_volume` - See Argument Reference above.
* `dst_volume.name` - See Argument Reference above.
* `dst_volume.description` - See Argument Reference above.
* `dst_volume.size_gb` - See Argument Reference above.
* `dst_volume.acl` - See Argument Reference above.
* `dst_volume.encryption` - See Argument Reference above.
* `dst_volume.encryption.enabled` - See Argument Reference above.
* `dst_volume.mount_protocol` - See Argument Reference above.
* `dst_volume.mount_protocol.protocol` - See Argument Reference above.
* `dst_volume.mount_protocol.cifs_auth_ids` - See Argument Reference above.
* `dst_volume.snapshot_policy` - See Argument Reference above.
* `dst_volume.snapshot_policy.max_scheduled_count` - See Argument Reference above.
* `dst_volume.snapshot_policy.reserve_percent` - See Argument Reference above.
* `dst_volume.snapshot_policy.schedule` - See Argument Reference above.
* `dst_volume.snapshot_policy.schedule.time` - See Argument Reference above.
* `dst_volume.snapshot_policy.schedule.time_offset` - See Argument Reference above.
* `dst_volume.snapshot_policy.schedule.weekdays` - See Argument Reference above.
