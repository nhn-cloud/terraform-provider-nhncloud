# Resource: nhncloud_compute_instance_v2

## Example Usage

### Create u2 Instance
```
resource "nhncloud_compute_instance_v2" "tf_instance_01"{
  name = "tf_instance_01"
  region    = "KR1"
  key_pair  = "terraform-keypair"
  image_id = data.nhncloud_images_image_v2.ubuntu_2004_20201222.id
  flavor_id = data.nhncloud_compute_flavor_v2.u2c2m4.id
  security_groups = ["default"]
  availability_zone = "kr-pub-a"

  network {
    name = data.nhncloud_networking_vpc_v2.default_network.name
    uuid = data.nhncloud_networking_vpc_v2.default_network.id
  }

  block_device {
    uuid = data.nhncloud_images_image_v2.ubuntu_2004_20201222.id
    source_type = "image"
    destination_type = "local"
    boot_index = 0
    delete_on_termination = true
    volume_size = 30
  }
}
```


### Flavors other than u2
#### Create an instance with network and block storage added
```
resource "nhncloud_compute_instance_v2" "tf_instance_02" {
  name      = "tf_instance_02"
  region    = "KR1"
  key_pair  = "terraform-keypair"
  flavor_id = data.nhncloud_compute_flavor_v2.m2c1m2.id
  security_groups = ["default","web"]

  network {
    name = data.nhncloud_networking_vpc_v2.default_network.name
    uuid = data.nhncloud_networking_vpc_v2.default_network.id
  }

  network {
    port = nhncloud_networking_port_v2.port_1.id
  }

  block_device {
    uuid                  = data.nhncloud_images_image_v2.ubuntu_2004_20201222.id
    source_type           = "image"
    destination_type      = "volume"
    boot_index            = 0
    volume_size           = 20
    delete_on_termination = true
  }

  block_device {
    source_type           = "blank"
    destination_type      = "volume"
    boot_index            = 1
    volume_size           = 20
    delete_on_termination = true
  }
}
```

## Argument Reference

* `region` - (Optional) The region of the instance to create<br>The default is the region configured in the provider.
* `flavor_name` - (Optional) The flavor name of the instance to create<br>Required if flavor_id is.
* `name` - (Required) The name of the instance to create.
* `flavor_id` - (Optional) The flavor ID of the instance to create<br>Required if flavor_name is empty.
* `image_name` - (Optional) The image name to use for creating an instance<br>Required if image_id is empty<br>Available only when the flavor is U2.
* `image_id` - (Optional) The image ID to use for creating an instance<br>Required if image_name is empty<br>Available only when the flavor is U2.
* `key_pair` - (Optional) The key pair name to use for accessing the instance<br>You can create a new key pair from **Compute > Instance > Key Pairs** on the NHN Cloud console,<br>or register an existing key pair<br>See `User Guide > Compute > Instance > Console User Guide` for more.
* `availability_zone` - (Optional) The availability zone of an instance to create.
* `network` - (Optional) VPC network information to be attached to an instance to create.<br>Go to **Network > VPC > Management**  on the console, select the VPC to be attached, and check the network name and UUID at the bottom.
* `network.name` - (Optional) The name of the VPC network <br>One among network.name, network.uuid, and network.port must be specified.
* `network.uuid` - (Optional) The ID of the VPC.
* `network.port` - (Optional) The ID of a port to be attached to VPC.
* `security_groups` - (Optional) List of the security group names for instance <br>Select a security group from **Network > VPC > Security Groups** on the console, and check detailed information at the bottom of the page.
* `user_data` - (Optional) 	The script to be executed after instance booting and its configuration<br>Base64-encoded string, which allows up to 65535 bytes.
* `block_device` - (Optional) Information object of the image or block storage to be used for an instance.
* `block_device.uuid` - (Optional) The ID of the original block storage <br>The block storage must be a bootable source if used as the root block storage. Volumes or snapshots which cannot be used to create images, such as those with WAF, MS-SQL images as the source, cannot be used.<br> The original other than `image` must have the same availability zone for the instance to create.
* `block_device.source_type` - (Optional) The type of the original block storage to create<br>`image`: Use an image to create a block storage<br>`volume`: Use the existing block storage, with the destination_type set to volume<br>`snapshot`: Use a snapshot to create block storage, with the destination_type set to volume.
* `block_device.destination_type` - (Optional) Requires different settings depending on the location of instance’s block storage or flavor<br>`local`: For U2 flavor<br>`volume`: For flavors other than U2.
* `block_device.boot_index` - (Optional) The order to boot the specified block storage<br>- If so, root block storage<br>- If not, additional block storage<br>The higher the number, the lower the booting priority.
* `block_device.volume_size` - (Optional) The block storage size for the instance to create<br>Available from 20 GB to 2,000 GB (required if the flavor is U2)<br>Since each flavor allows different volume size, see `User Guide > Compute > Instance Console User Guide`.
* `block_device.delete_on_termination` - (Optional) `true`: When deleting an instance, delete a block device<br>`false`: When deleting an instance, do not delete a block.
* `block_device.nhn_encryption` - (Optional) About block storage encryption.
* `block_device.nhn_encryption.skm_appkey` - (Required) The appKeys for Secure Key Manager products.
* `block_device.nhn_encryption.skm_key_id` - (Required) The key ID in Secure Key Manager.

## Attribute Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `access_ip_v4` - The first detected Fixed IPv4 address.
* `access_ip_v6` - The first detected Fixed IPv6 address.
* `security_groups` - See Argument Reference above.
* `flavor_id` - See Argument Reference above.
* `flavor_name` - See Argument Reference above.
* `network.uuid` - See Argument Reference above.
* `network.name` - See Argument Reference above.
* `network.port` - See Argument Reference above.
* `network.fixed_ip_v4` - The Fixed IPv4 address of the Instance on that network.
* `network.fixed_ip_v6` - The Fixed IPv6 address of the Instance on that network.
* `network.mac` - The MAC address of the NIC on that network.
* `created` - The time when the instance was created.
* `updated` - The time when the instance was last updated.