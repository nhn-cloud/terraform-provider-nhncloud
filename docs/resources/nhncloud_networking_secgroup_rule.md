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

* `security_group_id` - (Required) Security group ID containing the security rule.
* `direction` - (Required) Direction of packet to which the security rule is applied `ingress`, `egress`.
* `ethertype` - (Optional) Set to `IPv4`. Specified as `IPv4` if omitted.
* `protocol` - (Optional) Protocol name of the security rule. Applied to all protocols if omitted.
* `port_range_min` - (Optional) Minimum port range of the security rule.
* `port_range_max` - (Optional) Maximum port range of the security rule.
* `remote_ip_prefix` - (Optional) Destination IP prefix of the security rule.
* `remote_group_id` - (Optional) Remote security group ID of the security rule.
* `description` - (Optional) Security rule description.

## Attributes Reference

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