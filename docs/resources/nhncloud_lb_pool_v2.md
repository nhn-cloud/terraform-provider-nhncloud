# Resource: nhncloud_lb_pool_v2

## Example Usage

```
resource "nhncloud_lb_pool_v2" "tf_pool_01"{
  name = "tf_pool_01"
  description = "create pool by terraform."
  protocol = "HTTP"
  listener_id = nhncloud_lb_listener_v2.tf_listener_01.id
  lb_method = "LEAST_CONNECTIONS"
  persistence{
    type = "APP_COOKIE"
    cookie_name = "testCookie"
  }
  admin_state_up = true
}
```

## Argument Reference

* `name`  - (Optional) Load balancer name.
* `description`  - (Optional) Pool description.
* `protocol` - (Required) Protocol <br>One among `TCP`, `HTTP`, `HTTPS`, and `PROXY`.
* `listener_id` - (Required) The ID of the listener with which a pool to create is associated.
* `lb_method` - (Required) The load balancing method to distribute pool traffic to members <br>One among `ROUND_ROBIN`,`LEAST_CONNECTIONS`, and `SOURCE_IP`.
* `persistence` - (Optional) Session persistence of the pool to create.
* `persistence.type` - (Required) Session persistence type<br>One among `SOURCE_IP`, `HTTP_COOKIE`, and `APP_COOKIE` <br>Unavailable if the load balancing method is `SOURCE_IP`<br>HTTP_COOKIE and APP_COOKIE are unavailable if the protocol is `HTTPS` or `TCP`.
* `persistence.cookie_name` - (Optional) The name of cookie <br>persistence.cookie_name is available only when the session persistence type is APP_COOKIE.
* `admin_state_up` - (Optional) Administrator control status.
* `member_port` - (Optional) Member's receiving port. Traffic is sent to the port. Default is `-1`.

## Attribute Reference

The following attributes are exported:

* `id` - The unique ID for the pool.
* `tenant_id` - See Argument Reference above.
* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `protocol` - See Argument Reference above.
* `lb_method` - See Argument Reference above.
* `persistence` - See Argument Reference above.
* `admin_state_up` - See Argument Reference above.
* `member_port` - See Argument Reference above.
* `healthmonitor_id` - The health monitor ID of the pool.
* `operating_status` - The operating status of the member.