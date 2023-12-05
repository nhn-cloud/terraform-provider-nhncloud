package nhncloud

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/nhn/nhncloud.gophercloud/nhncloud/networking/v2/vpcs"
)

func resourceNetworkingVPCV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkingVPCV2Create,
		ReadContext:   resourceNetworkingVPCV2Read,
		UpdateContext: resourceNetworkingVPCV2Update,
		DeleteContext: resourceNetworkingVPCV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
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
				ForceNew: false,
			},

			"cidrv4": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsCIDR,
				},
			},

			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNetworkingVPCV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud networking client: %s", err)
	}

	createOpts := VPCCreateOpts{
		vpcs.CreateOpts{
			Name:     d.Get("name").(string),
			Cidrv4:   d.Get("cidrv4").(string),
			TenantID: d.Get("tenant_id").(string),
		},
	}

	// Declare a finalCreateOpts interface.
	var finalCreateOpts vpcs.CreateOptsBuilder
	finalCreateOpts = createOpts

	log.Printf("[DEBUG] nhncloud_networking_vpc_v2 create options: %#v", finalCreateOpts)
	n, err := vpcs.Create(networkingClient, finalCreateOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating nhncloud_networking_vpc_v2: %s", err)
	}

	log.Printf("[DEBUG] Waiting for nhncloud_networking_vpc_v2 %s to become available.", n.ID)

	// @tc-iaas-compute/1452
	// The only VPC status value is "available".
	// Therefore, other statuses (e.g. Pending) are regarded as empty strings.
	stateConf := &resource.StateChangeConf{
		Pending:    []string{""},
		Target:     []string{"available"},
		Refresh:    resourceNetworkingVPCV2StateRefreshFunc(networkingClient, n.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for nhncloud_networking_vpc_v2 %s to become available: %s", n.ID, err)
	}

	d.SetId(n.ID)

	log.Printf("[DEBUG] Created nhncloud_networking_vpc_v2 %s: %#v", n.ID, n)
	return resourceNetworkingVPCV2Read(ctx, d, meta)
}

func resourceNetworkingVPCV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud networking client: %s", err)
	}

	var network vpcExtendedDetail

	err = vpcs.Get(networkingClient, d.Id()).ExtractInto(&network)
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error getting nhncloud_networking_vpc_v2"))
	}

	log.Printf("[DEBUG] Retrieved nhncloud_networking_vpc_v2 %s: %#v", d.Id(), network)

	d.Set("region", GetRegion(d, config))
	d.Set("name", network.Name)
	d.Set("cidrv4", network.Cidrv4)
	d.Set("id", network.ID)

	return nil
}

func resourceNetworkingVPCV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud networking client: %s", err)
	}

	// Declare finalUpdateOpts interface and basic updateOpts structure.
	var (
		finalUpdateOpts vpcs.UpdateOptsBuilder
		updateOpts      vpcs.UpdateOpts
	)

	// Populate basic updateOpts.
	if d.HasChange("name") {
		name := d.Get("name").(string)
		updateOpts.Name = &name
	}
	if d.HasChange("cidrv4") {
		cidrv4 := d.Get("cidrv4").(string)
		updateOpts.Cidrv4 = &cidrv4
	}

	// Save basic updateOpts into finalUpdateOpts.
	finalUpdateOpts = updateOpts

	log.Printf("[DEBUG] nhncloud_networking_vpc_v2 %s update options: %#v", d.Id(), finalUpdateOpts)
	_, err = vpcs.Update(networkingClient, d.Id(), finalUpdateOpts).Extract()
	if err != nil {
		return diag.Errorf("Error updating nhncloud_networking_vpc_v2 %s: %s", d.Id(), err)
	}

	return resourceNetworkingVPCV2Read(ctx, d, meta)
}

func resourceNetworkingVPCV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud networking client: %s", err)
	}

	if err := vpcs.Delete(networkingClient, d.Id()).ExtractErr(); err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting nhncloud_networking_vpc_v2"))
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"available"},
		Target:     []string{"DELETED"},
		Refresh:    resourceNetworkingVPCV2StateRefreshFunc(networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for nhncloud_networking_vpc_v2 %s to Delete:  %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}
