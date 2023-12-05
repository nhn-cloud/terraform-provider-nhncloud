# Resource: nhncloud_lb_loadbalancer_v2

## Example Usage

```
resource "nhncloud_lb_loadbalancer_v2" "tf_loadbalancer_01"{
  name = "tf_loadbalancer_01"
  description = "create loadbalancer by terraform."
  vip_subnet_id = data.nhncloud_networking_vpcsubnet_v2.default_subnet.id
  vip_address = "192.168.0.10"  
  admin_state_up = true
}
```

## Argument Reference

* `name` - (Optional) The name of the load balancer.
* `description` - (Optional) The description of the load balancer.
* `tenant_id` - (Optional) The tenant ID for which the load balancer is to be created.
* `vip_subnet_id` - (Required) The subnet UUID to be used by the load balancer.
* `vip_address` - (Optional) The IP address of the load balancer.
* `security_group_ids` - (Optional) The list of security group IDs to be applied for the load balancer.<br>**Security groups must be specified by ID, not by name**.
* `admin_state_up` - (Optional) Administrator control status.

## Attribute Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `vip_subnet_id` - See Argument Reference above.
* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `tenant_id` - See Argument Reference above.
* `vip_address` - See Argument Reference above.
* `admin_state_up` - See Argument Reference above.
* `security_group_ids` - See Argument Reference above.
* `vip_port_id` - The Port ID of the Load Balancer IP.
