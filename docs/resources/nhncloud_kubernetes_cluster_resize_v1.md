# Resource: nhncloud_kubernetes_cluster_resize_v1

~> **Note** This resource performs a one-time operation. It does not maintain ongoing state and should be used when you need to resize a specific nodegroup within a cluster.

## Example Usage

```
# Get the existing nodegroup
data "nhncloud_kubernetes_nodegroup_v1" "existing_nodegroup" {
  cluster_id = nhncloud_kubernetes_cluster_v1.my_cluster.uuid
  name       = "default-worker"
}

# Perform cluster resize (scale up)
resource "nhncloud_kubernetes_cluster_resize_v1" "cluster_resize" {
  cluster_id   = nhncloud_kubernetes_cluster_v1.my_cluster.uuid
  nodegroup_id = data.nhncloud_kubernetes_nodegroup_v1.existing_nodegroup.id
  node_count   = 5
}

# Perform cluster resize (scale down with specific nodes)
resource "nhncloud_kubernetes_cluster_resize_v1" "cluster_resize_down" {
  cluster_id   = nhncloud_kubernetes_cluster_v1.my_cluster.uuid
  nodegroup_id = data.nhncloud_kubernetes_nodegroup_v1.existing_nodegroup.id
  node_count   = 3
  nodes_to_remove = ["d6075d02-a8d1-4b5a-b6e2-95d7acaa99a4", "149z0690-4a81-4b32-ac30-81cc18fe50cc"]
}
```

~> **Important** When using this resource with `nhncloud_kubernetes_nodegroup_v1`, make sure to set the same `node_count` value or use `lifecycle { ignore_changes = [node_count] }` in the nodegroup resource to avoid conflicts.

## Argument Reference

The following arguments are supported:

* `region` - (Optional) Region to perform the resize operation.
* `cluster_id` - (Required) Cluster UUID.
* `nodegroup_id` - (Required) Node group UUID to resize.
* `node_count` - (Required) Target node count. When changed, performs resize automatically.
* `nodes_to_remove` - (Optional) List of node UUIDs to remove when scaling down. If not specified, nodes are selected automatically.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `uuid` - Cluster UUID.


## Import

Cluster resize resources can be imported using the cluster ID:

```
$ terraform import nhncloud_kubernetes_cluster_resize_v1.cluster_resize <cluster_id>
```

