# Resource: nhncloud_kubernetes_cluster_resize_v1

~> **Note** This resource performs a one-time operation. It does not maintain ongoing state and should be used when you need to resize a specific nodegroup within a cluster.

## Example Usage

```
# Get the existing nodegroup
data "nhncloud_kubernetes_nodegroup_v1" "existing_nodegroup" {
  cluster_id = nhncloud_kubernetes_cluster_v1.my_cluster.id
  name       = "default-worker"
}

# Perform cluster resize (scale up)
resource "nhncloud_kubernetes_cluster_resize_v1" "cluster_resize" {
  cluster_id   = nhncloud_kubernetes_cluster_v1.my_cluster.id
  nodegroup_id = data.nhncloud_kubernetes_nodegroup_v1.existing_nodegroup.id
  node_count   = 5

  options {
    nodes_to_remove = []
  }
}

# Perform cluster resize (scale down with specific nodes)
resource "nhncloud_kubernetes_cluster_resize_v1" "cluster_resize_down" {
  cluster_id   = nhncloud_kubernetes_cluster_v1.my_cluster.id
  nodegroup_id = data.nhncloud_kubernetes_nodegroup_v1.existing_nodegroup.id
  node_count   = 2

  options {
    nodes_to_remove = ["node-uuid-1", "node-uuid-2"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to perform the resize operation.
* `cluster_id` - (Required) The UUID of the cluster that contains the nodegroup.
* `nodegroup_id` - (Required) The UUID of the nodegroup to resize.
* `node_count` - (Required) The target number of worker nodes after resize.
* `options` - (Optional) Resize options configuration.

### Options Configuration

The `options` block supports:

* `nodes_to_remove` - (Optional) A list of specific node UUIDs to remove when scaling down.
  If not specified when scaling down, nodes will be selected randomly for removal.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The UUID of the resize operation.


## Import

Cluster resize operations cannot be imported as they represent one-time actions.

