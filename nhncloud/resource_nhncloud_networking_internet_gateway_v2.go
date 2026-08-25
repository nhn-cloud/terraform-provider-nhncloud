package nhncloud

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/nhn-cloud/nhncloud.gophercloud/nhncloud/networking/v2/internetgateways"
)

func resourceNetworkingInternetGatewayV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkingInternetGatewayV2Create,
		ReadContext:   resourceNetworkingInternetGatewayV2Read,
		DeleteContext: resourceNetworkingInternetGatewayV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				// The region in which the resource is located.
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				// 인터넷 게이트웨이 이름
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"external_network_id": {
				// 인터넷 게이트웨이가 연결할 외부 네트워크 ID
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"id": {
				// 인터넷 게이트웨이 ID
				Type:     schema.TypeString,
				Computed: true,
			},
			"routingtable_id": {
				// 라우팅 테이블과 연결되어 있을 때 인터넷 게이트웨이를 연결한 라우팅 테이블 ID
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				// 인터넷 게이트웨이의 상태.
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				// 인터넷 게이트웨이 생성 시간(UTC)
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_id": {
				// 인터넷 게이트웨이가 속한 테넌트 ID
				Type:     schema.TypeString,
				Computed: true,
			},
			"migrate_status": {
				// 점검으로 인한 인터넷 게이트웨이 이동 시 처리 상태
				Type:     schema.TypeString,
				Computed: true,
			},
			"migrate_error": {
				// 점검으로 인한 인터넷 게이트웨이 이동 중 오류 발생 시 오류 메시지
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNetworkingInternetGatewayV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud networking client: %s", err)
	}

	createOpts := internetgateways.CreateOpts{
		Name:              d.Get("name").(string),
		ExternalNetworkID: d.Get("external_network_id").(string),
	}

	log.Printf("[DEBUG] nhncloud_networking_internet_gateway_v2 create options: %#v", createOpts)
	n, err := internetgateways.Create(networkingClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating nhncloud_networking_internet_gateway_v2: %s", err)
	}

	d.SetId(n.ID)
	log.Printf("[DEBUG] Created nhncloud_networking_internet_gateway_v2 %s: %#v", n.ID, n)
	return resourceNetworkingInternetGatewayV2Read(ctx, d, meta)
}

func resourceNetworkingInternetGatewayV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud networking client: %s", err)
	}

	n, err := internetgateways.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error getting nhncloud_networking_internet_gateway_v2"))
	}

	log.Printf("[DEBUG] Retrieved nhncloud_networking_internet_gateway_v2 %s: %#v", d.Id(), n)

	d.Set("region", GetRegion(d, config))
	d.Set("name", n.Name)
	d.Set("external_network_id", n.ExternalNetworkID)
	d.Set("routingtable_id", n.RoutingtableID)
	d.Set("state", n.State)
	d.Set("create_time", n.CreateTime)
	d.Set("tenant_id", n.TenantID)
	d.Set("migrate_status", n.MigrateStatus)
	d.Set("migrate_error", n.MigrateError)
	return nil
}

func resourceNetworkingInternetGatewayV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud networking client: %s", err)
	}

	if err := internetgateways.Delete(networkingClient, d.Id()).ExtractErr(); err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting nhncloud_networking_internet_gateway_v2"))
	}

	d.SetId("")
	return nil
}
