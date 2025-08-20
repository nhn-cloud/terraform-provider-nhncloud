package nhncloud

import (
	"context"
	"log"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/nhn-cloud/nhncloud.gophercloud/nhncloud/nks/v1/clusters"
)

func resourceNKSClusterV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNKSClusterV1Create,
		ReadContext:   resourceNKSClusterV1Read,
		DeleteContext: resourceNKSClusterV1Delete,
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

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"cluster_template_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"fixed_network": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"fixed_subnet": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"flavor_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"keypair": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"node_count": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"api_ep_ipacl": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable": {
							Type:     schema.TypeString,
							Required: true,
						},
						"action": {
							Type:     schema.TypeString,
							Required: true,
						},
						"ipacl_targets": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cidr_address": {
										Type:     schema.TypeString,
										Required: true,
									},
									"description": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},

			"addons": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"version": {
							Type:     schema.TypeString,
							Required: true,
						},
						"options": {
							Type:     schema.TypeMap,
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
		},
	}
}

func resourceNKSClusterV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	containerInfraClient, err := config.ContainerInfraV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud container-infra client: %s", err)
	}

	// Build create options
	createOpts := clusters.CreateOpts{
		Name:              d.Get("name").(string),
		ClusterTemplateID: d.Get("cluster_template_id").(string),
		FixedNetwork:      d.Get("fixed_network").(string),
		FixedSubnet:       d.Get("fixed_subnet").(string),
		FlavorID:          d.Get("flavor_id").(string),
		Keypair:           d.Get("keypair").(string),
		NodeCount:         d.Get("node_count").(int),
	}

	// Set labels
	if v, ok := d.GetOk("labels"); ok {
		labels := make(map[string]interface{})
		for k, val := range v.(map[string]interface{}) {
			labels[k] = val
		}
		createOpts.Labels = labels
	}

	// Set API endpoint IP ACL
	if v, ok := d.GetOk("api_ep_ipacl"); ok {
		ipaclList := v.([]interface{})
		if len(ipaclList) > 0 {
			ipaclMap := ipaclList[0].(map[string]interface{})
			ipacl := &clusters.APIEndpointIPACL{
				Enable: ipaclMap["enable"].(string),
				Action: ipaclMap["action"].(string),
			}

			if targets, ok := ipaclMap["ipacl_targets"]; ok {
				targetList := targets.([]interface{})
				ipacl.IPACLTargets = make([]clusters.IPACLTarget, len(targetList))
				for i, target := range targetList {
					targetMap := target.(map[string]interface{})
					ipacl.IPACLTargets[i] = clusters.IPACLTarget{
						CidrAddress: targetMap["cidr_address"].(string),
						Description: targetMap["description"].(string),
					}
				}
			}

			createOpts.APIEndpointIPACL = ipacl
		}
	}

	// Set addons
	if v, ok := d.GetOk("addons"); ok {
		addonList := v.([]interface{})
		createOpts.Addons = make([]clusters.Addon, len(addonList))
		for i, addon := range addonList {
			addonMap := addon.(map[string]interface{})
			addonOpts := clusters.Addon{
				Name:    addonMap["name"].(string),
				Version: addonMap["version"].(string),
			}

			if options, ok := addonMap["options"]; ok && options != nil {
				optionsMap := make(map[string]interface{})
				for k, val := range options.(map[string]interface{}) {
					optionsMap[k] = val
				}
				addonOpts.Options = optionsMap
			}

			createOpts.Addons[i] = addonOpts
		}
	}

	log.Printf("[DEBUG] Creating NKS cluster with options: %+v", createOpts)

	cluster, err := clusters.Create(containerInfraClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating NKS cluster: %s", err)
	}

	d.SetId(cluster.UUID)

	log.Printf("[DEBUG] Waiting for NKS cluster %s to be created", cluster.UUID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"CREATE_IN_PROGRESS"},
		Target:     []string{"CREATE_COMPLETE"},
		Refresh:    nksClusterV1StateRefreshFunc(containerInfraClient, cluster.UUID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for NKS cluster %s to be created: %s", cluster.UUID, err)
	}

	return resourceNKSClusterV1Read(ctx, d, meta)
}

func resourceNKSClusterV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	containerInfraClient, err := config.ContainerInfraV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud container-infra client: %s", err)
	}

	cluster, err := clusters.Get(containerInfraClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving NKS cluster"))
	}

	log.Printf("[DEBUG] Retrieved NKS cluster %s: %+v", d.Id(), cluster)

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

func resourceNKSClusterV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	containerInfraClient, err := config.ContainerInfraV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud container-infra client: %s", err)
	}

	log.Printf("[DEBUG] Deleting NKS cluster %s", d.Id())

	err = clusters.Delete(containerInfraClient, d.Id()).ExtractErr()
	if err != nil {
		return diag.Errorf("Error deleting NKS cluster %s: %s", d.Id(), err)
	}

	log.Printf("[DEBUG] Waiting for NKS cluster %s to be deleted", d.Id())

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"DELETE_IN_PROGRESS"},
		Target:     []string{"DELETED"},
		Refresh:    nksClusterV1StateRefreshFunc(containerInfraClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for NKS cluster %s to be deleted: %s", d.Id(), err)
	}

	return nil
}

func nksClusterV1StateRefreshFunc(client *gophercloud.ServiceClient, clusterID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		cluster, err := clusters.Get(client, clusterID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return cluster, "DELETED", nil
			}
			return nil, "", err
		}

		return cluster, cluster.Status, nil
	}
}
