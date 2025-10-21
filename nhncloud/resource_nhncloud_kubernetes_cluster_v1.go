package nhncloud

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud"
	"github.com/nhn-cloud/nhncloud.gophercloud/nhncloud/kubernetes/v1/clusters"
	"github.com/nhn-cloud/nhncloud.gophercloud/nhncloud/kubernetes/v1/nodegroups"
)

func resourceKubernetesClusterV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKubernetesClusterV1Create,
		ReadContext:   resourceKubernetesClusterV1Read,
		UpdateContext: resourceKubernetesClusterV1Update,
		DeleteContext: resourceKubernetesClusterV1Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"project_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
			},

			"user_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
			},

			"created_at": {
				Type:     schema.TypeString,
				ForceNew: false,
				Computed: true,
			},

			"updated_at": {
				Type:     schema.TypeString,
				ForceNew: false,
				Computed: true,
			},

			"api_address": {
				Type:     schema.TypeString,
				ForceNew: false,
				Computed: true,
			},

			"coe_version": {
				Type:     schema.TypeString,
				ForceNew: false,
				Computed: true,
			},

			"cluster_template_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_MAGNUM_CLUSTER_TEMPLATE", nil),
			},

			"container_version": {
				Type:     schema.TypeString,
				ForceNew: false,
				Computed: true,
			},

			"create_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"docker_volume_size": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"keypair": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"node_count": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
				Default:  1,
			},

			"node_addresses": {
				Type:     schema.TypeList,
				ForceNew: false,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"stack_id": {
				Type:     schema.TypeString,
				ForceNew: false,
				Computed: true,
			},

			"fixed_network": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"fixed_subnet": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"kubeconfig": {
				Type:      schema.TypeMap,
				Computed:  true,
				Sensitive: true,
				Elem:      &schema.Schema{Type: schema.TypeString},
			},

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

			"flavor_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
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
		},
	}
}

func resourceKubernetesClusterV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	kubernetesClient, err := config.ContainerInfraV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud kubernetes client: %s", err)
	}

	rawLabels := d.Get("labels").(map[string]interface{})
	labels, err := expandKubernetesV1LabelsMap(rawLabels)
	if err != nil {
		return diag.FromErr(err)
	}

	createOpts := clusters.CreateOpts{
		ClusterTemplateID: d.Get("cluster_template_id").(string),
		FlavorID:          d.Get("flavor_id").(string),
		Keypair:           d.Get("keypair").(string),
		Labels:            labels,
		Name:              d.Get("name").(string),
		FixedNetwork:      d.Get("fixed_network").(string),
		FixedSubnet:       d.Get("fixed_subnet").(string),
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

	s, err := clusters.Create(kubernetesClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating nhncloud_kubernetes_cluster_v1: %s", err)
	}

	// Store the Cluster ID.
	d.SetId(s.UUID)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"CREATE_IN_PROGRESS"},
		Target:       []string{"CREATE_COMPLETE"},
		Refresh:      kubernetesClusterV1StateRefreshFunc(kubernetesClient, s.UUID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        1 * time.Minute,
		PollInterval: 20 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"Error waiting for nhncloud_kubernetes_cluster_v1 %s to become ready: %s", s, err)
	}

	log.Printf("[DEBUG] Created nhncloud_kubernetes_cluster_v1 %s", s)

	return resourceKubernetesClusterV1Read(ctx, d, meta)
}

func getKubernetesDefaultNodegroupNodeCount(kubernetesClient *gophercloud.ServiceClient, clusterID string) (int, error) {
	kubernetesClient.Microversion = kubernetesV1NodeGroupMinMicroversion
	listOpts := nodegroups.ListOpts{}

	allPages, err := nodegroups.List(kubernetesClient, clusterID, listOpts).AllPages()
	if err != nil {
		return 0, err
	}

	ngs, err := nodegroups.ExtractNodegroups(allPages)
	if err != nil {
		return 0, err
	}

	for _, ng := range ngs {
		if ng.IsDefault && ng.Role != "master" {
			return ng.NodeCount, nil
		}
	}

	return 0, fmt.Errorf("Default worker nodegroup not found")
}

func resourceKubernetesClusterV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	kubernetesClient, err := config.ContainerInfraV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud kubernetes client: %s", err)
	}

	s, err := clusters.Get(kubernetesClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving nhncloud_kubernetes_cluster_v1"))
	}

	log.Printf("[DEBUG] Retrieved nhncloud_kubernetes_cluster_v1 %s: %#v", d.Id(), s)

	labels := s.Labels
	if err := d.Set("labels", labels); err != nil {
		return diag.Errorf("Unable to set nhncloud_kubernetes_cluster_v1 labels: %s", err)
	}

	nodeCount, err := getKubernetesDefaultNodegroupNodeCount(kubernetesClient, d.Id())
	if err != nil {
		log.Printf("[DEBUG] Can't retrieve node_count of the default worker node group %s: %s", d.Id(), err)

		nodeCount = s.NodeCount
	}

	d.Set("region", GetRegion(d, config))
	d.Set("name", s.Name)
	d.Set("project_id", s.ProjectID)
	d.Set("user_id", s.UserID)
	d.Set("api_address", s.APIAddress)
	d.Set("coe_version", s.COEVersion)
	d.Set("cluster_template_id", s.ClusterTemplateID)
	d.Set("container_version", s.ContainerVersion)
	d.Set("create_timeout", s.CreateTimeout)
	d.Set("docker_volume_size", s.DockerVolumeSize)
	d.Set("flavor", s.FlavorID)
	d.Set("keypair", s.KeyPair)
	d.Set("node_count", nodeCount)
	d.Set("node_addresses", s.NodeAddresses)
	d.Set("stack_id", s.StackID)
	d.Set("fixed_network", s.FixedNetwork)
	d.Set("fixed_subnet", s.FixedSubnet)
	d.Set("uuid", s.UUID)
	d.Set("status", s.Status)
	d.Set("status_reason", s.StatusReason)

	if err := d.Set("created_at", s.CreatedAt.Format(time.RFC3339)); err != nil {
		log.Printf("[DEBUG] Unable to set nhncloud_kubernetes_cluster_v1 created_at: %s", err)
	}
	if err := d.Set("updated_at", s.UpdatedAt.Format(time.RFC3339)); err != nil {
		log.Printf("[DEBUG] Unable to set nhncloud_kubernetes_cluster_v1 updated_at: %s", err)
	}

	return nil
}

func resourceKubernetesClusterV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	kubernetesClient, err := config.ContainerInfraV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud kubernetes client: %s", err)
	}

	updateOpts := []clusters.UpdateOptsBuilder{}

	if d.HasChange("node_count") {
		nodeCount := d.Get("node_count").(int)
		updateOpts = append(updateOpts, clusters.UpdateOpts{
			Op:    clusters.ReplaceOp,
			Path:  strings.Join([]string{"/", "node_count"}, ""),
			Value: nodeCount,
		})
	}

	if len(updateOpts) > 0 {
		log.Printf(
			"[DEBUG] Updating nhncloud_kubernetes_cluster_v1 %s with options: %#v", d.Id(), updateOpts)

		_, err = clusters.Update(kubernetesClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating nhncloud_kubernetes_cluster_v1 %s: %s", d.Id(), err)
		}

		stateConf := &resource.StateChangeConf{
			Pending:      []string{"UPDATE_IN_PROGRESS"},
			Target:       []string{"UPDATE_COMPLETE"},
			Refresh:      kubernetesClusterV1StateRefreshFunc(kubernetesClient, d.Id()),
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			Delay:        1 * time.Minute,
			PollInterval: 20 * time.Second,
		}
		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf(
				"Error waiting for nhncloud_kubernetes_cluster_v1 %s to become updated: %s", d.Id(), err)
		}
	}
	return resourceKubernetesClusterV1Read(ctx, d, meta)
}

func resourceKubernetesClusterV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	kubernetesClient, err := config.ContainerInfraV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud kubernetes client: %s", err)
	}

	if err := clusters.Delete(kubernetesClient, d.Id()).ExtractErr(); err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting nhncloud_kubernetes_cluster_v1"))
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"DELETE_IN_PROGRESS"},
		Target:       []string{"DELETE_COMPLETE"},
		Refresh:      kubernetesClusterV1StateRefreshFunc(kubernetesClient, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        30 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"Error waiting for nhncloud_kubernetes_cluster_v1 %s to become deleted: %s", d.Id(), err)
	}

	return nil
}
