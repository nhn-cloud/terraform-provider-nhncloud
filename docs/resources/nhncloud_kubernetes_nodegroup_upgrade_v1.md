# Resource: nhncloud_kubernetes_nodegroup_upgrade_v1

## Example Usage

```
# Get the existing nodegroup
data "nhncloud_kubernetes_nodegroup_v1" "existing_nodegroup" {
  cluster_id = nhncloud_kubernetes_cluster_v1.my_cluster.uuid
  name       = "default-worker"
}

# Perform nodegroup upgrade
resource "nhncloud_kubernetes_nodegroup_upgrade_v1" "nodegroup_upgrade" {
  cluster_id   = nhncloud_kubernetes_cluster_v1.my_cluster.uuid
  nodegroup_id = data.nhncloud_kubernetes_nodegroup_v1.existing_nodegroup.uuid
  version      = "v1.32.3"
  num_buffer_nodes          = 1
  num_max_unavailable_nodes = 1
}
```

~> **Important** When using this resource with `nhncloud_kubernetes_nodegroup_v1`, make sure to set the same `version` value or use `lifecycle { ignore_changes = [version] }` in the nodegroup resource to avoid conflicts.

## Argument Reference

The following arguments are supported:

* `region` - (Optional) Region to perform the upgrade operation.
* `cluster_id` - (Required) Cluster UUID.
* `nodegroup_id` - (Required) Node group UUID to upgrade.
* `version` - (Required) Target Kubernetes version (e.g., "v1.32.3").
* `num_buffer_nodes` - (Optional) Number of buffer nodes during upgrade. Min: 0, Max: (quota - current count), Default: 1.
* `num_max_unavailable_nodes` - (Optional) Maximum unavailable nodes during upgrade. Min: 1, Max: current node count, Default: 1.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `uuid` - Node group UUID.

## Import

Nodegroup upgrade resources can be imported using the nodegroup ID:

```
$ terraform import nhncloud_kubernetes_nodegroup_upgrade_v1.nodegroup_upgrade <nodegroup_id>
```

