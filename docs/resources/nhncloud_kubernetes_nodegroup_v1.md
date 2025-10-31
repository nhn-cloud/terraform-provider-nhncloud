# Resource: nhncloud_kubernetes_nodegroup_v1

## Example Usage

```
resource "nhncloud_kubernetes_nodegroup_v1" "my_nodegroup" {
  cluster_id = nhncloud_kubernetes_cluster_v1.my_cluster.uuid
  name       = "additional-workers"
  node_count = 3
  flavor_id   = "b71c2d4e-31e4-4d0e-ac2f-f057ec4b6d71"
  image_id    = "1a10bf47-2f28-1234-5678-e2dc43f61789"

  labels = {
    availability_zone = "kr-pub-a"
    boot_volume_size  = "50"
    boot_volume_type  = "General HDD"
    ca_enable = "False"
  }

  lifecycle {
    ignore_changes = [node_count]
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) Region to create the node group.
* `cluster_id` - (Required) Cluster UUID. Changing this creates a new node group.
* `name` - (Required) Node group name. Changing this creates a new node group.
* `node_count` - (Optional, Computed) Initial number of nodes. Defaults to 1 if not specified. After creation, this value becomes read-only and reflects the actual node count. To change the node count, use the `nhncloud_kubernetes_cluster_resize_v1` resource. See lifecycle configuration below.
* `flavor_id` - (Required) Instance flavor UUID. Changing this creates a new node group.
* `image_id` - (Required) Base image UUID. Changing this creates a new node group.
* `labels` - (Required) Node group labels (key-value pairs for configuration). Changing this creates a new node group.
* `version` - (Optional, Computed) Kubernetes version. After creation, this value becomes read-only and reflects the actual version. To upgrade the version, use the `nhncloud_kubernetes_nodegroup_upgrade_v1` resource.
* `min_node_count` - (Optional) Minimum node count for autoscaling. Can be updated.
* `max_node_count` - (Optional) Maximum node count for autoscaling. Can be updated.

### Labels Configuration

The `labels` block supports the following required arguments:

* `availability_zone` - (Required) Availability zone (e.g., "kr-pub-a").
* `boot_volume_size` - (Required) Boot volume size (GB).
* `boot_volume_type` - (Required) Boot volume type (e.g., "General HDD", "General SSD").
* `ca_enable` - (Required) Enable cluster autoscaler ("true"/"false").

Optional labels include:

* `ca_max_node_count` - Maximum node count for autoscaler.
* `ca_min_node_count` - Minimum node count for autoscaler.
* `ca_scale_down_enable` - Enable autoscaler scale down ("true"/"false").
* `ca_scale_down_delay_after_add` - Delay after scale out (minutes).
* `ca_scale_down_unneeded_time` - Unneeded time before scale down (minutes).
* `ca_scale_down_util_thresh` - Utilization threshold for scale down (%).
* `user_script_v2` - User script (base64 encoded).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `uuid` - Node group UUID.
* `status` - Node group status.
* `status_reason` - Node group status reason.
* `max_node_count` - Maximum node count.
* `min_node_count` - Minimum node count.
* `node_addresses` - Worker node IP address list.
* `project_id` - Project ID.
* `role` - Node group role.
* `stack_id` - Heat stack ID.
* `version` - Kubernetes version.
* `created_at` - Created time.
* `updated_at` - Last updated time.

## Lifecycle Configuration

### Managing Node Count and Version

Node count and version cannot be updated directly on the nodegroup resource after creation. Use dedicated resources for these operations:

- **For node scaling**: Use `nhncloud_kubernetes_cluster_resize_v1`
- **For version upgrades**: Use `nhncloud_kubernetes_nodegroup_upgrade_v1`

### Ignoring node_count Changes

When using autoscaler or managing node scaling with the resize resource, it's recommended to ignore changes to `node_count` to prevent Terraform from detecting drift:

```
resource "nhncloud_kubernetes_nodegroup_v1" "my_nodegroup" {
  cluster_id = nhncloud_kubernetes_cluster_v1.my_cluster.uuid
  name       = "workers"
  # node_count can be omitted - defaults to 1
  flavor_id  = "b71c2d4e-31e4-4d0e-ac2f-f057ec4b6d71"
  image_id   = "1a10bf47-2f28-1234-5678-e2dc43f61789"
  
  labels = {
    availability_zone = "kr-pub-a"
    boot_volume_size  = "50"
    boot_volume_type  = "General HDD"
    ca_enable = "False"
  }

  lifecycle {
    ignore_changes = [node_count]
  }
}

# Use resize resource for scaling
resource "nhncloud_kubernetes_cluster_resize_v1" "scale_workers" {
  cluster_id   = nhncloud_kubernetes_cluster_v1.my_cluster.id
  nodegroup_id = nhncloud_kubernetes_nodegroup_v1.my_nodegroup.id
  node_count   = 5
}
```

### Ignoring version Changes

When using `nhncloud_kubernetes_nodegroup_upgrade_v1` resource to manage Kubernetes version upgrades separately, you should also ignore version changes:

```hcl
resource "nhncloud_kubernetes_nodegroup_v1" "my_nodegroup" {
  cluster_id = nhncloud_kubernetes_cluster_v1.my_cluster.uuid
  name       = "workers"
  # version can be omitted - will use cluster's version
  # ... other configuration ...

  lifecycle {
    ignore_changes = [node_count, version]
  }
}

# Use upgrade resource for version upgrades
resource "nhncloud_kubernetes_nodegroup_upgrade_v1" "upgrade_workers" {
  cluster_id   = nhncloud_kubernetes_cluster_v1.my_cluster.id
  nodegroup_id = nhncloud_kubernetes_nodegroup_v1.my_nodegroup.id
  version      = "v1.32.3"
  
  num_buffer_nodes            = 1
  num_max_unavailable_nodes   = 1
}
```

### Combining Multiple Lifecycle Rules

You can combine multiple `ignore_changes` rules when needed:

```
resource "nhncloud_kubernetes_nodegroup_v1" "my_nodegroup" {
  cluster_id = nhncloud_kubernetes_cluster_v1.my_cluster.uuid
  name       = "additional-workers"
  node_count = 3
  version    = "v1.30.3"
  # ... other configuration ...

  lifecycle {
    ignore_changes = [node_count, version]
  }
}
```

**Note:** Without these lifecycle configurations, Terraform will attempt to revert changes to the values in your configuration file on every apply, potentially overriding autoscaler actions, upgrade operations, or manual changes.
