# Resource: nhncloud_networking_floatingip_v2

## Example Usage

```
resource "nhncloud_networking_floatingip_v2" "fip_01" {
  pool = "Public Network"
}
```

## Argument Reference

* `pool` - (Required) The IP pool to create a floating IP <br>From `Network > Floating IP` on the console, click `Create Floating IP` and check the IP pool.

## Attribute Reference

The following attributes are exported:

* `pool` - See Argument Reference above.
* `address` - The actual floating IP address itself.
* `port_id` - The ID of the associated port.
* `tenant_id` - The ID of the tenant in which to create the floating IP.
* `fixed_ip` - The fixed IP to which the floating IP maps.