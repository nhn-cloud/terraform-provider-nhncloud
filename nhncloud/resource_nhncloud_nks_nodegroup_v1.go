package nhncloud

import (
	"context"
	"log"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/nhn-cloud/nhncloud.gophercloud/nhncloud/nks/v1/nodegroups"
)

func resourceNKSNodegroupV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNKSNodegroupV1Create,
		ReadContext:   resourceNKSNodegroupV1Read,
		DeleteContext: resourceNKSNodegroupV1Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
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

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"node_count": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"flavor_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"image_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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

			"status_reason": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNKSNodegroupV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	kubernetesClient, err := config.ContainerInfraV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud Kubernetes client: %s", err)
	}

	// Set container-infra API microversion to latest for NKS compatibility
	kubernetesClient.Microversion = kubernetesV1NodeGroupMinMicroversion

	clusterIDOrName := d.Get("cluster_id").(string)

	// Build create options
	createOpts := nodegroups.CreateOpts{
		Name:      d.Get("name").(string),
		NodeCount: d.Get("node_count").(int),
		FlavorID:  d.Get("flavor_id").(string),
		ImageID:   d.Get("image_id").(string),
	}

	// Set labels
	if v, ok := d.GetOk("labels"); ok {
		labels := make(map[string]interface{})
		for k, val := range v.(map[string]interface{}) {
			labels[k] = val
		}
		createOpts.Labels = labels
	}

	log.Printf("[DEBUG] Creating NKS nodegroup in cluster %s with options: %+v", clusterIDOrName, createOpts)

	nodegroup, err := nodegroups.Create(kubernetesClient, clusterIDOrName, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating NKS nodegroup in cluster %s: %s", clusterIDOrName, err)
	}

	d.SetId(nodegroup.UUID)

	log.Printf("[DEBUG] Waiting for NKS nodegroup %s to be created", nodegroup.UUID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"CREATE_IN_PROGRESS"},
		Target:     []string{"CREATE_COMPLETE"},
		Refresh:    nksNodegroupV1StateRefreshFunc(kubernetesClient, clusterIDOrName, nodegroup.UUID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for NKS nodegroup %s to be created: %s", nodegroup.UUID, err)
	}

	return resourceNKSNodegroupV1Read(ctx, d, meta)
}

func resourceNKSNodegroupV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	kubernetesClient, err := config.ContainerInfraV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud Kubernetes client: %s", err)
	}

	// Set container-infra API microversion to latest for NKS compatibility
	kubernetesClient.Microversion = kubernetesV1NodeGroupMinMicroversion

	clusterIDOrName := d.Get("cluster_id").(string)

	nodegroup, err := nodegroups.Get(kubernetesClient, clusterIDOrName, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving NKS nodegroup"))
	}

	if nodegroup == nil {
		return diag.Errorf("Retrieved NKS nodegroup is nil for ID: %s", d.Id())
	}

	log.Printf("[DEBUG] Retrieved NKS nodegroup %s: %+v", d.Id(), nodegroup)

	d.Set("uuid", nodegroup.UUID)
	d.Set("name", nodegroup.Name)
	d.Set("cluster_id", nodegroup.ClusterID)
	d.Set("node_count", nodegroup.NodeCount)
	d.Set("flavor_id", nodegroup.FlavorID)
	d.Set("image_id", nodegroup.ImageID)
	d.Set("status", nodegroup.Status)
	d.Set("status_reason", nodegroup.StatusReason)
	d.Set("created_at", nodegroup.CreatedAt)
	d.Set("updated_at", nodegroup.UpdatedAt)
	d.Set("version", nodegroup.Version)
	d.Set("region", GetRegion(d, config))

	// Set labels
	if nodegroup.Labels != nil {
		labelMap := make(map[string]interface{})
		for k, v := range nodegroup.Labels {
			if str, ok := v.(string); ok {
				labelMap[k] = str
			}
		}
		d.Set("labels", labelMap)
	}

	return nil
}

func resourceNKSNodegroupV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	kubernetesClient, err := config.ContainerInfraV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud Kubernetes client: %s", err)
	}

	// Set container-infra API microversion to latest for NKS compatibility
	kubernetesClient.Microversion = kubernetesV1NodeGroupMinMicroversion

	clusterIDOrName := d.Get("cluster_id").(string)

	log.Printf("[DEBUG] Deleting NKS nodegroup %s from cluster %s", d.Id(), clusterIDOrName)

	err = nodegroups.Delete(kubernetesClient, clusterIDOrName, d.Id()).ExtractErr()
	if err != nil {
		return diag.Errorf("Error deleting NKS nodegroup %s: %s", d.Id(), err)
	}

	log.Printf("[DEBUG] Waiting for NKS nodegroup %s to be deleted", d.Id())

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"DELETE_IN_PROGRESS"},
		Target:     []string{"DELETED", "DELETE_COMPLETE"},
		Refresh:    nksNodegroupV1StateRefreshFunc(kubernetesClient, clusterIDOrName, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for NKS nodegroup %s to be deleted: %s", d.Id(), err)
	}

	return nil
}

func nksNodegroupV1StateRefreshFunc(client *gophercloud.ServiceClient, clusterIDOrName, nodegroupID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		nodegroup, err := nodegroups.Get(client, clusterIDOrName, nodegroupID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return nodegroup, "DELETED", nil
			}
			return nil, "", err
		}

		return nodegroup, nodegroup.Status, nil
	}
}
