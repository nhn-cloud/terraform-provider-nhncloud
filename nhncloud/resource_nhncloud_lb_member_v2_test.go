package nhncloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/nhn-cloud/nhncloud.gophercloud/nhncloud/networking/v2/extensions/lbaas_v2/pools"
)

func TestAccLBV2Member_basic(t *testing.T) {
	var member1 pools.Member
	var member2 pools.Member

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckLB(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckLBV2MemberDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccLbV2MemberConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2MemberExists("nhncloud_lb_member_v2.member_1", &member1),
					testAccCheckLBV2MemberExists("nhncloud_lb_member_v2.member_2", &member2),
					resource.TestCheckResourceAttr("nhncloud_lb_member_v2.member_1", "backup", "true"),
				),
			},
			{
				Config: TestAccLbV2MemberConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("nhncloud_lb_member_v2.member_1", "weight", "10"),
					resource.TestCheckResourceAttr("nhncloud_lb_member_v2.member_1", "backup", "false"),
					resource.TestCheckResourceAttr("nhncloud_lb_member_v2.member_2", "weight", "15"),
				),
			},
		},
	})
}

func TestAccLBV2Member_monitor(t *testing.T) {
	var member1 pools.Member
	var member2 pools.Member

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckLB(t)
			testAccPreCheckUseOctavia(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckLBV2MemberDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccLbV2MemberMonitor,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2MemberExists("nhncloud_lb_member_v2.member_1", &member1),
					testAccCheckLBV2MemberExists("nhncloud_lb_member_v2.member_2", &member2),
					resource.TestCheckResourceAttr("nhncloud_lb_member_v2.member_1", "monitor_address", "192.168.199.110"),
					resource.TestCheckResourceAttr("nhncloud_lb_member_v2.member_1", "monitor_port", "8080"),
					resource.TestCheckResourceAttr("nhncloud_lb_member_v2.member_2", "monitor_address", "192.168.199.111"),
					resource.TestCheckResourceAttr("nhncloud_lb_member_v2.member_2", "monitor_port", "8080"),
				),
			},
			{
				Config: TestAccLbV2MemberMonitorUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("nhncloud_lb_member_v2.member_1", "monitor_address", "192.168.199.110"),
					resource.TestCheckResourceAttr("nhncloud_lb_member_v2.member_1", "monitor_port", "8080"),
					resource.TestCheckResourceAttr("nhncloud_lb_member_v2.member_2", "monitor_address", "192.168.199.110"),
					resource.TestCheckResourceAttr("nhncloud_lb_member_v2.member_2", "monitor_port", "443"),
				),
			},
		},
	})
}

func testAccCheckLBV2MemberDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	lbClient, err := chooseLBV2AccTestClient(config, osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack load balancing client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nhncloud_lb_member_v2" {
			continue
		}

		poolID := rs.Primary.Attributes["pool_id"]
		_, err := pools.GetMember(lbClient, poolID, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Member still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckLBV2MemberExists(n string, member *pools.Member) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		lbClient, err := chooseLBV2AccTestClient(config, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack load balancing client: %s", err)
		}

		poolID := rs.Primary.Attributes["pool_id"]
		found, err := pools.GetMember(lbClient, poolID, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Member not found")
		}

		*member = *found

		return nil
	}
}

const TestAccLbV2MemberConfigBasic = `
resource "nhncloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "nhncloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  network_id = "${nhncloud_networking_network_v2.network_1.id}"
  cidr = "192.168.199.0/24"
  ip_version = 4
}

resource "nhncloud_lb_loadbalancer_v2" "loadbalancer_1" {
  name = "loadbalancer_1"
  vip_subnet_id = "${nhncloud_networking_subnet_v2.subnet_1.id}"
  vip_address = "192.168.199.10"

  timeouts {
    create = "15m"
    update = "15m"
    delete = "15m"
  }
}

resource "nhncloud_lb_listener_v2" "listener_1" {
  name = "listener_1"
  protocol = "HTTP"
  protocol_port = 8080
  loadbalancer_id = "${nhncloud_lb_loadbalancer_v2.loadbalancer_1.id}"
}

resource "nhncloud_lb_pool_v2" "pool_1" {
  name = "pool_1"
  protocol = "HTTP"
  lb_method = "ROUND_ROBIN"
  listener_id = "${nhncloud_lb_listener_v2.listener_1.id}"
}

resource "nhncloud_lb_member_v2" "member_1" {
  address = "192.168.199.110"
  protocol_port = 8080
  pool_id = "${nhncloud_lb_pool_v2.pool_1.id}"
  subnet_id = "${nhncloud_networking_subnet_v2.subnet_1.id}"
  weight = 0
  backup = true

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}

resource "nhncloud_lb_member_v2" "member_2" {
  address = "192.168.199.111"
  protocol_port = 8080
  pool_id = "${nhncloud_lb_pool_v2.pool_1.id}"
  subnet_id = "${nhncloud_networking_subnet_v2.subnet_1.id}"

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`

const TestAccLbV2MemberConfigUpdate = `
resource "nhncloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "nhncloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${nhncloud_networking_network_v2.network_1.id}"
}

resource "nhncloud_lb_loadbalancer_v2" "loadbalancer_1" {
  name = "loadbalancer_1"
  vip_subnet_id = "${nhncloud_networking_subnet_v2.subnet_1.id}"

  timeouts {
    create = "15m"
    update = "15m"
    delete = "15m"
  }
}

resource "nhncloud_lb_listener_v2" "listener_1" {
  name = "listener_1"
  protocol = "HTTP"
  protocol_port = 8080
  loadbalancer_id = "${nhncloud_lb_loadbalancer_v2.loadbalancer_1.id}"
}

resource "nhncloud_lb_pool_v2" "pool_1" {
  name = "pool_1"
  protocol = "HTTP"
  lb_method = "ROUND_ROBIN"
  listener_id = "${nhncloud_lb_listener_v2.listener_1.id}"
}

resource "nhncloud_lb_member_v2" "member_1" {
  address = "192.168.199.110"
  protocol_port = 8080
  weight = 10
  admin_state_up = "true"
  pool_id = "${nhncloud_lb_pool_v2.pool_1.id}"
  subnet_id = "${nhncloud_networking_subnet_v2.subnet_1.id}"
  backup = false

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}

resource "nhncloud_lb_member_v2" "member_2" {
  address = "192.168.199.111"
  protocol_port = 8080
  weight = 15
  admin_state_up = "true"
  pool_id = "${nhncloud_lb_pool_v2.pool_1.id}"
  subnet_id = "${nhncloud_networking_subnet_v2.subnet_1.id}"

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`

const TestAccLbV2MemberMonitor = `
resource "nhncloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "nhncloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  network_id = "${nhncloud_networking_network_v2.network_1.id}"
  cidr = "192.168.199.0/24"
  ip_version = 4
}

resource "nhncloud_lb_loadbalancer_v2" "loadbalancer_1" {
  name = "loadbalancer_1"
  vip_subnet_id = "${nhncloud_networking_subnet_v2.subnet_1.id}"
  vip_address = "192.168.199.10"

  timeouts {
    create = "15m"
    update = "15m"
    delete = "15m"
  }
}

resource "nhncloud_lb_listener_v2" "listener_1" {
  name = "listener_1"
  protocol = "HTTP"
  protocol_port = 8080
  loadbalancer_id = "${nhncloud_lb_loadbalancer_v2.loadbalancer_1.id}"
}

resource "nhncloud_lb_pool_v2" "pool_1" {
  name = "pool_1"
  protocol = "HTTP"
  lb_method = "ROUND_ROBIN"
  listener_id = "${nhncloud_lb_listener_v2.listener_1.id}"
}

resource "nhncloud_lb_member_v2" "member_1" {
  address = "192.168.199.110"
  protocol_port = 8080
  pool_id = "${nhncloud_lb_pool_v2.pool_1.id}"
  subnet_id = "${nhncloud_networking_subnet_v2.subnet_1.id}"
  weight = 0
  monitor_address = "192.168.199.110"
  monitor_port = 8080

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}

resource "nhncloud_lb_member_v2" "member_2" {
  address = "192.168.199.111"
  protocol_port = 8080
  pool_id = "${nhncloud_lb_pool_v2.pool_1.id}"
  subnet_id = "${nhncloud_networking_subnet_v2.subnet_1.id}"
  monitor_address = "192.168.199.111"
  monitor_port = 8080

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`

const TestAccLbV2MemberMonitorUpdate = `
resource "nhncloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "nhncloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${nhncloud_networking_network_v2.network_1.id}"
}

resource "nhncloud_lb_loadbalancer_v2" "loadbalancer_1" {
  name = "loadbalancer_1"
  vip_subnet_id = "${nhncloud_networking_subnet_v2.subnet_1.id}"

  timeouts {
    create = "15m"
    update = "15m"
    delete = "15m"
  }
}

resource "nhncloud_lb_listener_v2" "listener_1" {
  name = "listener_1"
  protocol = "HTTP"
  protocol_port = 8080
  loadbalancer_id = "${nhncloud_lb_loadbalancer_v2.loadbalancer_1.id}"
}

resource "nhncloud_lb_pool_v2" "pool_1" {
  name = "pool_1"
  protocol = "HTTP"
  lb_method = "ROUND_ROBIN"
  listener_id = "${nhncloud_lb_listener_v2.listener_1.id}"
}

resource "nhncloud_lb_member_v2" "member_1" {
  address = "192.168.199.110"
  protocol_port = 8080
  weight = 10
  admin_state_up = "true"
  pool_id = "${nhncloud_lb_pool_v2.pool_1.id}"
  subnet_id = "${nhncloud_networking_subnet_v2.subnet_1.id}"
  monitor_address = "192.168.199.110"
  monitor_port = 8080

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}

resource "nhncloud_lb_member_v2" "member_2" {
  address = "192.168.199.111"
  protocol_port = 8080
  weight = 15
  admin_state_up = "true"
  pool_id = "${nhncloud_lb_pool_v2.pool_1.id}"
  subnet_id = "${nhncloud_networking_subnet_v2.subnet_1.id}"
  monitor_address = "192.168.199.110"
  monitor_port = 443

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`
