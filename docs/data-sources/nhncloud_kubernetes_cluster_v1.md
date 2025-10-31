# Data Source: nhncloud_kubernetes_cluster_v1

## Example Usage

```
# Get cluster by UUID
data "nhncloud_kubernetes_cluster_v1" "existing_cluster" {
  uuid = "abcd1234-efgh-5678-ijkl-9012mnop3456"
}

# Get cluster by name
data "nhncloud_kubernetes_cluster_v1" "existing_cluster" {
  name = "my-nks-cluster"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the NKS client.
* `uuid` - (Optional) The UUID of the cluster. Either `uuid` or `name` must be specified.
* `name` - (Optional) The name of the cluster. Either `uuid` or `name` must be specified.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Cluster UUID.
* `uuid` - See Argument Reference above.
* `name` - See Argument Reference above.
* `status` - Cluster status.
* `status_reason` - Cluster status reason.
* `api_address` - Kubernetes API endpoint address.
* `cluster_template_id` - Cluster template ID.
* `create_timeout` - Cluster creation timeout (minutes).
* `discovery_url` - Discovery URL.
* `fixed_network` - VPC network UUID.
* `fixed_subnet` - VPC subnet UUID.
* `flavor_id` - Instance flavor UUID.
* `keypair` - Keypair name.
* `labels` - Cluster labels (key-value pairs).
* `node_addresses` - Worker node IP address list.
* `node_count` - Number of worker nodes.
* `project_id` - Project ID.
* `stack_id` - Heat stack ID.
* `user_id` - User ID.
* `created_at` - Created time.
* `updated_at` - Last updated time.

