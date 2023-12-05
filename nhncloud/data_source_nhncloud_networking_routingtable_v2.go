package nhncloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nhn/nhncloud.gophercloud/nhncloud/networking/v2/routingtables"
	"log"
)

func dataSourceNetworkingRoutingtableV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkingRoutingtableV2Read,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"default_table": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"distributed": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"gateway_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"state": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"vpcs": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"subnets": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"routes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},

						"mask": {
							Type:     schema.TypeInt,
							Computed: true,
							Optional: true,
						},

						"gateway": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},

						"gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},

						"routingtable_id": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},

						"cidr": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},

						"id": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
					},
				},
			},

			"create_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceNetworkingRoutingtableV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud networking client: %s", err)
	}

	listOpts := routingtables.ListOpts{}

	if v, ok := d.GetOk("id"); ok {
		listOpts.ID = v.(string)
	}

	if v, ok := d.GetOk("name"); ok {
		listOpts.Name = v.(string)
	}

	if v, ok := d.GetOkExists("distributed"); ok {
		dist := v.(bool)
		listOpts.Distributed = &dist
	}

	pages, err := routingtables.List(networkingClient, listOpts).AllPages()
	if err != nil {
		return diag.FromErr(err)
	}

	tmpAllRoutingtables, err := routingtables.ExtractRoutingtables(pages)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(tmpAllRoutingtables) != 1 {
		return diag.Errorf("There is not a single result: %d", len(tmpAllRoutingtables))
	}

	var routingtableDetail routingtables.RoutingtableDetail
	err = routingtables.Get(networkingClient, tmpAllRoutingtables[0].ID).ExtractInto(&routingtableDetail)
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error getting nhncloud_networking_routingtable_v2"))
	}

	log.Printf("[DEBUG] Retrieved nhncloud_networking_routingtable_v2 %s: %+v", routingtableDetail.ID, routingtableDetail)

	d.SetId(routingtableDetail.ID)

	d.Set("name", routingtableDetail.Name)
	d.Set("default_table", routingtableDetail.DefaultTable)
	d.Set("distributed", routingtableDetail.Distributed)
	d.Set("state", routingtableDetail.State)
	d.Set("tenant_id", routingtableDetail.TenantID)
	d.Set("gateway_id", routingtableDetail.GatewayID)

	if err := d.Set("vpcs", routingtableDetail.VPCs); err != nil {
		log.Printf("[DEBUG] Unable to set vpcs: %s", err)
	}
	if err := d.Set("subnets", routingtableDetail.Subnets); err != nil {
		log.Printf("[DEBUG] Unable to set subnets: %s", err)
	}

	//routes := make([]routingtables.Route, 0, len(routingtableDetail.Routes))
	//for _, v := range routingtableDetail.Routes {
	//	routes = append(routes, routingtables.Route{
	//		ID: v.ID,
	//		CIDR:  v.CIDR,
	//		Gateway: v.Gateway,
	//		GatewayID: v.GatewayID,
	//		RoutingtableID: v.RoutingtableID,
	//		TenantID: v.TenantID,
	//		Mask: v.Mask,
	//	})
	//}
	if err = d.Set("routes", routingtableDetail.Routes); err != nil {
		log.Printf("[DEBUG] Unable to set routes: %s", err)
	}
	//d.Set("region", GetRegion(d, config))

	//for key, val := range NewTransformer().Transform(&routingtableDetail) {
	//	d.Set(key, val)
	//}

	return nil
}
