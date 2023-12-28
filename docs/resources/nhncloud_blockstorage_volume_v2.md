# Resource: nhncloud_blockstorage_volume_v2

## Example Usage

```
# Create HDD-type Empty Block Storage
resource "nhncloud_blockstorage_volume_v2" "volume_01" {
  name = "tf_volume_01"
  size = 10
  availability_zone = "kr-pub-a"
  volume_type = "General HDD"
}

# Create SSD-type Empty Block Storage
resource "nhncloud_blockstorage_volume_v2" "volume_02" {
  name = "tf_volume_02"
  size = 10
  availability_zone = "kr-pub-b"
  volume_type = "General SSD"
}

# Create Block Storage with Snapshot
resource "nhncloud_blockstorage_volume_v2" "volume_03" {
  name = "tf_volume_03"
  description = "terraform create volume with snapshot test"
  snapshot_id = data.nhncloud_blockstorage_snapshot_v2.snapshot_01.id
  size = 30
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region of the block storage to create.<br>The default is the region configured in the provider.
* `name` - (Optional) The name of the block storage to create.
* `description` - (Optional) The description of the block storage.
* `size` - (Required) The size of the block storage to create (GB).
* `snapshot_id` - (Optional) The snapshot ID from which to create the block storage.
* `availability_zone` - (Optional) The availability zone of the block storage to create. If the value does not exist, random availability zone is used. <br>To check availability_zone, go to `Storage > Block Storage > Management` on the console and click **Create Block Storage**.
* `volume_type` - (Optional) The type of the block storage. 
  <br> `General HDD`: HDD block storage (default). 
  <br>`General SSD`: SSD block storage.
  <br>`Encrypted HDD`: Encrypted HDD block storage.
  <br>`Encrypted SSD`: Encrypted SSD block storage.
* `nhn_encryption` - (Optional) Block storage encryption.
* `nhn_encryption.skm_appkey` - (Required) The appKeys for Secure Key Manager products.
* `nhn_encryption.skm_key_id` - (Required) The key ID in Secure Key Manager.

## Attribute Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `size` - See Argument Reference above.
* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `availability_zone` - See Argument Reference above.
* `snapshot_id` - See Argument Reference above.
* `volume_type` - See Argument Reference above.
* `attachment` - If a volume is attached to an instance, this attribute will display the Attachment ID, Instance ID, and the Device as the Instance sees it.