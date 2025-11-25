# Resource: nhncloud_kubernetes_cluster_v1

## Example Usage

```
resource "nhncloud_kubernetes_cluster_v1" "my_cluster" {
  name                = "my-nks-cluster"
  cluster_template_id = "iaas_console"
  fixed_network       = "65e08bba-a7d5-4014-b28f-8f4fdaac001e"
  fixed_subnet        = "3a50e104-14d2-47a5-97a0-8ec87af8599e"
  flavor_id           = "b71c2d4e-31e4-4d0e-ac2f-f057ec4b6d71"
  keypair             = "my-keypair"
  node_count          = 2

  # Required labels for cluster configuration
  labels = {
    kube_tag                      = "v1.32.3"
    availability_zone             = "kr-pub-a"
    boot_volume_size              = "50"
    boot_volume_type              = "General HDD"
    ca_enable                     = "false"
    ca_max_node_count             = "9"
    ca_min_node_count             = "2"
    ca_scale_down_delay_after_add = "10"
    ca_scale_down_enable          = "false"
    ca_scale_down_unneeded_time   = "10"
    ca_scale_down_util_thresh     = "50"
    cert_manager_api              = "True"
    clusterautoscale              = "nodegroupfeature"
    external_network_id           = "751b8227-1234-5678-9349-dbf829d0aba5"
    external_subnet_id_list       = "59ddc195-1234-5678-9693-f09880747dc6"
    master_lb_floating_ip_enabled = "true"
    node_image                    = "1a10bf47-2f28-1234-5678-e2dc43f61789"
    strict_sg_rules               = "false"
  }

  # Required addons configuration
  addons {
    name    = "calico"
    version = "v3.28.2-nks1"
    options = {
      mode = "ebpf"
    }
  }

  addons {
    name    = "coredns"
    version = "1.8.4-nks1"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) Region to create the cluster.
* `name` - (Required) Cluster name. Changing this creates a new cluster.
* `cluster_template_id` - (Required) Cluster template ID (currently only "iaas_console" is supported). Changing this creates a new cluster.
* `fixed_network` - (Required) VPC network UUID. Changing this creates a new cluster.
* `fixed_subnet` - (Required) VPC subnet UUID. Changing this creates a new cluster.
* `flavor_id` - (Required) Instance flavor UUID for worker nodes. Changing this creates a new cluster.
* `keypair` - (Required) Keypair name for SSH access. Changing this creates a new cluster.
* `node_count` - (Optional) Initial number of worker nodes for the default node group. Defaults to 1 if not specified.
* `labels` - (Required) Cluster labels (key-value pairs for cluster configuration). Changing this creates a new cluster.
* `addons` - (Required) List of addons to install (CNI and CoreDNS are required). Changing this creates a new cluster.

### Labels Configuration

The `labels` block supports the following required arguments:

* `kube_tag` - (Required) Kubernetes version (e.g., "v1.32.3").
* `availability_zone` - (Required) Availability zone (e.g., "kr-pub-a").
* `boot_volume_size` - (Required) Boot volume size (GB).
* `boot_volume_type` - (Required) Boot volume type (e.g., "General HDD", "General SSD").
* `ca_enable` - (Required) Enable cluster autoscaler ("true"/"false").
* `cert_manager_api` - (Required) Enable CSR feature ("True"/"False").
* `master_lb_floating_ip_enabled` - (Required) Create public domain for API endpoint ("true"/"false").
* `node_image` - (Required) Base image UUID.

Optional labels include:

* `ca_max_node_count` - Maximum node count for autoscaler.
* `ca_min_node_count` - Minimum node count for autoscaler.
* `ca_scale_down_enable` - Enable autoscaler scale down ("true"/"false").
* `ca_scale_down_delay_after_add` - Delay after scale out (minutes).
* `ca_scale_down_unneeded_time` - Unneeded time before scale down (minutes).
* `ca_scale_down_util_thresh` - Utilization threshold for scale down (%).
* `external_network_id` - Internet Gateway network UUID.
* `external_subnet_id_list` - Colon-separated Internet Gateway subnet UUID list.
* `strict_sg_rules` - Apply strict security group rules ("true"/"false").

### Addons Configuration

The `addons` block supports:

* `name` - (Required) Addon name (e.g., "calico", "coredns").
* `version` - (Required) Addon version.
* `options` - (Optional) Addon-specific options (key-value pairs).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Cluster UUID.
* `status` - Cluster status.
* `status_reason` - Cluster status reason.
* `api_address` - Kubernetes API endpoint address.
* `node_addresses` - Worker node IP address list.
* `project_id` - Project ID.
* `stack_id` - Heat stack ID.
* `user_id` - User ID.
* `created_at` - Created time.
* `updated_at` - Last updated time.

## Lifecycle Configuration

### Ignoring node_count Changes

When using cluster autoscaler or managing node scaling separately with `nhncloud_kubernetes_cluster_resize_v1`, it's recommended to ignore changes to `node_count` to prevent Terraform from reverting autoscaler actions or manual scaling operations:

```
resource "nhncloud_kubernetes_cluster_v1" "my_cluster" {
  name       = "my-nks-cluster"
  node_count = 2
  # ... other configuration ...

  lifecycle {
    ignore_changes = [node_count]
  }
}
```

**Note:** Without these lifecycle configurations, Terraform will attempt to revert changes to the values in your configuration file on every apply, potentially overriding autoscaler actions, upgrade operations, or manual changes.
