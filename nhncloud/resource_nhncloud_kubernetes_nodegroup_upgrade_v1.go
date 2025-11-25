package nhncloud

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/nhn-cloud/nhncloud.gophercloud/nhncloud/kubernetes/v1/nodegroups"
)

func resourceKubernetesNodegroupUpgradeV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKubernetesNodegroupUpgradeV1Create,
		ReadContext:   resourceKubernetesNodegroupUpgradeV1Read,
		UpdateContext: resourceKubernetesNodegroupUpgradeV1Update,
		DeleteContext: resourceKubernetesNodegroupUpgradeV1Delete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"nodegroup_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},

			"num_buffer_nodes": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
				Default:  1,
			},

			"num_max_unavailable_nodes": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
				Default:  1,
			},

			// computed-only
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceKubernetesNodegroupUpgradeV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	kubernetesClient, err := config.ContainerInfraV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud Kubernetes client: %s", err)
	}

	// Set container-infra API microversion to latest for NKS compatibility
	kubernetesClient.Microversion = kubernetesV1NodeGroupMinMicroversion

	clusterIDOrName := d.Get("cluster_id").(string)
	nodegroupInput := d.Get("nodegroup_id").(string)

	// Extract nodegroup_id from "cluster_id/nodegroup_id" format if necessary
	// This allows using nhncloud_kubernetes_nodegroup_v1.resource_name.id directly
	nodegroupIDOrName := extractNodeGroupIDForUpgrade(nodegroupInput)

	// Build upgrade options
	upgradeOpts := nodegroups.UpgradeOpts{
		Version: d.Get("version").(string),
	}

	// Set upgrade options as direct fields
	if v, ok := d.GetOk("num_buffer_nodes"); ok {
		upgradeOpts.NumBufferNodes = v.(int)
	}
	if v, ok := d.GetOk("num_max_unavailable_nodes"); ok {
		upgradeOpts.NumMaxUnavailableNodes = v.(int)
	}

	log.Printf("[DEBUG] Checking cluster status before upgrading NKS nodegroup %s in cluster %s to version %s", nodegroupIDOrName, clusterIDOrName, upgradeOpts.Version)

	// Wait for cluster to be in a stable state before attempting upgrade
	clusterStateConf := &resource.StateChangeConf{
		Pending:    []string{"UPDATE_IN_PROGRESS", "CREATE_IN_PROGRESS"},
		Target:     []string{"UPDATE_COMPLETE", "CREATE_COMPLETE"},
		Refresh:    kubernetesClusterV1StateRefreshFunc(kubernetesClient, clusterIDOrName),
		Timeout:    10 * time.Minute, // Wait up to 10 minutes for cluster to stabilize
		Delay:      30 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	log.Printf("[DEBUG] Waiting for cluster %s to be in stable state before nodegroup upgrade", clusterIDOrName)
	_, err = clusterStateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for cluster %s to be in stable state before nodegroup upgrade: %s", clusterIDOrName, err)
	}

	log.Printf("[DEBUG] Upgrading NKS nodegroup %s in cluster %s to version %s", nodegroupIDOrName, clusterIDOrName, upgradeOpts.Version)

	nodegroup, err := nodegroups.Upgrade(kubernetesClient, clusterIDOrName, nodegroupIDOrName, upgradeOpts).Extract()
	if err != nil {
		return diag.Errorf("Error upgrading NKS nodegroup %s in cluster %s: %s", nodegroupIDOrName, clusterIDOrName, err)
	}

	d.SetId(nodegroup.UUID)

	if nodegroupIDOrName == "default-master" {
		log.Printf("[DEBUG] Waiting for cluster %s to complete default-master nodegroup upgrade", clusterIDOrName)

		clusterUpgradeStateConf := &resource.StateChangeConf{
			Pending:    []string{"UPDATE_IN_PROGRESS"},
			Target:     []string{"UPDATE_COMPLETE", "CREATE_COMPLETE"},
			Refresh:    kubernetesClusterV1StateRefreshFunc(kubernetesClient, clusterIDOrName),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      30 * time.Second,
			MinTimeout: 10 * time.Second,
		}

		_, err = clusterUpgradeStateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("Error waiting for cluster %s to complete default-master upgrade: %s", clusterIDOrName, err)
		}
	} else {
		log.Printf("[DEBUG] Waiting for nodegroup %s upgrade to complete", nodegroup.UUID)

		nodegroupUpgradeStateConf := &resource.StateChangeConf{
			Pending:    []string{"UPDATE_IN_PROGRESS"},
			Target:     []string{"UPDATE_COMPLETE"},
			Refresh:    kubernetesNodeGroupV1StateRefreshFunc(kubernetesClient, clusterIDOrName, nodegroup.UUID),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      30 * time.Second,
			MinTimeout: 10 * time.Second,
		}

		_, err = nodegroupUpgradeStateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("Error waiting for nodegroup %s upgrade to complete: %s", nodegroup.UUID, err)
		}
	}

	return resourceKubernetesNodegroupUpgradeV1Read(ctx, d, meta)
}

func resourceKubernetesNodegroupUpgradeV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	kubernetesClient, err := config.ContainerInfraV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud Kubernetes client: %s", err)
	}

	// Set container-infra API microversion to latest for NKS compatibility
	kubernetesClient.Microversion = kubernetesV1NodeGroupMinMicroversion

	clusterIDOrName := d.Get("cluster_id").(string)
	nodegroupInput := d.Get("nodegroup_id").(string)
	nodegroupIDOrName := extractNodeGroupIDForUpgrade(nodegroupInput)

	nodegroup, err := nodegroups.Get(kubernetesClient, clusterIDOrName, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(gophercloud.ErrDefault403); ok {
			// default-master nodegroup access is restricted by NHN Cloud policy
			if nodegroupIDOrName == "default-master" {
				log.Printf("[INFO] default-master nodegroup access is restricted by NHN Cloud policy")
				return nil
			}
			// For other nodegroups, 403 error is unexpected and should be reported
			return diag.Errorf("Error retrieving NKS nodegroup %s: %s", d.Id(), err)
		}
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving NKS nodegroup"))
	}

	log.Printf("[DEBUG] Retrieved NKS nodegroup %s: %+v", d.Id(), nodegroup)

	d.Set("uuid", nodegroup.UUID)
	d.Set("status", nodegroup.Status)
	d.Set("upgraded_at", nodegroup.UpdatedAt)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceKubernetesNodegroupUpgradeV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	kubernetesClient, err := config.ContainerInfraV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud Kubernetes client: %s", err)
	}

	kubernetesClient.Microversion = kubernetesV1NodeGroupMinMicroversion

	if d.HasChange("version") || d.HasChange("num_buffer_nodes") || d.HasChange("num_max_unavailable_nodes") {
		clusterIDOrName := d.Get("cluster_id").(string)
		nodegroupInput := d.Get("nodegroup_id").(string)
		nodegroupIDOrName := extractNodeGroupIDForUpgrade(nodegroupInput)

		// Build upgrade options
		upgradeOpts := nodegroups.UpgradeOpts{
			Version: d.Get("version").(string),
		}

		if v, ok := d.GetOk("num_buffer_nodes"); ok {
			upgradeOpts.NumBufferNodes = v.(int)
		}
		if v, ok := d.GetOk("num_max_unavailable_nodes"); ok {
			upgradeOpts.NumMaxUnavailableNodes = v.(int)
		}

		log.Printf("[DEBUG] Waiting for cluster %s to be in stable state before nodegroup upgrade", clusterIDOrName)

		// Wait for cluster to be in a stable state before attempting upgrade
		clusterStateConf := &resource.StateChangeConf{
			Pending:    []string{"UPDATE_IN_PROGRESS", "CREATE_IN_PROGRESS"},
			Target:     []string{"UPDATE_COMPLETE", "CREATE_COMPLETE"},
			Refresh:    kubernetesClusterV1StateRefreshFunc(kubernetesClient, clusterIDOrName),
			Timeout:    10 * time.Minute,
			Delay:      30 * time.Second,
			MinTimeout: 10 * time.Second,
		}

		_, err = clusterStateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("Error waiting for cluster %s to be in stable state before nodegroup upgrade: %s", clusterIDOrName, err)
		}

		log.Printf("[DEBUG] Upgrading NKS nodegroup %s in cluster %s to version %s", nodegroupIDOrName, clusterIDOrName, upgradeOpts.Version)

		nodegroup, err := nodegroups.Upgrade(kubernetesClient, clusterIDOrName, nodegroupIDOrName, upgradeOpts).Extract()
		if err != nil {
			return diag.Errorf("Error upgrading NKS nodegroup %s in cluster %s: %s", nodegroupIDOrName, clusterIDOrName, err)
		}

		// For default-master nodegroup, verify upgrade completion via cluster status
		// because direct nodegroup access is restricted by NHN Cloud policy
		if nodegroupIDOrName == "default-master" {
			log.Printf("[DEBUG] Waiting for cluster %s to complete default-master nodegroup upgrade", clusterIDOrName)
			log.Printf("[INFO] Verifying upgrade completion via cluster status due to master nodegroup access restrictions")

			clusterUpgradeStateConf := &resource.StateChangeConf{
				Pending:    []string{"UPDATE_IN_PROGRESS"},
				Target:     []string{"UPDATE_COMPLETE", "CREATE_COMPLETE"},
				Refresh:    kubernetesClusterV1StateRefreshFunc(kubernetesClient, clusterIDOrName),
				Timeout:    d.Timeout(schema.TimeoutUpdate),
				Delay:      30 * time.Second,
				MinTimeout: 10 * time.Second,
			}

			_, err = clusterUpgradeStateConf.WaitForStateContext(ctx)
			if err != nil {
				return diag.Errorf("Error waiting for cluster %s to complete default-master upgrade: %s", clusterIDOrName, err)
			}
		} else {
			// For regular nodegroups, verify upgrade completion via nodegroup status directly
			log.Printf("[DEBUG] Waiting for nodegroup %s upgrade to complete", nodegroup.UUID)

			nodegroupUpgradeStateConf := &resource.StateChangeConf{
				Pending:    []string{"UPDATE_IN_PROGRESS"},
				Target:     []string{"UPDATE_COMPLETE"},
				Refresh:    kubernetesNodeGroupV1StateRefreshFunc(kubernetesClient, clusterIDOrName, nodegroup.UUID),
				Timeout:    d.Timeout(schema.TimeoutUpdate),
				Delay:      30 * time.Second,
				MinTimeout: 10 * time.Second,
			}

			_, err = nodegroupUpgradeStateConf.WaitForStateContext(ctx)
			if err != nil {
				return diag.Errorf("Error waiting for nodegroup %s upgrade to complete: %s", nodegroup.UUID, err)
			}
		}
	}

	return resourceKubernetesNodegroupUpgradeV1Read(ctx, d, meta)
}

func resourceKubernetesNodegroupUpgradeV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// This is a one-time operation resource, so delete just removes it from state
	log.Printf("[DEBUG] Removing NKS nodegroup upgrade resource %s from state", d.Id())
	return nil
}

// extractNodeGroupIDForUpgrade extracts the nodegroup ID from either a simple ID or "cluster_id/nodegroup_id" format
// This allows users to reference nhncloud_kubernetes_nodegroup_v1 resources directly
func extractNodeGroupIDForUpgrade(input string) string {
	parts := strings.Split(input, "/")
	if len(parts) >= 2 {
		// Return the last part if in "cluster_id/nodegroup_id" format
		return parts[len(parts)-1]
	}
	// Return as-is if already a simple ID
	return input
}
