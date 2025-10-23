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

* `id` - The UUID of the nodegroup.
* `uuid` - See Argument Reference above.
* `name` - See Argument Reference above.
* `cluster_id` - See Argument Reference above.
* `status` - The status of the nodegroup.
* `status_reason` - The reason for the nodegroup status.
* `flavor_id` - The flavor ID for nodegroup instances.
* `image_id` - The image ID for nodegroup instances.
* `labels` - A map of nodegroup labels.
* `max_node_count` - The maximum number of nodes in the nodegroup.
* `min_node_count` - The minimum number of nodes in the nodegroup.
* `node_addresses` - A list of node addresses in the nodegroup.
* `node_count` - The current number of nodes in the nodegroup.
* `project_id` - The project ID of the nodegroup.
* `role` - The role of the nodegroup.
* `stack_id` - The stack ID associated with the nodegroup.
* `version` - The Kubernetes version of the nodegroup.
* `created_at` - The creation timestamp of the nodegroup.
* `updated_at` - The last update timestamp of the nodegroup.

