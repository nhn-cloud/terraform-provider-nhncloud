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

	"github.com/nhn-cloud/nhncloud.gophercloud/nhncloud/kubernetes/v1/nodegroups"
)

func resourceKubernetesNodeGroupV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKubernetesNodeGroupV1Create,
		ReadContext:   resourceKubernetesNodeGroupV1Read,
		UpdateContext: resourceKubernetesNodeGroupV1Update,
		DeleteContext: resourceKubernetesNodeGroupV1Delete,
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

			"project_id": {
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

			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Computed: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					labels := v.(map[string]interface{})
					requiredLabels := []string{
						"availability_zone",
						"boot_volume_type",
						"boot_volume_size",
						"ca_enable",
					}
					for _, key := range requiredLabels {
						if _, exists := labels[key]; !exists {
							errors = append(errors, fmt.Errorf("required label '%s' is missing in %s", key, k))
						}
					}
					return
				},
			},

			"node_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"image_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"flavor_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
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

			"version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			// Nodegroup upgrade options
			"num_max_unavailable_nodes": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},

			"num_buffer_nodes": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
		},
	}
}

func resourceKubernetesNodeGroupV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	kubernetesClient, err := config.ContainerInfraV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud kubernetes client: %s", err)
	}

	kubernetesClient.Microversion = kubernetesV1NodeGroupMinMicroversion

	rawLabels := d.Get("labels").(map[string]interface{})
	labels, err := expandKubernetesV1LabelsMap(rawLabels)
	if err != nil {
		return diag.FromErr(err)
	}

	createOpts := nodegroups.CreateOpts{
		Name:     d.Get("name").(string),
		Labels:   labels,
		ImageID:  d.Get("image_id").(string),
		FlavorID: d.Get("flavor_id").(string),
	}

	var nodeCount int
	if v, ok := d.GetOk("node_count"); ok {
		nodeCount = v.(int)
	} else {
		nodeCount = 1
	}

	createOpts.NodeCount = &nodeCount
	if nodeCount == 0 {
		kubernetesClient.Microversion = kubernetesV1ZeroNodeCountMicroversion
	}

	log.Printf("[DEBUG] nhncloud_kubernetes_nodegroup_v1 create options: %#v", createOpts)

	clusterID := d.Get("cluster_id").(string)
	nodeGroup, err := nodegroups.Create(kubernetesClient, clusterID, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating nhncloud_kubernetes_nodegroup_v1: %s", err)
	}

	id := fmt.Sprintf("%s/%s", clusterID, nodeGroup.UUID)
	d.SetId(id)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"CREATE_IN_PROGRESS"},
		Target:       []string{"CREATE_COMPLETE"},
		Refresh:      kubernetesNodeGroupV1StateRefreshFunc(kubernetesClient, clusterID, nodeGroup.UUID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        1 * time.Minute,
		PollInterval: 20 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"Error waiting for nhncloud_kubernetes_nodegroup_v1 %s to become ready: %s", nodeGroup.UUID, err)
	}

	log.Printf("[DEBUG] Created nhncloud_kubernetes_nodegroup_v1 %s", nodeGroup.UUID)

	return resourceKubernetesNodeGroupV1Read(ctx, d, meta)
}

func resourceKubernetesNodeGroupV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	kubernetesClient, err := config.ContainerInfraV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud kubernetes client: %s", err)
	}

	kubernetesClient.Microversion = kubernetesV1NodeGroupMinMicroversion

	clusterID, nodeGroupID, err := parseNodeGroupID(d.Id())
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error parsing ID of nhncloud_kubernetes_nodegroup_v1"))
	}

	nodeGroup, err := nodegroups.Get(kubernetesClient, clusterID, nodeGroupID).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving nhncloud_kubernetes_nodegroup_v1"))
	}

	log.Printf("[DEBUG] Retrieved nhncloud_kubernetes_nodegroup_v1 %s: %#v", d.Id(), nodeGroup)

	apiLabels := nodeGroup.Labels
	rawConfig := d.GetRawConfig()
	if !rawConfig.IsNull() && rawConfig.Type().HasAttribute("labels") {
		configLabelsAttr := rawConfig.GetAttr("labels")
		if configLabelsAttr.IsKnown() && !configLabelsAttr.IsNull() &&
			(configLabelsAttr.Type().IsObjectType() || configLabelsAttr.Type().IsMapType()) {

			filteredLabels := make(map[string]string)
			for key, val := range configLabelsAttr.AsValueMap() {
				if val.IsNull() || !val.IsKnown() {
					continue
				}
				filteredLabels[key] = val.AsString()
			}

			for key, apiVal := range apiLabels {
				if _, existsInConfig := filteredLabels[key]; existsInConfig {
					filteredLabels[key] = flattenKubernetesV1LabelValue(apiVal)
				}
			}

			if err := d.Set("labels", filteredLabels); err != nil {
				return diag.Errorf("Unable to set labels: %s", err)
			}
		}
	}

	d.Set("cluster_id", clusterID)
	d.Set("region", GetRegion(d, config))
	d.Set("name", nodeGroup.Name)
	d.Set("project_id", nodeGroup.ProjectID)
	d.Set("role", nodeGroup.Role)
	d.Set("node_count", nodeGroup.NodeCount)
	d.Set("min_node_count", nodeGroup.MinNodeCount)
	d.Set("max_node_count", nodeGroup.MaxNodeCount)
	d.Set("image_id", nodeGroup.ImageID)
	d.Set("flavor_id", nodeGroup.FlavorID)
	d.Set("docker_volume_size", nodeGroup.DockerVolumeSize)
	d.Set("uuid", nodeGroup.UUID)
	d.Set("status", nodeGroup.Status)
	d.Set("status_reason", nodeGroup.StatusReason)
	d.Set("version", nodeGroup.Version)

	if err := d.Set("created_at", nodeGroup.CreatedAt.Format(time.RFC3339)); err != nil {
		log.Printf("[DEBUG] Unable to set nhncloud_kubernetes_nodegroup_v1 created_at: %s", err)
	}
	if err := d.Set("updated_at", nodeGroup.UpdatedAt.Format(time.RFC3339)); err != nil {
		log.Printf("[DEBUG] Unable to set nhncloud_kubernetes_nodegroup_v1 updated_at: %s", err)
	}

	return nil
}

func resourceKubernetesNodeGroupV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	kubernetesClient, err := config.ContainerInfraV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud kubernetes client: %s", err)
	}

	kubernetesClient.Microversion = kubernetesV1NodeGroupMinMicroversion

	clusterID, nodeGroupID, err := parseNodeGroupID(d.Id())
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error parsing ID of nhncloud_kubernetes_nodegroup_v1"))
	}

	updateOpts := []nodegroups.UpdateOptsBuilder{}

	if d.HasChange("min_node_count") {
		minNodeCount := d.Get("min_node_count").(int)
		updateOpts = kubernetesNodeGroupV1AppendUpdateOpts(
			updateOpts, "min_node_count", minNodeCount)
	}

	if d.HasChange("max_node_count") {
		maxNodeCount := d.Get("max_node_count").(int)
		updateOpts = kubernetesNodeGroupV1AppendUpdateOpts(
			updateOpts, "max_node_count", maxNodeCount)
	}

	if len(updateOpts) > 0 {
		log.Printf(
			"[DEBUG] Updating nhncloud_kubernetes_nodegroup_v1 %s with options: %#v", d.Id(), updateOpts)

		_, err = nodegroups.Update(kubernetesClient, clusterID, nodeGroupID, updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating nhncloud_kubernetes_nodegroup_v1 %s: %s", d.Id(), err)
		}
	}

	return resourceKubernetesNodeGroupV1Read(ctx, d, meta)
}

func resourceKubernetesNodeGroupV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	kubernetesClient, err := config.ContainerInfraV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud kubernetes client: %s", err)
	}

	kubernetesClient.Microversion = kubernetesV1NodeGroupMinMicroversion

	clusterID, nodeGroupID, err := parseNodeGroupID(d.Id())
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error parsing ID of nhncloud_kubernetes_nodegroup_v1"))
	}

	if err := nodegroups.Delete(kubernetesClient, clusterID, nodeGroupID).ExtractErr(); err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting nhncloud_kubernetes_nodegroup_v1"))
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"DELETE_IN_PROGRESS"},
		Target:       []string{"DELETE_COMPLETE"},
		Refresh:      kubernetesNodeGroupV1StateRefreshFunc(kubernetesClient, clusterID, nodeGroupID),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        30 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"Error waiting for nhncloud_kubernetes_nodegroup_v1 %s to become deleted: %s", d.Id(), err)
	}

	return nil
}

func parseNodeGroupID(id string) (string, string, error) {
	idParts := strings.Split(id, "/")
	if len(idParts) < 2 {
		return "", "", fmt.Errorf("Unable to determine nodegroup ID %s", id)
	}

	return idParts[0], idParts[1], nil
}
