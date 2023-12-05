# Resource: nhncloud_lb_member_v2

## Example Usage

```
resource "nhncloud_lb_member_v2" "tf_member_01"{
  pool_id = nhncloud_lb_pool_v2.tf_pool_01.id
  subnet_id = data.nhncloud_networking_vpcsubnet_v2.default_subnet.id
  address = nhncloud_compute_instance_v2.tf_instance_01.access_ip_v4
  protocol_port = 8080
  weight = 4
  admin_state_up = true
}
```

## Argument Reference

* `pool_id` - (Required) The pool ID to which the member to create belongs.
* `subnet_id` - (Required) The subnet ID of the member to create.
* `address` - (Required) The IP address of the member to receive traffic from the load balancer.
* `protocol_port` - (Required) The port of the member to receive traffic.
* `weight` - (Optional) The weight of traffic to receive from the pool <br>The higher the weight, the more traffic you receive.
* `admin_state_up` - (Optional) Administrator control status.

## Attribute Reference

The following attributes are exported:

* `id` - The unique ID for the member.
* `weight` - See Argument Reference above.
* `admin_state_up` - See Argument Reference above.
* `subnet_id` - See Argument Reference above.
* `pool_id` - See Argument Reference above.
* `address` - See Argument Reference above.
* `protocol_port` - See Argument Reference above.
