# Resource: nhncloud_networking_internet_gateway_v2

## Example Usage

```
data "nhncloud_networking_network_v2" "ext" {
  external = true
}

resource "nhncloud_networking_internet_gateway_v2" "gw" {
  name = "tf-igw-01"
  external_network_id = data.nhncloud_networking_network_v2.ext.id
}
```

## Argument Reference

* `name` - (Required) The name of the internet gateway. Changing this creates a new internet gateway.
* `external_network_id` - (Required) The ID of the external network the internet gateway connects to. Changing this creates a new internet gateway.
* `region` - (Optional) The region in which to create the internet gateway. If omitted, the region configured on the provider is used. Changing this creates a new internet gateway.

## Attribute Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `external_network_id` - See Argument Reference above.
* `id` - The ID of the internet gateway.
* `routingtable_id` - The ID of the routing table the internet gateway is attached to, when one is connected.
* `state` - The status of the internet gateway. A newly created gateway stays `unavailable` until it is attached to a routing table.
* `create_time` - The time the internet gateway was created (UTC).
* `tenant_id` - The ID of the tenant that owns the internet gateway.
* `migrate_status` - The processing status while the internet gateway is being moved to another gateway server for maintenance.
* `migrate_error` - The error message when an error occurs while the internet gateway is being moved for maintenance.
