# Data Source: nhncloud_images_image_v2

## Example Usage

```
data "nhncloud_images_image_v2" "ubuntu_2004_20201222" {
  name = "Ubuntu Server 20.04.1 LTS (2020.12.22)"
  most_recent = true
}

# Query the oldest image from images with the same name
data "nhncloud_images_image_v2" "windows2016_20200218" {
  name = "Windows 2019 STD with MS-SQL 2019 Standard (2020.12.22) KO"
  sort_key = "created_at"
  sort_direction = "asc"
  owner = "c289b99209ca4e189095cdecebbd092d"
  tag = "_AVAILABLE_"
}
```

## Argument Reference

* `name` - (Optional) The name of the image to query. To check the image name, go to **Compute > Instance** on the NHN Cloud console and click **Create Instance**. The information can be found on the list of images provided by NHN Cloud. The image name must be entered as **<Image Description>** which appears on the NHN Cloud console. If the Language item exists, follow the format **"<Image Description> <Language>"** as shown in the above.
* `size_min` - (Optional) The minimum size of the image to query (bytes).
* `size_max` - (Optional) The maximum size of the image to query (bytes).
* `properties` - (Optional) The attributes of the image to query. Images that match all attributes are queried.
* `sort_key` - (Optional) Sort the list of images queried by particular attributes. The default is `name`.
* `sort_direction` - (Optional) The sorting order of the list of queried images. 
  * `asc`: Ascending order (Default).
  * `desc`: Descending order.
* `owner` - (Optional) The ID of the tenant which includes the image to query.
* `tag` - (Optional) Search images with a particular tag.
* `visibility` - (Optional) The visibility of the image to query <br>Select only one among public, private, and shared. If omitted, the list with all types of images is returned.
* `most_recent` - (Optional)
  * `true`: Select the most recently created image from the list of queried images 
  * `false`: Select images in the order they were queried.
* `member_status` - (Optional) The status of the image member to query. One among `accepted`,`pending`, `rejected`, and `all`.

## Attribute Reference

`id` is set to the ID of the found image. In addition, the following attributes
are exported:

* `checksum` - The checksum of the data associated with the image.
* `container_format`: The format of the image's container.
* `disk_format`: The format of the image's disk.
* `file` - The trailing path after the glance endpoint that represent the location of the image or the path to retrieve it.
* `metadata` - The metadata associated with the image. Image metadata allow for meaningfully define the image properties and tags. See https://docs.openstack.org/glance/latest/user/metadefs-concepts.html.
* `min_disk_gb` - The minimum amount of disk space required to use image.
* `min_ram_mb` - The minimum amount of ram required to use image.
* `properties` - The freeform information about the image.
* `protected` - Whether or not the image is protected.
* `schema` - The path to the JSON-schema that represent the image or image
* `size_bytes` - The size of the image (in bytes).
* `created_at` - The date the image was created.
* `updated_at` - The date the image was last updated.
