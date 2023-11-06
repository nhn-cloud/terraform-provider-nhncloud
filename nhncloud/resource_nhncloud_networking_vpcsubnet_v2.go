package nhncloud

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	gophercloud "github.com/nhn/nhncloud.gophercloud"
	"github.com/nhn/nhncloud.gophercloud/openstack/networking/v2/vpcsubnets"
)

func resourceNetworkingVPCSubnetV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkingVPCSubnetV2Create,
		ReadContext:   resourceNetworkingVPCSubnetV2Read,
		UpdateContext: resourceNetworkingVPCSubnetV2Update,
		DeleteContext: resourceNetworkingVPCSubnetV2Delete,
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

			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"cidr": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsCIDR,
				},
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},

			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"id": {
				Type:     schema.TypeString,
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
				Optional: true,
				Computed: true,
				ForceNew: false,
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

			"shared": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"gateway": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNetworkingVPCSubnetV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	createOpts := vpcsubnets.CreateOpts{
		VpcId:    d.Get("vpc_id").(string),
		Name:     d.Get("name").(string),
		TenantID: d.Get("tenant_id").(string),
	}

	// Set CIDR if provided. Check if inferred subnet would match the provided cidr.
	if v, ok := d.GetOk("cidr"); ok {
		cidr := v.(string)
		_, netAddr, err := net.ParseCIDR(cidr)
		if err != nil {
			return diag.Errorf("Invalid CIDR %s: %s", cidr, err)
		}
		if netAddr.String() != cidr {
			return diag.Errorf("cidr %s doesn't match subnet address %s for openstack_networking_vpcsubnet_v2", cidr, netAddr.String())
		}
		createOpts.CIDR = cidr
	}

	log.Printf("[DEBUG] openstack_networking_vpcsubnet_v2 create options: %#v", createOpts)
	s, err := vpcsubnets.Create(networkingClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating openstack_networking_vpcsubnet_v2: %s", err)
	}

	log.Printf("[DEBUG] Waiting for openstack_networking_vpcsubnet_v2 %s to become available", s.ID)
	// Backend returns "available", but callback represents it to "ACTIVE"
	stateConf := &resource.StateChangeConf{
		Target:     []string{"ACTIVE"},
		Refresh:    networkingVPCSubnetV2StateRefreshFunc(networkingClient, s.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for openstack_networking_vpcsubnet_v2 %s to become available: %s", s.ID, err)
	}

	// Set routingtable_id if provided.
	if v, ok := d.GetOk("routingtable_id"); ok {
		if err = attachRoutingtable(networkingClient, v.(string), s.ID); err != nil {
			diagnostics := resourceNetworkingVPCSubnetV2Delete(ctx, d, meta)
			if diagnostics != nil {
				return diag.Errorf("Tried deleting the vpcsubnet created, due to the routingtable attaching error, but even it failed: %s: %s", diagnostics[0].Summary, err)
			}
			return diag.Errorf("Error creating openstack_networking_vpcsubnet_v2, as Error attaching routingtable: %s", err)
		}
	}

	d.SetId(s.ID)

	log.Printf("[DEBUG] Created openstack_networking_vpcsubnet_v2 %s: %#v", s.ID, s)
	return resourceNetworkingVPCSubnetV2Read(ctx, d, meta)
}

func resourceNetworkingVPCSubnetV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	subnetDetail, err := vpcsubnets.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error getting openstack_networking_vpcsubnet_v2"))
	}

	log.Printf("[DEBUG] Retrieved openstack_networking_vpcsubnet_v2 %s: %#v", d.Id(), subnetDetail)

	d.Set("region", GetRegion(d, config))

	for key, val := range NewTransformer().Transform(subnetDetail) {
		d.Set(key, val)
	}

	return nil
}

func resourceNetworkingVPCSubnetV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	var (
		updateOpts      vpcsubnets.UpdateOpts
		finalUpdateOpts vpcsubnets.UpdateOptsBuilder
	)

	if d.HasChange("name") {
		name := d.Get("name").(string)
		updateOpts.Name = &name

		finalUpdateOpts = updateOpts
		_, err = vpcsubnets.Update(networkingClient, d.Id(), finalUpdateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating openstack_networking_vpcsubnet_v2 %s: %s", d.Id(), err)
		}
	}

	if d.HasChange("routingtable_id") {
		if v, ok := d.GetOk("routingtable_id"); ok {
			detachRoutingtable(networkingClient, d.Id())
			if err = attachRoutingtable(networkingClient, v.(string), d.Id()); err != nil {
				return diag.Errorf("Error updating openstack_networking_vpcsubnet_v2, as Error attaching routingtable: %s", err)
			}
		} else {
			if err = detachRoutingtable(networkingClient, d.Id()); err != nil {
				return diag.Errorf("Error updating openstack_networking_vpcsubnet_v2, as Error detaching routingtable: %s", err)
			}
		}
	}

	return resourceNetworkingVPCSubnetV2Read(ctx, d, meta)
}

func resourceNetworkingVPCSubnetV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	if err := vpcsubnets.Delete(networkingClient, d.Id()).ExtractErr(); err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting openstack_networking_vpcsubnet_v2"))
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    networkingVPCSubnetV2StateRefreshFuncDelete(networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for openstack_networking_vpcsubnet_v2 %s to become deleted: %s", d.Id(), err)
	}

	return nil
}

func attachRoutingtable(networkingClient *gophercloud.ServiceClient, routingtableId string, subnetId string) error {
	var (
		attachOpts      vpcsubnets.AttachOpts
		finalAttachOpts vpcsubnets.AttachOptsBuilder
	)

	attachOpts.RoutingtableID = routingtableId
	finalAttachOpts = attachOpts

	return vpcsubnets.Attach(networkingClient, subnetId, finalAttachOpts).ExtractErr()
}

func detachRoutingtable(networkingClient *gophercloud.ServiceClient, subnetId string) error {
	return vpcsubnets.Detach(networkingClient, subnetId).ExtractErr()
}
