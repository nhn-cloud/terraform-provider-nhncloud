package nhncloud

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/nhn-cloud/nhncloud.gophercloud/nhncloud/kubernetes/v1/clusters"
)

func dataSourceKubernetesCluster() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKubernetesClusterRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"name", "uuid"},
			},

			"uuid": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"name", "uuid"},
			},

			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"user_id": {
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

			"api_address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"coe_version": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"cluster_template_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"container_version": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"create_timeout": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"discovery_url": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"docker_volume_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"keypair": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"labels": {
				Type:     schema.TypeMap,
				Computed: true,
			},

			"node_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"node_addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"stack_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"fixed_network": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"fixed_subnet": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"flavor_id": {
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

			"addons": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"options": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceKubernetesClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	kubernetesClient, err := config.ContainerInfraV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud kubernetes client: %s", err)
	}

	var clusterIDOrName string
	if uuid, ok := d.GetOk("uuid"); ok {
		clusterIDOrName = uuid.(string)
	} else if name, ok := d.GetOk("name"); ok {
		clusterIDOrName = name.(string)
	} else {
		return diag.Errorf("Either 'uuid' or 'name' must be specified")
	}

	c, err := clusters.Get(kubernetesClient, clusterIDOrName).Extract()
	if err != nil {
		return diag.Errorf("Error getting nhncloud_kubernetes_cluster_v1 %s: %s", clusterIDOrName, err)
	}

	d.SetId(c.UUID)

	d.Set("project_id", c.ProjectID)
	d.Set("user_id", c.UserID)
	d.Set("api_address", c.APIAddress)
	d.Set("coe_version", c.COEVersion)
	d.Set("cluster_template_id", c.ClusterTemplateID)
	d.Set("container_version", c.ContainerVersion)
	d.Set("create_timeout", c.CreateTimeout)
	d.Set("docker_volume_size", c.DockerVolumeSize)
	d.Set("flavor", c.FlavorID)
	d.Set("keypair", c.KeyPair)
	d.Set("node_count", c.NodeCount)
	d.Set("node_addresses", c.NodeAddresses)
	d.Set("stack_id", c.StackID)
	d.Set("fixed_network", c.FixedNetwork)
	d.Set("fixed_subnet", c.FixedSubnet)
	d.Set("flavor_id", c.FlavorID)
	d.Set("uuid", c.UUID)
	d.Set("name", c.Name)
	d.Set("status", c.Status)
	d.Set("status_reason", c.StatusReason)

	// TODO: Set api_ep_ipacl and addons when the structure is available in gophercloud

	if err := d.Set("labels", c.Labels); err != nil {
		log.Printf("[DEBUG] Unable to set labels for nhncloud_kubernetes_cluster_v1 %s: %s", c.UUID, err)
	}
	if err := d.Set("created_at", c.CreatedAt.Format(time.RFC3339)); err != nil {
		log.Printf("[DEBUG] Unable to set created_at for nhncloud_kubernetes_cluster_v1 %s: %s", c.UUID, err)
	}
	if err := d.Set("updated_at", c.UpdatedAt.Format(time.RFC3339)); err != nil {
		log.Printf("[DEBUG] Unable to set updated_at for nhncloud_kubernetes_cluster_v1 %s: %s", c.UUID, err)
	}

	d.Set("region", GetRegion(d, config))

	return nil
}
