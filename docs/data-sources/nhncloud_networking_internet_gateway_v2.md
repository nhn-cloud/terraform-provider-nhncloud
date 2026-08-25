# Data Source: nhncloud_networking_internet_gateway_v2

## Example Usage

```
data "nhncloud_networking_internet_gateway_v2" "gw" {
  name = "tf-igw-01"
}
```

## Argument Reference

* `region` - (Optional) The region name in which the internet gateway to query exists.
* `id` - (Optional) The internet gateway ID to query.
* `name` - (Optional) The internet gateway name to query.
* `external_network_id` - (Optional) The external network ID connected to the internet gateway to query.
* `routingtable_id` - (Optional) The ID of the routing table the internet gateway to query is attached to.
* `tenant_id` - (Optional) The tenant ID to which the internet gateway to query belongs.

## Attribute Reference

`id` is set to the ID of the found internet gateway. In addition, the following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `external_network_id` - See Argument Reference above.
* `routingtable_id` - See Argument Reference above.
* `tenant_id` - See Argument Reference above.
* `state` - The status of the internet gateway.
* `create_time` - The time the internet gateway was created (UTC).
* `migrate_status` - The processing status while the internet gateway is being moved to another gateway server for maintenance.
* `migrate_error` - The error message when an error occurs while the internet gateway is being moved for maintenance.
