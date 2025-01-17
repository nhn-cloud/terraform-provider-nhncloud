package nhncloud

import (
	"context"
	"github.com/gophercloud/gophercloud"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/nhn-cloud/nhncloud.gophercloud/nhncloud/networking/v2/vpcsubnets"
)

func dataSourceNetworkingVPCSubnetV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkingVPCSubnetV2Read,

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

			"cidr": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsCIDR,
				},
			},

			"shared": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"external": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"state": {
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
				Optional: true,
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

			"gateway": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNetworkingVPCSubnetV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud networking client: %s", err)
	}

	listOpts := vpcsubnets.ListOpts{}

	if v, ok := d.GetOk("tenant_id"); ok {
		listOpts.TenantID = v.(string)
	}

	if v, ok := d.GetOk("id"); ok {
		listOpts.ID = v.(string)
	}

	if v, ok := d.GetOk("name"); ok {
		listOpts.Name = v.(string)
	}

	if v, ok := d.GetOkExists("shared"); ok {
		shared := v.(bool)
		listOpts.Shared = &shared
	}

	pages, err := vpcsubnets.List(networkingClient, listOpts).AllPages()
	if err != nil {
		return diag.Errorf("Unable to retrieve nhncloud_networking_vpcsubnet_v2: %s", err)
	}

	// List API which is pre-flight
	allSubnets, err := vpcsubnets.ExtractSubnets(pages)
	if err != nil {
		return diag.Errorf("Unable to extract nhncloud_networking_vpcsubnet_v2: %s", err)
	}

	if len(allSubnets) < 1 {
		return diag.Errorf("Your query returned no nhncloud_networking_vpcsubnet_v2. " +
			"Please change your search criteria and try again.")
	}

	if len(allSubnets) > 1 {
		return diag.Errorf("Your query returned more than one nhncloud_networking_vpcsubnet_v2." +
			" Please try a more specific search criteria")
	}

	subnet := allSubnets[0]
	subnetDetail, err := vpcsubnets.Get(networkingClient, subnet.ID).Extract()
	if err != nil {
		if _, ok := err.(gophercloud.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error retrieving nhncloud_networking_vpcsubnet_v2: %s", err)
	}

	log.Printf("[DEBUG] Retrieved nhncloud_networking_vpcsubnet_v2 %s: %+v", subnet.ID, subnet)

	d.SetId(subnet.ID)
	d.Set("region", GetRegion(d, config))

	for key, val := range NewTransformer().Transform(subnetDetail) {
		d.Set(key, val)
	}

	return nil
}
