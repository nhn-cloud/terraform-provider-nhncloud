# Data Source: nhncloud_keymanager_container_v1

## Example Usage

```
data "nhncloud_keymanager_container_v1" "container_01" {
  name      = "terraform_container_01"
}
```

## Argument Reference

* `region` - (Optional) The region name to which the secret container you want to look up belongs.
* `name` - (Optional) The secret container name to query.

## Attribute Reference

`id` is set to the ID of the found secret container. In addition, the following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
