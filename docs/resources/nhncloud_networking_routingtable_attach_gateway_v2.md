# Resource: nhncloud_networking_routingtable_v2

## Example Usage

```
resource "nhncloud_networking_vpc_v2" "resource-vpc-01" {
  ...
}

resource "nhncloud_networking_routingtable_v2" "resource-rt-01" {
  name = "resource-rt-01"
  vpc_id = nhncloud_networking_vpc_v2.resource-vpc-01.id
  distributed = false
}
```

## Argument Reference

* `routingtable_id` - (Required) Routing table ID to modify.
* `gateway_id` - (Required) Internet gateway ID to be associated with routing table.
  In the console, select the Internet gateway you want to use from the **Network > Internet Gateway** menu, 
  and you can see the ID of the gateway in the details screen below.

## Attribute Reference

`id` is set to the ID of the found attachment ID of the gateway and routingtable. In addition, the following attributes are exported:

* `routingtable_id` - See Argument Reference above.
* `gateway_id` - See Argument Reference above.