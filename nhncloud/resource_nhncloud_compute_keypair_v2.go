package nhncloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/nhn-cloud/nhncloud.gophercloud/nhncloud/compute/v2/keypairs"
)

func resourceComputeKeypairV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceComputeKeypairV2Create,
		ReadContext:   resourceComputeKeypairV2Read,
		DeleteContext: resourceComputeKeypairV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
				ForceNew: true,
			},

			"public_key": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"value_specs": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},

			// computed-only
			"private_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"user_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceComputeKeypairV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	computeClient, err := config.ComputeV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud compute client: %s", err)
	}

	userID := d.Get("user_id").(string)
	if userID != "" {
		computeClient.Microversion = computeKeyPairV2UserIDMicroversion
	}

	name := d.Get("name").(string)
	createOpts := ComputeKeyPairV2CreateOpts{
		keypairs.CreateOpts{
			Name:      name,
			PublicKey: d.Get("public_key").(string),
			UserID:    d.Get("user_id").(string),
		},
		MapValueSpecs(d),
	}

	log.Printf("[DEBUG] nhncloud_compute_keypair_v2 create options: %#v", createOpts)

	kp, err := keypairs.Create(computeClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Unable to create nhncloud_compute_keypair_v2 %s: %s", name, err)
	}

	d.SetId(kp.Name)
	d.Set("user_id", d.Get("user_id").(string))

	// Private Key is only available in the response to a create.
	d.Set("private_key", kp.PrivateKey)

	return resourceComputeKeypairV2Read(ctx, d, meta)
}

func resourceComputeKeypairV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	computeClient, err := config.ComputeV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud compute client: %s", err)
	}

	userID := d.Get("user_id").(string)
	if userID != "" {
		computeClient.Microversion = computeKeyPairV2UserIDMicroversion
	}

	log.Printf("[DEBUG] Microversion %s", computeClient.Microversion)

	kpopts := keypairs.GetOpts{
		UserID: userID,
	}

	kp, err := keypairs.Get(computeClient, d.Id(), kpopts).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving nhncloud_compute_keypair_v2"))
	}

	log.Printf("[DEBUG] Retrieved nhncloud_compute_keypair_v2 %s: %#v", d.Id(), kp)

	d.Set("name", kp.Name)
	d.Set("public_key", kp.PublicKey)
	d.Set("fingerprint", kp.Fingerprint)
	d.Set("region", GetRegion(d, config))
	if userID != "" {
		d.Set("user_id", kp.UserID)
	}

	return nil
}

func resourceComputeKeypairV2Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	computeClient, err := config.ComputeV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud compute client: %s", err)
	}

	userID := d.Get("user_id").(string)
	if userID != "" {
		computeClient.Microversion = computeKeyPairV2UserIDMicroversion
	}

	log.Printf("[DEBUG] User ID %s", userID)
	log.Printf("[DEBUG] Microversion %s", computeClient.Microversion)

	kpopts := keypairs.DeleteOpts{
		UserID: userID,
	}

	err = keypairs.Delete(computeClient, d.Id(), kpopts).ExtractErr()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting nhncloud_compute_keypair_v2"))
	}

	return nil
}
