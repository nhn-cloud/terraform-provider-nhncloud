package nhncloud

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/nhn/nhncloud.gophercloud/nhncloud/networking/v2/vpcsubnets"
)

func resourceNetworkingVPCSubnetRoutingtableV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkingVPCSubnetRoutingtableV2Create,
		ReadContext:   resourceNetworkingVPCSubnetRoutingtableV2Read,
		DeleteContext: resourceNetworkingVPCSubnetRoutingtableV2Delete,
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

			"routingtable_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceNetworkingVPCSubnetRoutingtableV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	var (
		attachOpts      vpcsubnets.AttachOpts
		finalAttachOpts vpcsubnets.AttachOptsBuilder
	)

	attachOpts.RoutingtableID = d.Get("routingtable_id").(string)
	finalAttachOpts = attachOpts

	subnetID := d.Get("subnet_id").(string)
	if err = vpcsubnets.Attach(networkingClient, subnetID, finalAttachOpts).ExtractErr(); err != nil {
		return diag.Errorf("Error attaching vpc subnet to routingtable: %s", err)
	}

	d.SetId(subnetID)
	log.Printf("[DEBUG] Attached vpc subnet to routingtable")

	return nil
}

func resourceNetworkingVPCSubnetRoutingtableV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	//@tc-iaas-compute/1452
	// This method should be implemented, but I never think that it would be invoked.
	// So I don't implement business logic, but let it do nothing.
	return nil
}

func resourceNetworkingVPCSubnetRoutingtableV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	if err = vpcsubnets.Detach(networkingClient, d.Id()).ExtractErr(); err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error detaching vpc subnet from routingtable"))
	}

	d.SetId("")
	log.Printf("[DEBUG] Detached vpc subnet from routingtable")

	return nil
}
