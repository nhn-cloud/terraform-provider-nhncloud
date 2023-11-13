# Data Source: nhncloud_blockstorage_snapshot_v2

## Example Usage

```
data "nhncloud_blockstorage_snapshot_v2" "my_snapshot" {
  name = "my-snapshot"
  volume_id = data.nhncloud_blockstorage_volume_v2.volume_00.id
  status = "available"
  most_recent = true
}
```

## Argument Reference

* `region` - (Optional) Region name that snapshot to query belongs to.
* `name` - (Optional) Name of snapshot to query.
* `volume_id` - (Optional) ID of original block storage of snapshot to query.
* `status` - (Optional) Status of snapshot to query.
* `most_recent` - (Optional) 
  * `true`: Select the most recently created snapshot from the queried snapshot list.
  * `false`: Select snapshots in the queried order.

## Attribute Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `status` - See Argument Reference above.
* `volume_id` - See Argument Reference above.
* `size` - The size of the snapshot.
* `metadata` - The snapshot's metadata.