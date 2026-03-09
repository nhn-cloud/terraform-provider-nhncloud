# Resource: nhncloud_nas_storage_volume_v1

## Example Usage

```hcl
resource "nhncloud_nas_storage_volume_v1" "tf_nas_volume_01" {
  region = "KR1"
  name = "tf_nas_volume_01"
  description = "create nas volume by terraform"
  size_gb = 300

  acl = ["10.10.10.0/24"]

  encryption {
    enabled = trueㅣ
  }

  mount_protocol {
    protocol = "cifs"
    cifs_auth_ids = ["auth_id"]
  }

  snapshot_policy {
    max_scheduled_count = 3
    reserve_percent = 10

    schedule {
      time = "00:00"
      time_offset = "+09:00"
      weekdays = [1, 3, 5]
    }
  }
}
```

## Argument Reference

* `region` - (Optional) The region of the NAS storage to create.<br>The default is the region configured in the provider.
* `name` - (Required) The name of the NAS storage to create.
* `description` - (Optional) The description of the NAS storage.
* `size_gb` - (Required) The size(GB) of the NAS storage to create. NAS storage can be set from a minimum of 300GB to a maximum of 10,000GB, in 100GB increments.
* `acl` - (Optional) The list of the IPs or CIDR blocks that allow read and write permissions.
* `encryption` - (Optional) Encryption settings object when creating the NAS storage.
* `encryption.enabled` - (Optional) Whether to enable encryption settings. After the encryption keystore is set up, setting its field to `true` enables encryption.
* `mount_protocol` - (Required) Protocol settings object when creating the NAS storage.
* `mount_protocol.protocol` - (Required) Specifying protocols when mounting NAS storage. One among `nfs` and `cifs`
* `mount_protocol.cifs_auth_ids` - (Optional) The list of CIFS Authentication IDS. No input required for NFS protocol selection.
* `snapshot_policy` - (Optional) Snapshot Settings object when creating the NAS storage.
* `snapshot_policy.max_scheduled_count` - (Optional) The maximum number of snapshots that can be saved. You can set a maximum of 30, and the first automatically created snaphot will be deleted when the maximum number of saves is reached,
* `snapshot_policy.reserve_percent` - (Optional) Snapshot capacity ratio. The default is 5.
* `snapshot_policy.schedule` - (Optional) Snapshot auto-create objects. If `null`, snapshot auto-creation will not be configured.
* `snapshot_policy.schedule.time` - (Required) Snapshot auto-create time.
* `snapshot_policy.schedule.time_offset` - (Required) Time zone for snaphost auto-create.
* `snapshot_policy.schedule.weekdays` - (Required) Days of the week that snapshots are automatically created.
An empty list means every day, and the days of the week are specified as a list of numbers from 0 (Sunday) to 6 (Saturday).

## Attribute Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `id` - The unique ID for the NAS storage.
* `name` - See Argument Refernce above.
* `description` - See Argument Refernce above.
* `size_gb` - See Argument Refernce above.
* `project_id` - The project ID to which the NAS storage belongs.
* `tenant_id` - The tenant ID to which the NAS storage belongs
* `acl` - See Argument Refernce above.
* `encryption` - See Argument Refernce above.
* `encryption.enabled` - See Argument Refernce above.
* `encryption.keys` - NAS Storage encryption keys information.
* `encryption.keys.key_id` - ID for the key used for encryption.
* `encryption.keys.key_store_id` - ID for the key store used for encryption.
* `mount_protocol` - See Argument Refernce above.
* `mount_protocol.protocol` - See Argument Refernce above.
* `mount_protocol.cifs_auth_ids` - See Argument Refernce above.
* `snapshot_policy` - See Argument Refernce above.
* `snapshot_policy.max_scheduled_count` - See Argument Refernce above.
* `snapshot_policy.reserve_percent` - See Argument Refernce above.
* `snapshot_policy.schedule` - See Argument Refernce above.
* `snapshot_policy.schedule.time` - See Argument Refernce above.
* `snapshot_policy.schedule.time_offset` - See Argument Refernce above.
* `snapshot_policy.schedule.weekdays` - See Argument Refernce above.
* `created_at` - The date the NAS storage was created
* `updated_at` - The date the NAS storage was last updated.
