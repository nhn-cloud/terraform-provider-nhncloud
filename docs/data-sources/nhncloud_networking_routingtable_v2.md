# Data Source: nhncloud_networking_routingtable_v2

## Example Usage

```
data "nhncloud_networking_routingtable_v2" "default_rt" {
  id = "bf15f6f6-1339-4057-a7fe-5811d39bab18"
}
```

## Argument Reference

* `tenant_id` - (Optional) The tenant ID to which the routing table to query belongs.
* `id` - (Optional) The routing table ID to query.
* `name` - (Optional) The routing table name to query.

## Attribute Reference

`id` is set to the ID of the found routing table. In addition, the following attributes are exported:

* `tenant_id` - See Argument Reference above.
* `id` - See Argument Reference above.
* `name` - See Argument Reference above.