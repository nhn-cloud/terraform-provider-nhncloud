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

* `region` - (Optional) Region name that VPC to query belongs to.
* `tenant_id` - (Optional) Tenant ID that VPC to query belongs to.
* `id` - (Optional) VPC ID to query.
* `name` - (Optional) VPC name to query.

## Attribute Reference

`id` is set to the ID of the found network. In addition, the following attributes
are exported:

* `name` - See Argument Reference above.
* `region` - See Argument Reference above.
