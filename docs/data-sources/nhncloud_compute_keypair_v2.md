# Data Source: nhncloud_compute_keypair_v2

## Example Usage

```
data "nhncloud_compute_keypair_v2" "my_keypair"{
  name = "my_keypair"
}
```

## Argument Reference

* `name` - (Required) The unique name of the keypair.
* `region` - (Optional) The region name to which the keypair to query belongs.

## Attribute Reference

* `name` - See Argument Reference above.
* `region` - See Argument Reference above.
* `public_key` - The OpenSSH-formatted public key of the keypair.
* `fingerprint` - The fingerprint of the public key.