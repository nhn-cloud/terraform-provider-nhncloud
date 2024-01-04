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

* `name` - (Required) The name of the routing table.
* `vpc_id` - (Optional) The VPC ID of the routing table.
* `distributed` - (Optional) The routing method for the routing table. (default: `true`)
  * `true`: decentralized
  * `false`: centralized

## Attribute Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `shared` - Whether to share the routing table.
* `tenant_id` - See Argument Reference above.