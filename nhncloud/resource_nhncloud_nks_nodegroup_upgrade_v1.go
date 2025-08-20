package nhncloud

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/nhn-cloud/nhncloud.gophercloud/nhncloud/nks/v1/nodegroups"
)

func resourceNKSNodegroupUpgradeV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNKSNodegroupUpgradeV1Create,
		ReadContext:   resourceNKSNodegroupUpgradeV1Read,
		DeleteContext: resourceNKSNodegroupUpgradeV1Delete,

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

			"version": {
				Type:     schema.TypeString,
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
						"num_buffer_nodes": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1,
						},
						"num_max_unavailable_nodes": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1,
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

			"upgraded_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNKSNodegroupUpgradeV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	containerInfraClient, err := config.ContainerInfraV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud container-infra client: %s", err)
	}

	clusterIDOrName := d.Get("cluster_id").(string)
	nodegroupIDOrName := d.Get("nodegroup_id").(string)

	// Build upgrade options
	upgradeOpts := nodegroups.UpgradeOpts{
		Version: d.Get("version").(string),
	}

	// Set upgrade options
	if v, ok := d.GetOk("options"); ok {
		optionsList := v.([]interface{})
		if len(optionsList) > 0 {
			optionsMap := optionsList[0].(map[string]interface{})
			upgradeOpts.Options = &nodegroups.UpgradeOptions{
				NumBufferNodes:         optionsMap["num_buffer_nodes"].(int),
				NumMaxUnavailableNodes: optionsMap["num_max_unavailable_nodes"].(int),
			}
		}
	}

	log.Printf("[DEBUG] Upgrading NKS nodegroup %s in cluster %s to version %s", nodegroupIDOrName, clusterIDOrName, upgradeOpts.Version)

	nodegroup, err := nodegroups.Upgrade(containerInfraClient, clusterIDOrName, nodegroupIDOrName, upgradeOpts).Extract()
	if err != nil {
		return diag.Errorf("Error upgrading NKS nodegroup %s in cluster %s: %s", nodegroupIDOrName, clusterIDOrName, err)
	}

	d.SetId(nodegroup.UUID)

	log.Printf("[DEBUG] Waiting for NKS nodegroup %s upgrade to complete", nodegroup.UUID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"UPDATE_IN_PROGRESS"},
		Target:     []string{"UPDATE_COMPLETE"},
		Refresh:    nksNodegroupV1StateRefreshFunc(containerInfraClient, clusterIDOrName, nodegroup.UUID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      30 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for NKS nodegroup %s upgrade to complete: %s", nodegroup.UUID, err)
	}

	return resourceNKSNodegroupUpgradeV1Read(ctx, d, meta)
}

func resourceNKSNodegroupUpgradeV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	containerInfraClient, err := config.ContainerInfraV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud container-infra client: %s", err)
	}

	clusterIDOrName := d.Get("cluster_id").(string)

	nodegroup, err := nodegroups.Get(containerInfraClient, clusterIDOrName, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving NKS nodegroup"))
	}

	log.Printf("[DEBUG] Retrieved NKS nodegroup %s: %+v", d.Id(), nodegroup)

	d.Set("uuid", nodegroup.UUID)
	d.Set("status", nodegroup.Status)
	d.Set("upgraded_at", nodegroup.UpdatedAt)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceNKSNodegroupUpgradeV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// This is a one-time operation resource, so delete just removes it from state
	log.Printf("[DEBUG] Removing NKS nodegroup upgrade resource %s from state", d.Id())
	return nil
}
