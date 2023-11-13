package nhncloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/nhn/nhncloud.gophercloud/openstack/networking/v2/vpcs"
)

func dataSourceNetworkingVPCV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkingVPCV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"external": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"shared": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"routingtables": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"default_table": {
							Type:     schema.TypeBool,
							Computed: true,
						},

						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"subnets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"external": {
							Type:     schema.TypeBool,
							Computed: true,
						},

						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"enable_dhcp": {
							Type:     schema.TypeBool,
							Computed: true,
						},

						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"gateway": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"routingtable_gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"routingtable_default_table": {
							Type:     schema.TypeBool,
							Computed: true,
						},

						"routingtable_explicit": {
							Type:     schema.TypeBool,
							Computed: true,
						},

						"routingtable_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"routingtable_name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"routes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"subnet_id": {
										Type:     schema.TypeString,
										Computed: true,
									},

									"tenant_id": {
										Type:     schema.TypeString,
										Computed: true,
									},

									"mask": {
										Type:     schema.TypeInt,
										Computed: true,
									},

									"gateway": {
										Type:     schema.TypeString,
										Computed: true,
									},

									"cidr": {
										Type:     schema.TypeString,
										Computed: true,
									},

									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},

						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"available_ip_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"vpc_shared": {
							Type:     schema.TypeBool,
							Computed: true,
						},

						"vpc_state": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"vpc_cidrv4": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"vpc_name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"shared": {
							Type:     schema.TypeBool,
							Computed: true,
						},

						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"cidrv4": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNetworkingVPCV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud networking client: %s", err)
	}

	// @tc-iaas-compute/1452
	// It uses list api, but finally fetches a single data source,
	// so I don't think query options other than identifiers are necessary.
	var listOpts vpcs.ListOptsBuilder
	listOpts = vpcs.ListOpts{
		ID:       d.Get("id").(string),
		Name:     d.Get("name").(string),
		TenantID: d.Get("tenant_id").(string),
	}

	pages, err := vpcs.List(networkingClient, listOpts).AllPages()
	if err != nil {
		return diag.FromErr(err)
	}

	tmpAllNetworks, err := vpcs.ExtractNetworks(pages)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(tmpAllNetworks) != 1 {
		return diag.Errorf("There is not a single result: %d", len(tmpAllNetworks))
	}

	// @tc-iaas-compute/1452
	// Next step which is for GET(detail), not LIST(summary)
	var network vpcs.NetworkDetail
	err = vpcs.Get(networkingClient, tmpAllNetworks[0].ID).ExtractInto(&network)
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error getting nhncloud_networking_vpc_v2"))
	}

	log.Printf("[DEBUG] Retrieved nhncloud_networking_vpc_v2 %s: %+v", network.ID, network)

	d.SetId(network.ID)
	d.Set("region", GetRegion(d, config))

	for key, val := range NewTransformer().Transform(&network) {
		d.Set(key, val)
	}

	return nil
}
