# Resource: nhncloud_nas_storage_volume_interface_v1

## Exmpale Usage

```hcl
# Create NAS Storage Volume
resource "nhncloud_nas_storage_volume_v1" "tf_nas_volume_01" {
  ...
  mount_protocol {
    protocol = "nfs"
  }
  ...
}

# Create NAS Storage Volume Interface
resource "nhncloud_nas_storage_volume_interface_v1" "tf_nas_interface_01" {
  region = "KR1"
  volume_id = nhncloud_nas_storage_volume_v1.tf_nas_volume_01.id
  subnet_id = data.nhncloud_networking_vpcsubnet_v2.default_subnet.id
}

# Create Instance with user script which connects to the interface after booting 
resource "nhncloud_compute_instance_v2" "tf_instance_01" {
  ...
  user_data = <<-EOT
    #! /bin/bash
    sudo service rpcbind start
    sudo mount -t nfs ${nhncloud_nas_storage_volume_interface_v1.tf_nas_interface_01.path} /mount_point
    EOT
  ...
}
```

## Argument Reference

* `region` - (Optional) The region of the NAS storage interface to create.<br>The default is the region configured in the provider.
* `volume_id` - (Required) The NAS storage ID to which the interface is connected
* `subnet_id` - (Required) The subnet ID associated with the NAS storage

## Attribute Reference

The flowwing attributes are exported:

* `region` - See Argument Reference above.
* `volume_id` - See Argument Reference above.
* `id` - The unique ID for the interface.
* `path` - The connection path that the instance will use when mounting
* `subnet_id` - See Argument Reference above.
* `tenant_id` - The tenant ID to which the interface belongs.
