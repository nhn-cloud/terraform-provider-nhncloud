package nhncloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/nhn-cloud/nhncloud.gophercloud/nhncloud/nks/v1/clusters"
)

func dataSourceNKSClusterV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNKSClusterV1Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
			"cluster_template_id": {
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

			"keypair": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"node_count": {
				Type:     schema.TypeInt,
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

			"labels": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"api_ep_ipacl": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"action": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipacl_targets": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cidr_address": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
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

func dataSourceNKSClusterV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	containerInfraClient, err := config.ContainerInfraV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud container-infra client: %s", err)
	}

	var clusterIDOrName string
	if uuid, ok := d.GetOk("uuid"); ok {
		clusterIDOrName = uuid.(string)
	} else if name, ok := d.GetOk("name"); ok {
		clusterIDOrName = name.(string)
	} else {
		return diag.Errorf("Either 'uuid' or 'name' must be specified")
	}

	log.Printf("[DEBUG] Retrieving NKS cluster: %s", clusterIDOrName)

	cluster, err := clusters.Get(containerInfraClient, clusterIDOrName).Extract()
	if err != nil {
		return diag.Errorf("Unable to retrieve NKS cluster %s: %s", clusterIDOrName, err)
	}

	log.Printf("[DEBUG] Retrieved NKS cluster %s: %+v", clusterIDOrName, cluster)

	d.SetId(cluster.UUID)
	d.Set("uuid", cluster.UUID)
	d.Set("name", cluster.Name)
	d.Set("cluster_template_id", cluster.ClusterTemplateID)
	d.Set("fixed_network", cluster.FixedNetwork)
	d.Set("fixed_subnet", cluster.FixedSubnet)
	d.Set("flavor_id", cluster.FlavorID)
	d.Set("keypair", cluster.Keypair)
	d.Set("node_count", cluster.NodeCount)
	d.Set("status", cluster.Status)
	d.Set("status_reason", cluster.StatusReason)
	d.Set("created_at", cluster.CreatedAt)
	d.Set("updated_at", cluster.UpdatedAt)
	d.Set("region", GetRegion(d, config))

	// Set labels
	if cluster.Labels != nil {
		labelMap := make(map[string]interface{})
		for k, v := range cluster.Labels {
			if str, ok := v.(string); ok {
				labelMap[k] = str
			}
		}
		d.Set("labels", labelMap)
	}

	// Set API endpoint IP ACL
	if cluster.APIEndpointIPACL != nil {
		ipacl := []map[string]interface{}{
			{
				"enable": cluster.APIEndpointIPACL.Enable,
				"action": cluster.APIEndpointIPACL.Action,
			},
		}

		if cluster.APIEndpointIPACL.IPACLTargets != nil {
			targets := make([]map[string]interface{}, len(cluster.APIEndpointIPACL.IPACLTargets))
			for i, target := range cluster.APIEndpointIPACL.IPACLTargets {
				targets[i] = map[string]interface{}{
					"cidr_address": target.CidrAddress,
					"description":  target.Description,
				}
			}
			ipacl[0]["ipacl_targets"] = targets
		}

		d.Set("api_ep_ipacl", ipacl)
	}

	// Set addons
	if cluster.Addons != nil {
		addonList := make([]map[string]interface{}, len(cluster.Addons))
		for i, addon := range cluster.Addons {
			addonMap := map[string]interface{}{
				"name":    addon.Name,
				"version": addon.Version,
			}

			if addon.Options != nil {
				optionsMap := make(map[string]interface{})
				for k, v := range addon.Options {
					if str, ok := v.(string); ok {
						optionsMap[k] = str
					}
				}
				addonMap["options"] = optionsMap
			}

			addonList[i] = addonMap
		}
		d.Set("addons", addonList)
	}

	return nil
}
