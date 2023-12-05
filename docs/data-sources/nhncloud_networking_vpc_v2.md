# Data Source: nhncloud_networking_vpc_v2

## Example Usage

```
data "nhncloud_networking_vpc_v2" "default_network" {
  region = "KR1"
  tenant_id = "ba3be1254ab141bcaef674e74630a31f"
  id = "e34fc878-89f6-4d17-a039-3830a0b78346"
  name = "Default Network"
}
```

## Argument Reference

* `region` - (Optional) The region name in which the VPC to query exists.
* `tenant_id` - (Optional) The tenant ID to which the VPC to query belongs.
* `id` - (Optional) The VPC ID to query.
* `name` - (Optional) The VPC name to query.

## Attribute Reference

`id` is set to the ID of the found network. In addition, the following attributes are exported:

* `name` - See Argument Reference above.
* `region` - See Argument Reference above.
