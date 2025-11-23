# Data Source: nhncloud_kubernetes_nodegroup_v1

## Example Usage

```
# Get nodegroup by UUID
data "nhncloud_kubernetes_nodegroup_v1" "existing_nodegroup" {
  cluster_id = "abcd1234-efgh-5678-ijkl-9012mnop3456"
  uuid       = "mnop3456-ijkl-5678-efgh-1234abcd9012"
}

# Get nodegroup by name
data "nhncloud_kubernetes_nodegroup_v1" "existing_nodegroup" {
  cluster_id = nhncloud_kubernetes_cluster_v1.my_cluster.id
  name       = "default-worker"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the NKS client.
* `cluster_id` - (Required) The UUID of the cluster that the nodegroup belongs to.
* `uuid` - (Optional) The UUID of the nodegroup. Either `uuid` or `name` must be specified.
* `name` - (Optional) The name of the nodegroup. Either `uuid` or `name` must be specified.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Node group UUID.
* `uuid` - See Argument Reference above.
* `name` - See Argument Reference above.
* `cluster_id` - See Argument Reference above.
* `status` - Node group status.
* `status_reason` - Node group status reason.
* `flavor_id` - Instance flavor UUID.
* `image_id` - Base image UUID.
* `labels` - Node group labels (key-value pairs).
* `max_node_count` - Maximum node count.
* `min_node_count` - Minimum node count.
* `node_addresses` - Worker node IP address list.
* `node_count` - Current node count.
* `project_id` - Project ID.
* `role` - Node group role.
* `stack_id` - Heat stack ID.
* `version` - Kubernetes version.
* `created_at` - Created time.
* `updated_at` - Last updated time.

