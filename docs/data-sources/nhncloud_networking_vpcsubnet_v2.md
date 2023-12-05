# Data Source: nhncloud_networking_vpcsubnet_v2

## Example Usage

```
data "nhncloud_networking_vpcsubnet_v2" "default_subnet" {
  region = "KR1"
  tenant_id = "ba3be1254ab141bcaef674e74630a31f"
  id = "05f6fdc3-641f-48df-b986-773b6489654f"
  name = "Default Network"
  shared = true
}
```

## Argument Reference

* `region` - (Optional) The region name in which the subnet to query exists.
* `tenant_id` - (Optional) The tenant ID to which the subnet to query belongs.
* `id` - (Optional) The subnet ID to query.
* `name` - (Optional) The subnet name to query.
* `shared` - (Optional) Whether to share the subnet to query.

## Attribute Reference

`id` is set to the ID of the found subnet. In addition, the following attributes are exported:

* `region` - See Argument Reference above.
