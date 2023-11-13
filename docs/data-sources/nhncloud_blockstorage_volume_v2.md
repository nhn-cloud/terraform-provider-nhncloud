# Data Source: nhncloud_blockstorage_volume_v2

## Example Usage

```
data "nhncloud_blockstorage_volume_v2" "volume_00" {
  name = "ssd_volume1"
  status = "available"
}
```

## Argument Reference

* `region` - (Optional) Region name that block storage to query belongs to.
* `name` - (Optional) Name of block storage to query.
* `status` - (Optional) Status of block storage to query.
* `metadata` - (Optional) Metadata related to block storage to query.

## Attribute Reference

`id` is set to the ID of the found volume. In addition, the following attributes
are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `status` - See Argument Reference above.
* `metadata` - See Argument Reference above.
* `volume_type` - The type of the volume.
* `bootable` - Indicates if the volume is bootable.
* `size` - The size of the volume in GBs.