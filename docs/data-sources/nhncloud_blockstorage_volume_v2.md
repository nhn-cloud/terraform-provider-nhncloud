# Data Source: nhncloud_blockstorage_volume_v2

## Example Usage

```
data "nhncloud_blockstorage_volume_v2" "volume_00" {
  name = "ssd_volume1"
  status = "available"
}
```

## Argument Reference

* `region` - (Optional) The region name in which the block storage to query exists.
* `name` - (Optional) The name of the block storage to query.
* `status` - (Optional) The status of the block storage to query.
* `metadata` - (Optional) The metadata related to the block storage to query.

## Attribute Reference

`id` is set to the ID of the found volume. In addition, the following attributes
are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `status` - See Argument Reference above.
* `metadata` - See Argument Reference above.
* `volume_type` - The type of the volume.
* `bootable` - Indicates if the volume is bootable.
* `size` - The size of the volume in GB.