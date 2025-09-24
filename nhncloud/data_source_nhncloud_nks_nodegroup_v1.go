package nhncloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/nhn-cloud/nhncloud.gophercloud/nhncloud/nks/v1/nodegroups"
)

func dataSourceNKSNodegroupV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNKSNodegroupV1Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"uuid": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"uuid", "name"},
			},

			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"uuid", "name"},
			},

			// computed-only
			"node_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"flavor_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"image_id": {
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

			"labels": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceNKSNodegroupV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	kubernetesClient, err := config.ContainerInfraV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud Kubernetes client: %s", err)
	}

	// Set container-infra API microversion to latest for NKS compatibility
	kubernetesClient.Microversion = kubernetesV1NodeGroupMinMicroversion

	clusterIDOrName := d.Get("cluster_id").(string)

	var nodegroupIDOrName string
	if uuid, ok := d.GetOk("uuid"); ok {
		nodegroupIDOrName = uuid.(string)
	} else if name, ok := d.GetOk("name"); ok {
		nodegroupIDOrName = name.(string)
	} else {
		return diag.Errorf("Either 'uuid' or 'name' must be specified")
	}

	log.Printf("[DEBUG] Retrieving NKS nodegroup: cluster=%s, nodegroup=%s", clusterIDOrName, nodegroupIDOrName)

	nodegroup, err := nodegroups.Get(kubernetesClient, clusterIDOrName, nodegroupIDOrName).Extract()
	if err != nil {
		return diag.Errorf("Unable to retrieve NKS nodegroup %s in cluster %s: %s", nodegroupIDOrName, clusterIDOrName, err)
	}

	if nodegroup == nil {
		return diag.Errorf("Retrieved NKS nodegroup is nil for %s in cluster %s", nodegroupIDOrName, clusterIDOrName)
	}

	log.Printf("[DEBUG] Retrieved NKS nodegroup %s: %+v", nodegroupIDOrName, nodegroup)

	d.SetId(nodegroup.UUID)
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
