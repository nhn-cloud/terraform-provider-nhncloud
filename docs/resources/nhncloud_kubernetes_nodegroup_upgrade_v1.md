# Resource: nhncloud_kubernetes_nodegroup_upgrade_v1

~> **Note** This resource performs a one-time operation. It does not maintain ongoing state and should be used when you need to upgrade a nodegroup to a specific Kubernetes version.

## Example Usage

```
# Get the existing nodegroup
data "nhncloud_kubernetes_nodegroup_v1" "existing_nodegroup" {
  cluster_id = nhncloud_kubernetes_cluster_v1.my_cluster.id
  name       = "default-worker"
}

# Perform nodegroup upgrade
resource "nhncloud_kubernetes_nodegroup_upgrade_v1" "nodegroup_upgrade" {
  cluster_id   = nhncloud_kubernetes_cluster_v1.my_cluster.id
  nodegroup_id = data.nhncloud_kubernetes_nodegroup_v1.existing_nodegroup.id
  version      = "v1.32.3"

  options {
    num_buffer_nodes          = 1
    num_max_unavailable_nodes = 1
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to perform the upgrade operation.
* `cluster_id` - (Required) The UUID of the cluster that contains the nodegroup.
* `nodegroup_id` - (Required) The UUID of the nodegroup to upgrade.
* `version` - (Required) The target Kubernetes version for the upgrade.
* `options` - (Optional) Upgrade options configuration.

### Options Configuration

The `options` block supports:

* `num_buffer_nodes` - (Optional) The number of buffer nodes to create during upgrade. 
  Minimum: 0, Maximum: (max node quota per nodegroup - current node count), Default: 1.
* `num_max_unavailable_nodes` - (Optional) The maximum number of nodes that can be unavailable during upgrade.
  Minimum: 1, Maximum: current node count of the nodegroup, Default: 1.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The UUID of the upgrade operation.

## Import

Nodegroup upgrade operations cannot be imported as they represent one-time actions.

