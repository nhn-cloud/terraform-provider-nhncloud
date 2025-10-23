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

* `id` - The UUID of the cluster.
* `uuid` - See Argument Reference above.
* `name` - See Argument Reference above.
* `status` - The status of the cluster.
* `status_reason` - The reason for the cluster status.
* `api_address` - The API server endpoint address.
* `cluster_template_id` - The cluster template ID used for creating the cluster.
* `create_timeout` - The timeout for cluster creation.
* `discovery_url` - The discovery URL for the cluster.
* `fixed_network` - The fixed network UUID for the cluster.
* `fixed_subnet` - The fixed subnet UUID for the cluster.
* `flavor_id` - The flavor ID for cluster nodes.
* `keypair` - The keypair name for cluster nodes.
* `labels` - A map of cluster labels.
* `master_addresses` - A list of master node addresses.
* `master_count` - The number of master nodes.
* `node_addresses` - A list of worker node addresses.
* `node_count` - The number of worker nodes.
* `project_id` - The project ID of the cluster.
* `stack_id` - The stack ID associated with the cluster.
* `user_id` - The user ID who created the cluster.
* `created_at` - The creation timestamp of the cluster.
* `updated_at` - The last update timestamp of the cluster.

