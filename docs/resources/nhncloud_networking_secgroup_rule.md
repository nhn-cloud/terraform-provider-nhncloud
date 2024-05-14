# nhncloud_networking_secgroup_rule_v2

## Example Usage

```
resource "nhncloud_networking_secgroup_rule_v2" "resource-sg-rule-01" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 22
  port_range_max    = 22
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = data.nhncloud_networking_secgroup_v2.sg-01.id
}

###################### Data Sources ######################

data "nhncloud_networking_secgroup_v2" "sg-01" {
  name = "sg-01"
}
```

## Argument Reference

The following arguments are supported:

* `security_group_id` - (Required) The security group ID to which the security rule belongs.
* `direction` - (Required) The direction of the packet to which the security rule applies. `ingress` or `egress`.
* `ethertype` - (Optional) Set to `IPv4`. Specified as `IPv4` if omitted.
* `protocol` - (Optional) The protocol name of the security rule. Applies to all protocols if omitted.
* `port_range_min` - (Optional) The minimum port range of the security rule.
* `port_range_max` - (Optional) The maximum port range of the security rule.
* `remote_ip_prefix` - (Optional) The destination IP prefix of the security rule.
* `remote_group_id` - (Optional) The remote security group ID to which the security rule belongs.
* `description` - (Optional) The description for the security rule.

## Attribute Reference

The following attributes are exported:

* `security_group_id` - See Argument Reference above.
* `direction` - See Argument Reference above.
* `ethertype` - See Argument Reference above.
* `protocol` - See Argument Reference above.
* `port_range_min` - See Argument Reference above.
* `port_range_max` - See Argument Reference above.
* `remote_ip_prefix` - See Argument Reference above.
* `remote_group_id` - See Argument Reference above.
* `description` - See Argument Reference above.