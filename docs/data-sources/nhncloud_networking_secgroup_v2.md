# Data Source: nhncloud_networking_secgroup_v2

## Example Usage

```
data "nhncloud_networking_secgroup_v2" "default_sg" {
  name = "default"
}
```

## Argument Reference

* `region` - (Optional) The region name to which the security group to query belongs.
* `tenant_id` - (Optional) The tenant ID to which the security group to query belongs.
* `name` - (Optional) The security group name to query.

## Attribute Reference

`id` is set to the ID of the found security group. In addition, the following attributes are exported:

* `region` - See Argument Reference above.
* `tenant_id` - See Argument Reference above.
* `name` - See Argument Reference above.
