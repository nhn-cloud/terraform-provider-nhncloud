package nhncloud

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/nhn-cloud/nhncloud.gophercloud/nhncloud/nks/v1/clusters"
)

func resourceNKSClusterResizeV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNKSClusterResizeV1Create,
		ReadContext:   resourceNKSClusterResizeV1Read,
		DeleteContext: resourceNKSClusterResizeV1Delete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
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
				ForceNew: true,
			},

			"options": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nodes_to_remove": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},

			// computed-only
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"resized_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNKSClusterResizeV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	kubernetesClient, err := config.ContainerInfraV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud Kubernetes client: %s", err)
	}

	// Set container-infra API microversion to latest for NKS compatibility
	kubernetesClient.Microversion = kubernetesV1NodeGroupMinMicroversion

	clusterIDOrName := d.Get("cluster_id").(string)

	// Build resize options
	resizeOpts := clusters.ResizeOpts{
		NodeCount: d.Get("node_count").(int),
		NodeGroup: d.Get("nodegroup_id").(string),
	}

	// Set resize options
	if v, ok := d.GetOk("options"); ok {
		optionsList := v.([]interface{})
		if len(optionsList) > 0 {
			optionsMap := optionsList[0].(map[string]interface{})
			if nodesToRemove, ok := optionsMap["nodes_to_remove"]; ok {
				nodeList := nodesToRemove.([]interface{})
				resizeOpts.Options = &clusters.ResizeOptions{
					NodesToRemove: make([]string, len(nodeList)),
				}
				for i, node := range nodeList {
					resizeOpts.Options.NodesToRemove[i] = node.(string)
				}
			}
		}
	}

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
		Refresh:    nksClusterV1StateRefreshFunc(kubernetesClient, cluster.UUID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      30 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for NKS cluster %s resize to complete: %s", cluster.UUID, err)
	}

	return resourceNKSClusterResizeV1Read(ctx, d, meta)
}

func resourceNKSClusterResizeV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func resourceNKSClusterResizeV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// This is a one-time operation resource, so delete just removes it from state
	log.Printf("[DEBUG] Removing NKS cluster resize resource %s from state", d.Id())
	return nil
}
