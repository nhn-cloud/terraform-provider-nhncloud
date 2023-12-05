# Resource: nhncloud_networking_port_v2

## Example Usage

```
resource "nhncloud_networking_port_v2" "port_1" {
  name = "tf_port_1"
  network_id = data.nhncloud_networking_vpc_v2.default_network.id
  admin_state_up = "true"
}
```

## Argument Reference

* `name` - (Required) The port name to create.
* `description` - (Optional) The port description.
* `network_id` - (Required) The ID of the VPC network to create a port.
* `tenant_id` - (Optional) The tenant ID of the port to create.
* `device_id` - (Optional) The device ID to which the created port will be connected.
* `fixed_ip` - (Optional) Setting information of the fixed IP of a port to create<br>Must not include the `no_fixed_ip` attribute.
* `fixed_ip.subent_id` - (Required) The Subnet ID of a fixed IP.
* `fixed_ip.ip_address` - (Optional) The address of fixed IP to configure.
* `no_fixed_ip` - (Optional) `true`: Port without fixed IP<br>Must not include the `fixed_ip` attribute.
* `admin_state_up` - (Optional) Administrator control status<br> `true`: Running<br>`false`: Suspended.

## Attribute Reference

The following attributes are exported:

* `description` - See Argument Reference above.
* `admin_state_up` - See Argument Reference above.
* `tenant_id` - See Argument Reference above.
* `device_id` - See Argument Reference above.
* `fixed_ip` - See Argument Reference above.
* `all_fixed_ips` - The collection of Fixed IP addresses on the port in the order returned by the Network v2 API.