package nhncloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/nhn-cloud/nhncloud.gophercloud/nhncloud/networking/v2/internetgateways"
)

func dataSourceNetworkingInternetGatewayV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkingInternetGatewayV2Read,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"region": {
				// The region in which the resource is located.
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				// 인터넷 게이트웨이 이름
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"external_network_id": {
				// 인터넷 게이트웨이가 연결할 외부 네트워크 ID
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"routingtable_id": {
				// 라우팅 테이블과 연결되어 있을 때 인터넷 게이트웨이를 연결한 라우팅 테이블 ID
				Type:     schema.TypeString,
				Optional: true,
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
				Optional: true,
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

func dataSourceNetworkingInternetGatewayV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating NHN Cloud networking client: %s", err)
	}

	listOpts := internetgateways.ListOpts{}

	if v, ok := d.GetOk("id"); ok {
		listOpts.ID = v.(string)
	}
	if v, ok := d.GetOk("name"); ok {
		listOpts.Name = v.(string)
	}
	if v, ok := d.GetOk("external_network_id"); ok {
		listOpts.ExternalNetworkID = v.(string)
	}
	if v, ok := d.GetOk("routingtable_id"); ok {
		listOpts.RoutingtableID = v.(string)
	}
	if v, ok := d.GetOk("tenant_id"); ok {
		listOpts.TenantID = v.(string)
	}

	pages, err := internetgateways.List(networkingClient, listOpts).AllPages()
	if err != nil {
		return diag.FromErr(err)
	}

	allInternetGateways, err := internetgateways.ExtractInternetGateways(pages)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(allInternetGateways) != 1 {
		return diag.Errorf("Your query returned %d internet gateways. Please change your search criteria and try again.", len(allInternetGateways))
	}

	// The list item already carries the full field set for an internet gateway
	// (identical to Get), so use it directly instead of a second round trip.
	n := allInternetGateways[0]

	d.SetId(n.ID)
	d.Set("region", GetRegion(d, config))
	d.Set("name", n.Name)
	d.Set("external_network_id", n.ExternalNetworkID)
	d.Set("routingtable_id", n.RoutingtableID)
	d.Set("state", n.State)
	d.Set("create_time", n.CreateTime)
	d.Set("tenant_id", n.TenantID)
	d.Set("migrate_status", n.MigrateStatus)
	d.Set("migrate_error", n.MigrateError)
	log.Printf("[DEBUG] Retrieved nhncloud_networking_internet_gateway_v2 %s: %#v", d.Id(), n)
	return nil
}
