# Resource: nhncloud_networking_vpc_v2

## Example Usage

```
resource "nhncloud_networking_vpc_v2" "resource-vpc-01" {
  name = "tf-vpc-01"
  cidrv4 = "10.0.0.0/8"
}
```

## Argument Reference

* `name` - (Required) The name for the VPC.
* `cidrv4` - (Required) The IP range for the VPC.
* `region` - (Optional) The region name of the VPC.
* `tenant_id` - (Optional) The tenant ID of the VPC.

## Attribute Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `shared` - Whether to share VPC.
* `tenant_id` - See Argument Reference above.