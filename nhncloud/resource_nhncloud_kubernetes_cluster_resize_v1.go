package nhncloud

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/nhn-cloud/nhncloud.gophercloud/nhncloud/kubernetes/v1/clusters"
)

func resourceKubernetesClusterResizeV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKubernetesClusterResizeV1Create,
		ReadContext:   resourceKubernetesClusterResizeV1Read,
		UpdateContext: resourceKubernetesClusterResizeV1Update,
		DeleteContext: resourceKubernetesClusterResizeV1Delete,

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

			"node_count": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: false,
			},

			"nodes_to_remove": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceKubernetesClusterResizeV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	kubernetesClient, err := config.ContainerInfraV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud Kubernetes client: %s", err)
	}

	kubernetesClient.Microversion = kubernetesV1NodeGroupMinMicroversion

	clusterIDOrName := d.Get("cluster_id").(string)
	nodegroupInput := d.Get("nodegroup_id").(string)

	// Extract nodegroup_id from "cluster_id/nodegroup_id" format if necessary
	// This allows using nhncloud_kubernetes_nodegroup_v1.resource_name.id directly
	nodegroupID := extractNodeGroupID(nodegroupInput)

	// Build resize options
	resizeOpts := clusters.ResizeOpts{
		NodeGroup: nodegroupID,
	}

	if v, ok := d.GetOk("node_count"); ok {
		nodeCount := v.(int)
		resizeOpts.NodeCount = &nodeCount
	}

	var nodesToRemove []string
	if raw, ok := d.GetOk("nodes_to_remove"); ok {
		for _, v := range raw.([]interface{}) {
			nodesToRemove = append(nodesToRemove, v.(string))
		}
	}
	resizeOpts.NodesToRemove = nodesToRemove

	log.Printf("[DEBUG] Resizing NKS cluster %s to %d nodes", clusterIDOrName, resizeOpts.NodeCount)

	cluster, err := clusters.Resize(kubernetesClient, clusterIDOrName, resizeOpts).Extract()
	if err != nil {
		return diag.Errorf("Error resizing NKS cluster %s: %s", clusterIDOrName, err)
	}

	d.SetId(cluster.UUID)

	log.Printf("[DEBUG] Waiting for NKS cluster %s resize to complete", cluster.UUID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"UPDATE_IN_PROGRESS"},
		Target:     []string{"UPDATE_COMPLETE"},
		Refresh:    kubernetesClusterV1StateRefreshFunc(kubernetesClient, cluster.UUID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      30 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for NKS cluster %s resize to complete: %s", cluster.UUID, err)
	}

	return resourceKubernetesClusterResizeV1Read(ctx, d, meta)
}

func resourceKubernetesClusterResizeV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	kubernetesClient, err := config.ContainerInfraV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud Kubernetes client: %s", err)
	}

	// Set container-infra API microversion to latest for NKS compatibility
	kubernetesClient.Microversion = kubernetesV1NodeGroupMinMicroversion

	cluster, err := clusters.Get(kubernetesClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving NKS cluster"))
	}

	log.Printf("[DEBUG] Retrieved NKS cluster %s: %+v", d.Id(), cluster)

	d.Set("uuid", cluster.UUID)
	d.Set("status", cluster.Status)
	d.Set("resized_at", cluster.UpdatedAt)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceKubernetesClusterResizeV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	kubernetesClient, err := config.ContainerInfraV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud Kubernetes client: %s", err)
	}

	kubernetesClient.Microversion = kubernetesV1NodeGroupMinMicroversion

	if d.HasChange("node_count") || d.HasChange("nodes_to_remove") {
		nodegroupInput := d.Get("nodegroup_id").(string)
		nodegroupID := extractNodeGroupID(nodegroupInput)

		// Build resize options
		resizeOpts := clusters.ResizeOpts{
			NodeGroup: nodegroupID,
		}

		if v, ok := d.GetOk("node_count"); ok {
			nodeCount := v.(int)
			resizeOpts.NodeCount = &nodeCount
		}

		var nodesToRemove []string
		if raw, ok := d.GetOk("nodes_to_remove"); ok {
			for _, v := range raw.([]interface{}) {
				nodesToRemove = append(nodesToRemove, v.(string))
			}
		}
		resizeOpts.NodesToRemove = nodesToRemove

		log.Printf("[DEBUG] Updating NKS cluster %s nodegroup %s to %d nodes", d.Id(), nodegroupID, resizeOpts.NodeCount)

		_, err = clusters.Resize(kubernetesClient, d.Id(), resizeOpts).Extract()
		if err != nil {
			return diag.Errorf("Error resizing NKS cluster %s: %s", d.Id(), err)
		}

		log.Printf("[DEBUG] Waiting for NKS cluster %s resize to complete", d.Id())

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"UPDATE_IN_PROGRESS"},
			Target:     []string{"UPDATE_COMPLETE"},
			Refresh:    kubernetesClusterV1StateRefreshFunc(kubernetesClient, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      30 * time.Second,
			MinTimeout: 10 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("Error waiting for NKS cluster %s resize to complete: %s", d.Id(), err)
		}
	}

	return resourceKubernetesClusterResizeV1Read(ctx, d, meta)
}

func resourceKubernetesClusterResizeV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// This is a one-time operation resource, so delete just removes it from state
	log.Printf("[DEBUG] Removing NKS cluster resize resource %s from state", d.Id())
	return nil
}

// extractNodeGroupID extracts the nodegroup ID from either a simple ID or "cluster_id/nodegroup_id" format
// This allows users to reference nhncloud_kubernetes_nodegroup_v1 resources directly
func extractNodeGroupID(input string) string {
	parts := strings.Split(input, "/")
	if len(parts) >= 2 {
		// Return the last part if in "cluster_id/nodegroup_id" format
		return parts[len(parts)-1]
	}

	return input
}
