# Resource: nhncloud_keymanager_container_v1

## Example Usage

```
resource "nhncloud_keymanager_secret_v1" "secret_01" {
...
}

resource "nhncloud_keymanager_container_v1" "container_01" {
  name      = "terraform_container_01"
  type      = "generic"
  secret_refs {
    secret_ref = nhncloud_keymanager_secret_v1.secret_01.secret_ref
  }
}
```

## Argument Reference

* `type` - (Required) The container type. </br>One of `generic`, `rsa`, `certificate`.
* `name` - (Optional) The container name.
* `secret_refs` - (Optional) List of secrets to register in the container.
* `secret_refs.secret_ref` - (Optional) The secret address.
* `secret_refs.name` - (Optional) The secret name specified by the container. </br>If container type is `certificate`: Specify `as` `certificate`, `private_key`, `private_key_passphrase`, and `intermediates`. </br>If container type is `rsa`: Specify `as` `private_key`, `private_key_passphrase`, and `public_key`.

## Attribute Reference

The following attributes are exported:

* `container_ref` - The container reference <br>where to find the container.
* `name` - See Argument Reference above.
* `type` - See Argument Reference above.
* `secret_refs` - See Argument Reference above.
* `creator_id` - The creator of the container.
* `status` - The status of the container.