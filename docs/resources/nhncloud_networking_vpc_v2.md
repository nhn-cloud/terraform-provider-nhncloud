# Resource: nhncloud_networking_vpc_v2

## Example Usage

```
resource "nhncloud_networking_vpc_v2" "resource-vpc-01" {
  name = "tf-vpc-01"
  cidrv4 = "10.0.0.0/8"
}
```

## Argument Reference

* `name` - (Required) VPC name.
* `cidrv4` - (Required) VPC IP range.
* `region` - (Optional) VPC region name.
* `tenant_id` - (Optional) VPC tenant ID.

## Attribute Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `shared` - Whether to share VPC.
* `tenant_id` - See Argument Reference above.