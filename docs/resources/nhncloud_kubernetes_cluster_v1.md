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

  # Optional API endpoint IP access control
  api_ep_ipacl {
    enable = "True"
    action = "ALLOW"
    ipacl_targets {
      cidr_address = "192.168.0.5"
      description  = "My Friend"
    }
    ipacl_targets {
      cidr_address = "10.10.22.3/24"
      description  = "Your Friends"
    }
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

* `region` - (Optional) The region in which to create the cluster.
* `name` - (Required) The name of the cluster. Changing this creates a new cluster.
* `cluster_template_id` - (Required) The cluster template ID to use. Currently only "iaas_console" is supported. Changing this creates a new cluster.
* `fixed_network` - (Required) The VPC network UUID. Changing this creates a new cluster.
* `fixed_subnet` - (Required) The VPC subnet UUID. Changing this creates a new cluster.
* `flavor_id` - (Required) The instance flavor UUID for cluster nodes. Changing this creates a new cluster.
* `keypair` - (Required) The SSH keypair name for cluster nodes. Changing this creates a new cluster.
* `node_count` - (Required) The number of worker nodes. Changing this creates a new cluster.
* `labels` - (Required) A map of cluster configuration labels. Changing this creates a new cluster.
* `api_ep_ipacl` - (Optional) API endpoint IP access control configuration. Changing this creates a new cluster.
* `addons` - (Required) A list of addons to install on the cluster. Changing this creates a new cluster.

### Labels Configuration

The `labels` block supports the following required arguments:

* `kube_tag` - (Required) The Kubernetes version.
* `availability_zone` - (Required) The availability zone for the default worker nodegroup.
* `boot_volume_size` - (Required) The boot volume size in GB for the default worker nodegroup.
* `boot_volume_type` - (Required) The boot volume type for the default worker nodegroup.
* `ca_enable` - (Required) Whether to enable cluster autoscaler ("true"/"false").
* `cert_manager_api` - (Required) Whether to enable CSR feature ("True").
* `master_lb_floating_ip_enabled` - (Required) Whether to assign public IP to API endpoint ("true"/"false").
* `node_image` - (Required) The base image UUID for the default worker nodegroup.

Optional labels include:

* `ca_max_node_count` - The maximum number of nodes for autoscaler.
* `ca_min_node_count` - The minimum number of nodes for autoscaler.
* `ca_scale_down_enable` - Whether to enable scale-down ("true"/"false").
* `ca_scale_down_delay_after_add` - Scale-down delay after adding nodes (seconds).
* `ca_scale_down_unneeded_time` - Unneeded time before scale-down (seconds).
* `ca_scale_down_util_thresh` - CPU utilization threshold for scale-down (%).
* `external_network_id` - External network UUID connected to Internet Gateway.
* `external_subnet_id_list` - External subnet UUID list (colon-separated).
* `strict_sg_rules` - Whether to create minimal security group rules ("true"/"false").

### API Endpoint IP ACL Configuration

The `api_ep_ipacl` block supports:

* `enable` - (Required) Whether to enable IP access control.
* `action` - (Required) The IP access control action type.
* `ipacl_targets` - (Required) A list of IP access control targets.

The `ipacl_targets` block supports:

* `cidr_address` - (Required) IP address or CIDR range.
* `description` - (Required) Description for the IP access control target.

### Addons Configuration

The `addons` block supports:

* `name` - (Required) The addon name.
* `version` - (Required) The addon version.
* `options` - (Optional) A map of addon-specific options.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The UUID of the cluster.
* `status` - The status of the cluster.
* `status_reason` - The reason for the cluster status.
* `api_address` - The API server endpoint address.
* `master_addresses` - A list of master node addresses.
* `master_count` - The number of master nodes.
* `node_addresses` - A list of worker node addresses.
* `project_id` - The project ID of the cluster.
* `stack_id` - The stack ID associated with the cluster.
* `user_id` - The user ID who created the cluster.
* `created_at` - The creation timestamp of the cluster.
* `updated_at` - The last update timestamp of the cluster.

