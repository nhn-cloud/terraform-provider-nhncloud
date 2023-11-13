# Resource: nhncloud_networking_floatingip_associate_v2

## Example Usage

```
# Create Network Port
resource "nhncloud_networking_port_v2" "port_1" {
  ...
}

# Create Instance
resource "nhncloud_compute_instance_v2" "tf_instance_01" {
    ...
    network {
    port = nhncloud_networking_port_v2.port_1.id
  }
  ...
}

# Create Floating IP
resource "nhncloud_networking_floatingip_v2" "fip_01" {
  ...
}

# Associate Floating IP
resource "nhncloud_networking_floatingip_associate_v2" "fip_associate" {
  floating_ip = nhncloud_networking_floatingip_v2.fip_01.address
  port_id = nhncloud_networking_port_v2.port_1.id
}
```

## Argument Reference

* `floating_ip` - (Required) Floating IP to associate.
* `port_id` - (Required) UUID of port to be associated with floating IP.

## Attribute Reference

The following attributes are exported:

* `floating_ip` - See Argument Reference above.
* `port_id` - See Argument Reference above.