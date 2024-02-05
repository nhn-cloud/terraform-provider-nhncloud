# Resource: nhncloud_keymanager_secret_v1

## Example Usage

```
resource "nhncloud_keymanager_secret_v1" "secret_01" {
  algorithm            = "aes"
  bit_length           = 256
  mode                 = "cbc"
  name                 = "mysecret"
  payload              = "foobar"
  payload_content_type = "text/plain"
  secret_type          = "passphrase"
}
```

## Argument Reference

* `name` - (Optional) The secret name.
* `expiration` - (Optional) The expiration date. Request in the ISO 8601 format.
* `algorithm` - (Optional) The encryption algorithm.
* `bit_length` - (Optional) The encryption key length.
* `mode` - (Optional) How block encryption works.
* `payload` - (Optional) The encryption key payload.
* `payload_content_type` - (Optional) The encryption key payload content type. </br>Required when entering a payload. </br>The list of supported content types: `text/plain`, `application/octet-stream`, `application/pkcs8`, `application/pkix-cert`
* `payload_content_encoding` - (Optional) Encoding encryption key payload. </br>Required if the payload_content_type is not `text/plain`</br>Only supports `base64`.
* `secret_type` - (Optional) The secret type. </br>One of the following: `symmetric`, `public`, `private`, `passphrase`, `certificate`, and `opaque`.



## Attribute Reference

The following attributes are exported:

* `secret_ref` - Reference for the secret </br>Where the secret is located.
* `name` - See Argument Reference above.
* `bit_length` - See Argument Reference above.
* `algorithm` - See Argument Reference above.
* `mode` - See Argument Reference above.
* `secret_type` - See Argument Reference above.
* `payload` - See Argument Reference above.
* `payload_content_type` - See Argument Reference above.
* `payload_content_encoding` - See Argument Reference above.
* `expiration` - See Argument Reference above.
* `content_types` - The map of the content types assigned on the secret.
* `creator_id` - The creator of the secret.
* `status` - The status of the secret.
* `created_at` - The date the secret was created.
* `updated_at` - The date the secret was last updated.