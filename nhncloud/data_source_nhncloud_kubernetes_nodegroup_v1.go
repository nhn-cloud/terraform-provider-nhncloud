package nhncloud

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/nhn-cloud/nhncloud.gophercloud/nhncloud/kubernetes/v1/nodegroups"
)

func dataSourceKubernetesNodeGroupV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKubernetesNodeGroupRead,
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

			"project_id": {
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

			"docker_volume_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"labels": {
				Type:     schema.TypeMap,
				Computed: true,
			},

			"role": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"node_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"min_node_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"max_node_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"image_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"flavor_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			// NKS specific fields
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"status_reason": {
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

func dataSourceKubernetesNodeGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	kubernetesClient, err := config.ContainerInfraV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud kubernetes client: %s", err)
	}

	kubernetesClient.Microversion = kubernetesV1NodeGroupMinMicroversion

	clusterID := d.Get("cluster_id").(string)

	var nodegroupIDOrName string
	if uuid, ok := d.GetOk("uuid"); ok {
		nodegroupIDOrName = uuid.(string)
	} else if name, ok := d.GetOk("name"); ok {
		nodegroupIDOrName = name.(string)
	} else {
		return diag.Errorf("Either 'uuid' or 'name' must be specified")
	}

	nodeGroup, err := nodegroups.Get(kubernetesClient, clusterID, nodegroupIDOrName).Extract()
	if err != nil {
		return diag.Errorf("Error getting nhncloud_kubernetes_nodegroup_v1 %s: %s", nodegroupIDOrName, err)
	}

	d.SetId(nodeGroup.UUID)

	d.Set("project_id", nodeGroup.ProjectID)
	d.Set("docker_volume_size", nodeGroup.DockerVolumeSize)
	d.Set("role", nodeGroup.Role)
	d.Set("node_count", nodeGroup.NodeCount)
	d.Set("min_node_count", nodeGroup.MinNodeCount)
	d.Set("max_node_count", nodeGroup.MaxNodeCount)
	d.Set("image", nodeGroup.ImageID)
	d.Set("flavor", nodeGroup.FlavorID)
	d.Set("image_id", nodeGroup.ImageID)
	d.Set("flavor_id", nodeGroup.FlavorID)
	d.Set("uuid", nodeGroup.UUID)
	d.Set("name", nodeGroup.Name)
	d.Set("status", nodeGroup.Status)
	d.Set("status_reason", nodeGroup.StatusReason)
	d.Set("version", nodeGroup.Version)

	if err := d.Set("labels", nodeGroup.Labels); err != nil {
		log.Printf("[DEBUG] Unable to set labels for nhncloud_kubernetes_nodegroup_v1 %s: %s", nodeGroup.UUID, err)
	}
	if err := d.Set("created_at", nodeGroup.CreatedAt.Format(time.RFC3339)); err != nil {
		log.Printf("[DEBUG] Unable to set created_at for nhncloud_kubernetes_nodegroup_v1 %s: %s", nodeGroup.UUID, err)
	}
	if err := d.Set("updated_at", nodeGroup.UpdatedAt.Format(time.RFC3339)); err != nil {
		log.Printf("[DEBUG] Unable to set updated_at for nhncloud_kubernetes_nodegroup_v1 %s: %s", nodeGroup.UUID, err)
	}

	d.Set("region", GetRegion(d, config))

	return nil
}
