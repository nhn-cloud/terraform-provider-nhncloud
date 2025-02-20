# Resource: nhncloud_lb_monitor_v2

## Example Usage

```
resource "nhncloud_lb_monitor_v2" "tf_monitor_01"{
  name = "tf_monitor_01"
  pool_id = nhncloud_lb_pool_v2.tf_pool_01.id
  type = "HTTP"
  delay = 20
  timeout = 10
  max_retries = 5
  url_path = "/"
  http_method = "GET"
  expected_codes = "200-202"
  admin_state_up = true
  health_check_port = 2000
}
```

## Argument Reference

* `name` - (Optional) The name of the health monitor to create.
* `pool_id` - (Required) The ID of the pool to be connected with the health monitor to create.
* `type` - (Required) Support `TCP`, `HTTP`, and `HTTPS` only.
* `delay` - (Required) Interval of status check.
* `timeout` - (Required) Timeout for status check (seconds)<br>Timeout must have smaller value than delay.
* `max_retries` - (Required) The number of maximum retries, between 1 and 10.
* `url_path` - (Optional) The request URL for status check.
* `http_method` - (Optional) The HTTP method to use for status check<br>The default is GET.
* `expected_codes` - (Optional) HTTP(S) response code of members to be considered as normal status <br/>expected_codes can be set as list (`200,201,202`) or range (`200-202`).
* `admin_state_up` - (Optional) Administrator control status.
* `host_header` - (Optional) The host header field value to use for status check. When the status check type is set with `TCP`,  the value set in this field will be ignored.
* `health_check_port` - (Optional) The member port to be health-checked.

## Attribute Reference

The following attributes are exported:

* `id` - The unique ID for the monitor.
* `type` - See Argument Reference above.
* `delay` - See Argument Reference above.
* `timeout` - See Argument Reference above.
* `max_retries` - See Argument Reference above.
* `max_retries_down` - See Argument Reference above.
* `url_path` - See Argument Reference above.
* `http_method` - See Argument Reference above.
* `expected_codes` - See Argument Reference above.
* `admin_state_up` - See Argument Reference above.
* `host_header` - See Argument Reference above.
* `health_check_port` - See Argument Reference above.