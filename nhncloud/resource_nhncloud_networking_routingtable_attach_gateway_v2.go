package nhncloud

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/nhn-cloud/nhncloud.gophercloud/nhncloud/networking/v2/routingtables"
)

func resourceNetworkingRoutingtableAttachGatewayV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkingRoutingtableAttachGatewayV2Create,
		ReadContext:   resourceNetworkingRoutingtableAttachGatewayV2Read,
		DeleteContext: resourceNetworkingRoutingtableAttachGatewayV2Delete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"routingtable_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceNetworkingRoutingtableAttachGatewayV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud networking client: %s", err)
	}

	routingtableID := d.Get("routingtable_id").(string)

	attachOpts := routingtables.AttachGatewayOpts{
		GatewayID: d.Get("gateway_id").(string),
	}

	log.Printf("[DEBUG] nhncloud_networking_routingtable_attach_gateway_v2 %s attach gateway options: %#v", routingtableID, attachOpts)
	_, err = routingtables.AttachGateway(networkingClient, routingtableID, attachOpts).Extract()
	if err != nil {
		return diag.Errorf("Error attaching gateway nhncloud_networking_routingtable_attach_gateway_v2 %s: %s", routingtableID, err)
	}
	d.SetId(routingtableID)

	return resourceNetworkingRoutingtableV2Read(ctx, d, meta)
}

func resourceNetworkingRoutingtableAttachGatewayV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud networking client: %s", err)
	}

	var routingtable routingtableExtended

	err = routingtables.Get(networkingClient, d.Id()).ExtractInto(&routingtable)
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error getting nhncloud_networking_routingtable_attach_gateway_v2"))
	}

	log.Printf("[DEBUG] Retrieved nhncloud_networking_routingtable_attach_gateway_v2 %s: %#v", d.Id(), routingtable)

	d.Set("name", routingtable.Name)
	d.Set("id", routingtable.ID)
	d.Set("distributed", routingtable.Distributed)
	d.Set("tenant_id", routingtable.TenantID)
	d.Set("default_table", routingtable.DefaultTable)
	d.Set("state", routingtable.State)
	d.Set("create_time", routingtable.CreateTime)

	return nil
}

func resourceNetworkingRoutingtableAttachGatewayV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud networking client: %s", err)
	}

	log.Printf("[DEBUG] nhncloud_networking_routingtable_attach_gateway_v2 %s", d.Id())
	_, err = routingtables.DetachGateway(networkingClient, d.Id()).Extract()
	if err != nil {
		return diag.Errorf("Error detaching gateway nhncloud_networking_routingtable_attach_gateway_v2 %s: %s", d.Id(), err)
	}

	return resourceNetworkingRoutingtableV2Read(ctx, d, meta)
}
