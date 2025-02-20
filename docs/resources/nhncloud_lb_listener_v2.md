# Resource: nhncloud_lb_listener_v2

## Example Usage

```
# HTTP Listener
resource "nhncloud_lb_listener_v2" "tf_listener_http_01"{
  name = "tf_listener_01"
  description = "create listener by terraform."
  protocol = "HTTP"
  protocol_port = 80
  loadbalancer_id = nhncloud_lb_loadbalancer_v2.tf_loadbalancer_01.id
  default_pool_id = ""
  connection_limit = 2000
  timeout_client_data = 5000
  timeout_member_connect = 5000
  timeout_member_data = 5000
  timeout_tcp_inspect = 5000
  admin_state_up = true
}

# Terminated HTTPS Listener
resource "nhncloud_lb_listener_v2" "tf_listener_01"{
  name = "tf_listener_01"
  description = "create listener by terraform."
  protocol = "TERMINATED_HTTPS"
  protocol_port = 443
  loadbalancer_id = nhncloud_lb_loadbalancer_v2.tf_loadbalancer_01.id
  default_pool_id = ""
  connection_limit = 2000
  timeout_client_data = 5000
  timeout_member_connect = 5000
  timeout_member_data = 5000
  timeout_tcp_inspect = 5000
  default_tls_container_ref = "https://kr1-api-key-manager-infrastructure.nhncloudservice.com/v1/containers/3258d456-06f4-48c5-8863-acf9facb26de"
  sni_container_refs = null
  admin_state_up = true
}
```


## Argument Reference

* `name` - (Optional) The name of the listener to create.
* `description` - (Optional) The description of the listener.
* `protocol` - (Required) The protocol of the listener to create <br>One among `TCP`, `HTTP`, `HTTPS`, and `TERMINATED_HTTPS`.
* `protocol_port` - (Required) The port of the listener to create.
* `loadbalancer_id` - (Required) The ID of the load balancer to be connected with the listener to create.
* `default_pool_id` - (Optional) The ID of the default pool to be connected with the listener to create.
* `connection_limit` - (Optional) The maximum connection count allowed for the listener to create.
* `timeout_client_data` - (Optional) The timeout setting when the client is inactive (ms).
* `timeout_member_connect` - (Optional) The timeout setting when the member is connected (ms).
* `timeout_member_data` - (Optional) The timeout setting when the member is inactive (ms).
* `timeout_tcp_inspect` - (Optional) The timeout to wait for additional TCP packets for content inspection (ms).
* `default_tls_container_ref` - (Optional) The path of TLC certificate to be used when the protocol is `TERMINATED_HTTPS`.
* `sni_container_refs` - (Optional) The list of SNI certificate paths.
* `insert_headers` - (Optional) The list of headers to be added before a request is sent to a backend member.
* `admin_state_up` - (Optional) Administrator control status.
* `keepalive_timeout` - (Optional) Keepalive timeout of listener.


## Attribute Reference

The following attributes are exported:

* `id` - The unique ID for the Listener.
* `protocol` - See Argument Reference above.
* `protocol_port` - See Argument Reference above.
* `tenant_id` - See Argument Reference above.
* `name` - See Argument Reference above.
* `default_port_id` - See Argument Reference above.
* `description` - See Argument Reference above.
* `connection_limit` - See Argument Reference above.
* `timeout_client_data` - See Argument Reference above.
* `timeout_member_connect` - See Argument Reference above.
* `timeout_member_data` - See Argument Reference above.
* `timeout_tcp_inspect` - See Argument Reference above.
* `default_tls_container_ref` - See Argument Reference above.
* `sni_container_refs` - See Argument Reference above.
* `admin_state_up` - See Argument Reference above.
* `insert_headers` - See Argument Reference above.
* `allowed_cidrs` - See Argument Reference above.
* `keepalive_timeout` - See Argument Reference above.