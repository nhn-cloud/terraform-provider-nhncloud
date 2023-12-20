# Resource: nhncloud_networking_secgroup_v2

## Example Usage

```
resource "nhncloud_networking_secgroup_v2" "resource-sg-01" {
  name      = "sg-01"
}
```

## Argument Reference

* `name` - (Required) The security group name.
* `region` - (Optional) Name of region to which security group is assigned.

## Attribute Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.