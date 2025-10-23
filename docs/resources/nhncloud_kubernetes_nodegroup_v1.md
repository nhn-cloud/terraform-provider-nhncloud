# Resource: nhncloud_kubernetes_nodegroup_v1

## Example Usage

### Basic Nodegroup

```hcl
resource "nhncloud_kubernetes_nodegroup_v1" "my_nodegroup" {
  cluster_id  = nhncloud_kubernetes_cluster_v1.my_cluster.id
  name        = "additional-workers"
  node_count  = 3
  flavor_id   = "b71c2d4e-31e4-4d0e-ac2f-f057ec4b6d71"
  image_id    = "1a10bf47-2f28-1234-5678-e2dc43f61789"

  labels = {
    availability_zone             = "kr-pub-a"
    boot_volume_size              = "50"
    boot_volume_type              = "General HDD"
    ca_enable                     = "true"
    ca_max_node_count             = "10"
    ca_min_node_count             = "1"
    ca_scale_down_delay_after_add = "10"
    ca_scale_down_enable          = "true"
    ca_scale_down_unneeded_time   = "10"
    ca_scale_down_util_thresh     = "50"
  }
}
```

### Nodegroup with Version Upgrade

```hcl
resource "nhncloud_kubernetes_nodegroup_v1" "upgradeable_nodegroup" {
  cluster_id  = nhncloud_kubernetes_cluster_v1.my_cluster.id
  name        = "upgradeable-workers"
  node_count  = 3
  flavor_id   = "b71c2d4e-31e4-4d0e-ac2f-f057ec4b6d71"
  image_id    = "1a10bf47-2f28-1234-5678-e2dc43f61789"
  
  # Kubernetes 버전 지정 (업그레이드 가능)
  version = "v1.30.3"
  
  # 업그레이드 옵션
  upgrade_max_unavailable_nodes = 1  # 업그레이드 중 최대 사용 불가 노드 수
  upgrade_buffer_nodes          = 1  # 업그레이드 중 추가 버퍼 노드 수

  labels = {
    availability_zone             = "kr-pub-a"
    boot_volume_size              = "50"
    boot_volume_type              = "General HDD"
    ca_enable                     = "true"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to create the nodegroup.
* `cluster_id` - (Required) The UUID of the cluster that the nodegroup belongs to. Changing this creates a new nodegroup.
* `name` - (Required) The name of the nodegroup. Changing this creates a new nodegroup.
* `node_count` - (Required) The number of worker nodes in the nodegroup. This can be updated.
* `flavor_id` - (Required) The instance flavor UUID for nodegroup instances. Changing this creates a new nodegroup.
* `image_id` - (Required) The image UUID for nodegroup instances. Changing this creates a new nodegroup.
* `labels` - (Required) A map of nodegroup configuration labels. Changing this creates a new nodegroup.
* `version` - (Optional) The Kubernetes version of the nodegroup. When changed, triggers an upgrade operation.
* `min_node_count` - (Optional) The minimum number of nodes in the nodegroup for autoscaling.
* `max_node_count` - (Optional) The maximum number of nodes in the nodegroup for autoscaling.

### Upgrade Options

* `upgrade_max_unavailable_nodes` - (Optional) Maximum number of nodes that can be unavailable during upgrade. Default: 1.
* `upgrade_buffer_nodes` - (Optional) Number of additional buffer nodes to create during upgrade. Default: 0.

### Labels Configuration

The `labels` block supports the following required arguments:

* `availability_zone` - (Required) The availability zone for the nodegroup.
* `boot_volume_size` - (Required) The boot volume size in GB.
* `boot_volume_type` - (Required) The boot volume type.
* `ca_enable` - (Required) Whether to enable cluster autoscaler ("true"/"false").

Optional labels include:

* `ca_max_node_count` - The maximum number of nodes for autoscaler.
* `ca_min_node_count` - The minimum number of nodes for autoscaler.
* `ca_scale_down_enable` - Whether to enable scale-down ("true"/"false").
* `ca_scale_down_delay_after_add` - Scale-down delay after adding nodes (seconds).
* `ca_scale_down_unneeded_time` - Unneeded time before scale-down (seconds).
* `ca_scale_down_util_thresh` - CPU utilization threshold for scale-down (%).
* `user_script_v2` - User script for nodegroup instances.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The UUID of the nodegroup.
* `status` - The status of the nodegroup.
* `status_reason` - The reason for the nodegroup status.
* `max_node_count` - The maximum number of nodes in the nodegroup.
* `min_node_count` - The minimum number of nodes in the nodegroup.
* `node_addresses` - A list of node addresses in the nodegroup.
* `project_id` - The project ID of the nodegroup.
* `role` - The role of the nodegroup.
* `stack_id` - The stack ID associated with the nodegroup.
* `version` - The Kubernetes version of the nodegroup.
* `created_at` - The creation timestamp of the nodegroup.
* `updated_at` - The last update timestamp of the nodegroup.

## Kubernetes Version Upgrade

This resource supports Kubernetes version upgrades by changing the `version` attribute. When the version is changed, the provider will:

1. Trigger a nodegroup upgrade operation using the NHN Cloud API
2. Wait for the upgrade to complete
3. Update the resource state with the new version

### Supported Versions

Check the [NHN Cloud documentation](https://docs.nhncloud.com/ko/Container/NKS/ko/public-api/#_67) for currently supported Kubernetes versions.

### Upgrade Considerations

- **Irreversible**: Upgrades cannot be rolled back to previous versions
- **Sequential**: Avoid skipping versions; upgrade sequentially
- **Time**: Upgrades can take significant time depending on node count and workloads
- **Availability**: Some workload interruption may occur during upgrade

### Upgrade Example

```hcl
# Initial deployment with v1.29.3
resource "nhncloud_kubernetes_nodegroup_v1" "example" {
  # ... other configuration ...
  version = "v1.29.3"
  upgrade_max_unavailable_nodes = 1
  upgrade_buffer_nodes = 1
}

# To upgrade to v1.30.3, change the version and apply:
# version = "v1.30.3"
```
