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
#### Create instance with network and block storage added
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

* `region` - (Optional) Region of instance to create<br>The default is the region configured in provider.
* `flavor_name` - (Optional) Flavor name of instance to create<br>Required if flavor_id is.
* `name` - (Required) Name of instance to.
* `flavor_id` - (Optional) Flavor ID of instance to create<br>Required if flavor_name is.
* `image_name` - (Optional) Image name to use for creating an instance<br>Required if image_id is empty<br>Available only when the flavor is.
* `image_id` - (Optional) Image ID to use for creating an instance<br>Required if image_name is empty<br>Available only when the flavor is.
* `key_pair` - (Optional) Key pair name to use for accessing the instance<br>You can create a new key pair from **Compute > Instance > Key Pairs** on NHN Cloud console,<br>or register an existing key pair<br>See `User Guide > Compute > Instance > Console User Guide` for more.
* `availability_zone` - (Optional) Availability zone of an instance to.
* `network` - (Optional) VPC network information to be attached to an instance to create.<br>Go to **Network > VPC > Management**  on the console, select VPC to be attached, and check the network name and UUID at the bottom.
* `network.name` - (Optional) Name of VPC network <br>One among network.name, network.uuid, and network.port must be specified.
* `network.uuid` - (Optional) ID of VPC.
* `network.port` - (Optional) ID of a port to be attached to VPC.
* `security_groups` - (Optional) List of the security group names for instance <br>Select a security group from **Network > VPC > Security Groups** on the console, and check detailed information at the bottom of the page.
* `user_data` - (Optional) 	Script to be executed after instance booting and its configuration<br>Base64-encoded string, which allows up to 65535 bytes<br.
* `block_device` - (Optional) Information object of image or block storage to be used for an.
* `block_device.uuid` - (Optional) ID of original block storage <br>The block storage must be a bootable source if used as the root block storage. Volumes or snapshots which cannot be used to create images, such as those with WAF, MS-SQL images as the source, cannot be used.<br> The original other than `image` must have the same availability zone for the instance to create.
* `block_device.source_type` - (Optional) Type of original block storage to create<br>`image`: Use an image to create a block storage<br>`volume`: Use the existing block storage, with the destination_type set to volume<br>`snapshot`: Use a snapshot to create a block storage, with the destination_type set to.
* `block_device.destination_type` - (Optional) Requires different settings depending on the location of instance’s block storage or flavor<br>`local`: For U2 flavor<br>`volume`: For flavors other than.
* `block_device.boot_index` - (Optional) Order to boot the specified block storage<br>- If , root block storage<br>- If not, additional block storage<br>The higher the number, the lower the booting priority<br>
* `block_device.volume_size` - (Optional) Block storage size for instance to create<br>Available from 20GB to 2,000GB (required if the flavor is U2)<br>Since each flavor allows different volume size, see `User Guide > Compute > Instance Console User Guide.
* `block_device.delete_on_termination` - (Optional) `true`: When deleting an instance, delete a block device<br>`false`: When deleting an instance, do not delete a block.

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
* `created` - The creation time of the instance.
* `updated` - The time when the instance was last updated.