# Resource: nhncloud_compute_volume_attach_v2

## Example Usage

```
# Create Instance
resource "nhncloud_compute_instance_v2" "tf_instance_01" {
  ...
}

# Create Block Storage
resource "nhncloud_blockstorage_volume_v2" "volume_01" {
  ...
}

# Attach Block Storage
resource "nhncloud_compute_volume_attach_v2" "volume_to_instance"{
  instance_id = nhncloud_compute_instance_v2.tf_instance_02.id
  volume_id = nhncloud_blockstorage_volume_v2.volume_01.id
  vendor_options {
    ignore_volume_confirmation = true
  }
}
```

## Argument Reference

* `instance_id` - (Required) Target instance to attach the block storage.
* `volume_id` - (Required) UUID of block storage to be attached.

## Attribute Reference

In addition to the above, the following attributes are exported:

* `data` - This is a map of key/value pairs that contain the connection
  information. You will want to pass this information to a provisioner
  script to finalize the connection. See below for more information.
* `driver_volume_type` - The storage driver that the volume is based on.
* `mount_point_base` - A mount point base name for shared storage.