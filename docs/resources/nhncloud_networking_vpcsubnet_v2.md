# Resource: nhncloud_networking_vpcsubnet_v2

## Example Usage

```
resource "nhncloud_networking_vpcsubnet_v2" "resource-vpcsubnet-01" {
  name      = "tf-vpcsubnet-01"
  vpc_id    = "def56b5e-0f1d-4a31-8005-4d716127f177"
  cidr      = "10.10.10.0/24"
  routingtable_id = "c3ed678d-de8b-4bf7-abea-b7c1118f0828"
}
```

## Argument Reference

* `vpc_id` - (Requried) VPC ID to which subnet is assigned.
* `cidr` - (Requried) IP range of subnet.
* `name` - (Requried) Name of subnet.
* `region` - (Optional) Name of region to which subnet is assigned.
* `tenant_id` - (Optional) Tenant ID to which subnet is assigned.
* `routingtable_id` - (Optional) Routing table ID.

## Attribute Reference

The following attributes are exported:

* `vpc_id` - See Argument Reference above.
* `cidr` - See Argument Reference above.
* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `shared` - Whether to share subnet.
* `tenant_id` - See Argument Reference above.
* `gateway_ip` - Gateway IP of subnet.