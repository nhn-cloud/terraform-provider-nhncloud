package nhncloud

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/nhn-cloud/nhncloud.gophercloud/nhncloud/networking/v2/routingtables"
)

func resourceNetworkingRoutingtableV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkingRoutingtableV2Create,
		ReadContext:   resourceNetworkingRoutingtableV2Read,
		UpdateContext: resourceNetworkingRoutingtableV2Update,
		DeleteContext: resourceNetworkingRoutingtableV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"distributed": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceNetworkingRoutingtableV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud networking client: %s", err)
	}

	createOpts := RoutingtableCreateOpts{
		routingtables.CreateOpts{
			Name:  d.Get("name").(string),
			VpcID: d.Get("vpc_id").(string),
		},
	}

	if dRaw, ok := d.GetOkExists("distributed"); ok {
		d := dRaw.(bool)
		createOpts.Distributed = &d
	}

	var finalCreateOpts routingtables.CreateOptsBuilder
	finalCreateOpts = createOpts

	log.Printf("[DEBUG] nhncloud_networking_routingtable_v2 create options: %#v", finalCreateOpts)
	n, err := routingtables.Create(networkingClient, finalCreateOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating nhncloud_networking_routingtable_v2: %s", err)
	}

	log.Printf("[DEBUG] Waiting for nhncloud_networking_routingtable_v2 %s to become available.", n.ID)

	d.SetId(n.ID)

	log.Printf("[DEBUG] Created nhncloud_networking_routingtable_v2 %s: %#v", n.ID, n)
	return resourceNetworkingRoutingtableV2Read(ctx, d, meta)
}

func resourceNetworkingRoutingtableV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud networking client: %s", err)
	}

	var routingtable routingtableExtended

	err = routingtables.Get(networkingClient, d.Id()).ExtractInto(&routingtable)
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error getting nhncloud_networking_routingtable_v2"))
	}

	log.Printf("[DEBUG] Retrieved nhncloud_networking_routingtable_v2 %s: %#v", d.Id(), routingtable)

	d.Set("name", routingtable.Name)
	d.Set("id", routingtable.ID)
	d.Set("distributed", routingtable.Distributed)
	d.Set("tenant_id", routingtable.TenantID)
	d.Set("default_table", routingtable.DefaultTable)
	d.Set("state", routingtable.State)
	d.Set("create_time", routingtable.CreateTime)

	return nil
}

func resourceNetworkingRoutingtableV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud networking client: %s", err)
	}

	// Declare finalUpdateOpts interface and basic updateOpts structure.
	var (
		finalUpdateOpts routingtables.UpdateOptsBuilder
		updateOpts      routingtables.UpdateOpts
	)

	// Populate basic updateOpts.
	if d.HasChange("name") {
		name := d.Get("name").(string)
		updateOpts.Name = &name
	}
	if d.HasChange("distributed") {
		distributed := d.Get("distributed").(bool)
		updateOpts.Distributed = &distributed
	}

	// Save basic updateOpts into finalUpdateOpts.
	finalUpdateOpts = updateOpts

	log.Printf("[DEBUG] nhncloud_networking_routingtable_v2 %s update options: %#v", d.Id(), finalUpdateOpts)
	_, err = routingtables.Update(networkingClient, d.Id(), finalUpdateOpts).Extract()
	if err != nil {
		return diag.Errorf("Error updating nhncloud_networking_routingtable_v2 %s: %s", d.Id(), err)
	}

	return resourceNetworkingRoutingtableV2Read(ctx, d, meta)
}

func resourceNetworkingRoutingtableV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud networking client: %s", err)
	}

	if err := routingtables.Delete(networkingClient, d.Id()).ExtractErr(); err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting nhncloud_networking_routingtable_v2"))
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"available"},
		Target:     []string{"DELETED"},
		Refresh:    resourceNetworkingRoutingtableV2StateRefreshFunc(networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for nhncloud_networking_routingtable_v2 %s to Delete:  %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}
