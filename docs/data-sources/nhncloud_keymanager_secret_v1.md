# Data Source: nhncloud_keymanager_secret_v1

## Example Usage

```
data "nhncloud_keymanager_secret_v1" "secret_01" {
  name      = "terraform_secret_01"
}
```

## Argument Reference

* `region` - (Optional) The region name to which the secret to query belongs.
* `name` - (Optional) The secret name to query.

## Attribute Reference

`id` is set to the ID of the found secret. In addition, the following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
