package nhncloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/nhn-cloud/nhncloud.gophercloud/nhncloud/compute/v2/keypairs"
)

func dataSourceComputeKeypairV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceComputeKeypairV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			// computed-only
			"fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"public_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceComputeKeypairV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	computeClient, err := config.ComputeV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud compute client: %s", err)
	}

	name := d.Get("name").(string)
	kp, err := keypairs.Get(computeClient, name, keypairs.GetOpts{}).Extract()
	if err != nil {
		return diag.Errorf("Error retrieving nhncloud_compute_keypair_v2 %s: %s", name, err)
	}

	d.SetId(name)

	log.Printf("[DEBUG] Retrieved nhncloud_compute_keypair_v2 %s: %#v", d.Id(), kp)

	d.Set("fingerprint", kp.Fingerprint)
	d.Set("public_key", kp.PublicKey)
	d.Set("region", GetRegion(d, config))

	return nil
}
