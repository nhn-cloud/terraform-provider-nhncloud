package nhncloud

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/groups"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	octavialoadbalancers "github.com/nhn-cloud/nhncloud.gophercloud/nhncloud/loadbalancer/v2/loadbalancers"
	"github.com/nhn-cloud/nhncloud.gophercloud/nhncloud/networking/v2/extensions/lbaas_v2/loadbalancers"
)

func TestAccLBV2LoadBalancer_basic(t *testing.T) {
	var lb loadbalancers.LoadBalancer

	lbProvider := "haproxy"
	if osUseOctavia != "" {
		lbProvider = "octavia"
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckLB(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckLBV2LoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLbV2LoadBalancerConfigBasic(lbProvider),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2LoadBalancerExists("nhncloud_lb_loadbalancer_v2.loadbalancer_1", &lb),
					testAccCheckLBV2LoadBalancerHasTag("nhncloud_lb_loadbalancer_v2.loadbalancer_1", "tag1"),
					testAccCheckLBV2LoadBalancerTagCount("nhncloud_lb_loadbalancer_v2.loadbalancer_1", 1),
				),
			},
			{
				Config: testAccLbV2LoadBalancerConfigUpdate(lbProvider),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2LoadBalancerHasTag("nhncloud_lb_loadbalancer_v2.loadbalancer_1", "tag1"),
					testAccCheckLBV2LoadBalancerHasTag("nhncloud_lb_loadbalancer_v2.loadbalancer_1", "tag2"),
					testAccCheckLBV2LoadBalancerTagCount("nhncloud_lb_loadbalancer_v2.loadbalancer_1", 2),
					resource.TestCheckResourceAttr(
						"nhncloud_lb_loadbalancer_v2.loadbalancer_1", "name", "loadbalancer_1_updated"),
					resource.TestMatchResourceAttr(
						"nhncloud_lb_loadbalancer_v2.loadbalancer_1", "vip_port_id",
						regexp.MustCompile("^[a-f0-9-]+")),
				),
			},
		},
	})
}

func TestAccLBV2LoadBalancer_secGroup(t *testing.T) {
	var lb loadbalancers.LoadBalancer
	var sg1, sg2 groups.SecGroup

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckLB(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckLBV2LoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLbV2LoadBalancerSecGroup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2LoadBalancerExists(
						"nhncloud_lb_loadbalancer_v2.loadbalancer_1", &lb),
					testAccCheckNetworkingV2SecGroupExists(
						"nhncloud_networking_secgroup_v2.secgroup_1", &sg1),
					testAccCheckNetworkingV2SecGroupExists(
						"nhncloud_networking_secgroup_v2.secgroup_1", &sg2),
					resource.TestCheckResourceAttr(
						"nhncloud_lb_loadbalancer_v2.loadbalancer_1", "security_group_ids.#", "1"),
					testAccCheckLBV2LoadBalancerHasSecGroup(&lb, &sg1),
				),
			},
			{
				Config: testAccLbV2LoadBalancerSecGroupUpdate1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2LoadBalancerExists(
						"nhncloud_lb_loadbalancer_v2.loadbalancer_1", &lb),
					testAccCheckNetworkingV2SecGroupExists(
						"nhncloud_networking_secgroup_v2.secgroup_2", &sg1),
					testAccCheckNetworkingV2SecGroupExists(
						"nhncloud_networking_secgroup_v2.secgroup_2", &sg2),
					resource.TestCheckResourceAttr(
						"nhncloud_lb_loadbalancer_v2.loadbalancer_1", "security_group_ids.#", "2"),
					testAccCheckLBV2LoadBalancerHasSecGroup(&lb, &sg1),
					testAccCheckLBV2LoadBalancerHasSecGroup(&lb, &sg2),
				),
			},
			{
				Config: testAccLbV2LoadBalancerSecGroupUpdate2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2LoadBalancerExists(
						"nhncloud_lb_loadbalancer_v2.loadbalancer_1", &lb),
					testAccCheckNetworkingV2SecGroupExists(
						"nhncloud_networking_secgroup_v2.secgroup_2", &sg1),
					testAccCheckNetworkingV2SecGroupExists(
						"nhncloud_networking_secgroup_v2.secgroup_2", &sg2),
					resource.TestCheckResourceAttr(
						"nhncloud_lb_loadbalancer_v2.loadbalancer_1", "security_group_ids.#", "1"),
					testAccCheckLBV2LoadBalancerHasSecGroup(&lb, &sg2),
				),
			},
		},
	})
}

func TestAccLBV2LoadBalancer_vip_network(t *testing.T) {
	var lb loadbalancers.LoadBalancer

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckLB(t)
			testAccPreCheckUseOctavia(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckLBV2LoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLbV2LoadBalancerConfigVIPNetwork,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2LoadBalancerExists("nhncloud_lb_loadbalancer_v2.loadbalancer_1", &lb),
				),
			},
		},
	})
}

func TestAccLBV2LoadBalancer_vip_port_id(t *testing.T) {
	var lb loadbalancers.LoadBalancer
	var port ports.Port

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckLB(t)
			testAccPreCheckUseOctavia(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckLBV2LoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLbV2LoadBalancerConfigVIPPortID,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2LoadBalancerExists(
						"nhncloud_lb_loadbalancer_v2.loadbalancer_1", &lb),
					testAccCheckNetworkingV2PortExists(
						"nhncloud_networking_port_v2.port_1", &port),
					resource.TestCheckResourceAttrPtr(
						"nhncloud_lb_loadbalancer_v2.loadbalancer_1", "vip_port_id", &port.ID),
				),
			},
		},
	})
}

func testAccCheckLBV2LoadBalancerDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	lbClient, err := chooseLBV2AccTestClient(config, osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack load balancing client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nhncloud_lb_loadbalancer_v2" {
			continue
		}

		lb, err := loadbalancers.Get(lbClient, rs.Primary.ID).Extract()
		if err == nil && lb.ProvisioningStatus != "DELETED" {
			return fmt.Errorf("LoadBalancer still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckLBV2LoadBalancerExists(
	n string, lb *loadbalancers.LoadBalancer) resource.TestCheckFunc {
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

		found, err := loadbalancers.Get(lbClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Loadbalancer not found")
		}

		*lb = *found

		return nil
	}
}
func testAccCheckLBV2LoadBalancerHasTag(n, tag string) resource.TestCheckFunc {
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

		found, err := octavialoadbalancers.Get(lbClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Loadbalancer not found")
		}

		for _, v := range found.Tags {
			if tag == v {
				return nil
			}
		}

		return fmt.Errorf("Tag not found: %s", tag)
	}
}

func testAccCheckLBV2LoadBalancerTagCount(n string, expected int) resource.TestCheckFunc {
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

		found, err := octavialoadbalancers.Get(lbClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Loadbalancer not found")
		}

		if len(found.Tags) != expected {
			return fmt.Errorf("Expecting %d tags, found %d", expected, len(found.Tags))
		}

		return nil
	}
}

func testAccCheckLBV2LoadBalancerHasSecGroup(
	lb *loadbalancers.LoadBalancer, sg *groups.SecGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		networkingClient, err := config.NetworkingV2Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack networking client: %s", err)
		}

		port, err := ports.Get(networkingClient, lb.VipPortID).Extract()
		if err != nil {
			return err
		}

		for _, p := range port.SecurityGroups {
			if p == sg.ID {
				return nil
			}
		}

		return fmt.Errorf("LoadBalancer does not have the security group")
	}
}

func testAccLbV2LoadBalancerConfigBasic(lbProvider string) string {
	return fmt.Sprintf(`
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
      loadbalancer_provider = "%s"
      vip_subnet_id = "${nhncloud_networking_subnet_v2.subnet_1.id}"
	  tags = ["tag1"]

      timeouts {
        create = "15m"
        update = "15m"
        delete = "15m"
      }
    }`, lbProvider)
}

func testAccLbV2LoadBalancerConfigUpdate(lbProvider string) string {
	return fmt.Sprintf(`
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
      name = "loadbalancer_1_updated"
      loadbalancer_provider = "%s"
      admin_state_up = "true"
      vip_subnet_id = "${nhncloud_networking_subnet_v2.subnet_1.id}"
	  tags = ["tag1", "tag2"]

      timeouts {
        create = "15m"
        update = "15m"
        delete = "15m"
      }
    }`, lbProvider)
}

const testAccLbV2LoadBalancerSecGroup = `
resource "nhncloud_networking_secgroup_v2" "secgroup_1" {
  name = "secgroup_1"
  description = "secgroup_1"
}

resource "nhncloud_networking_secgroup_v2" "secgroup_2" {
  name = "secgroup_2"
  description = "secgroup_2"
}

resource "nhncloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "nhncloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  network_id = "${nhncloud_networking_network_v2.network_1.id}"
  cidr = "192.168.199.0/24"
}

resource "nhncloud_lb_loadbalancer_v2" "loadbalancer_1" {
    name = "loadbalancer_1"
    vip_subnet_id = "${nhncloud_networking_subnet_v2.subnet_1.id}"
    security_group_ids = [
      "${nhncloud_networking_secgroup_v2.secgroup_1.id}"
    ]

    timeouts {
      create = "15m"
      update = "15m"
      delete = "15m"
    }
}
`

const testAccLbV2LoadBalancerSecGroupUpdate1 = `
resource "nhncloud_networking_secgroup_v2" "secgroup_1" {
  name = "secgroup_1"
  description = "secgroup_1"
}

resource "nhncloud_networking_secgroup_v2" "secgroup_2" {
  name = "secgroup_2"
  description = "secgroup_2"
}

resource "nhncloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "nhncloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  network_id = "${nhncloud_networking_network_v2.network_1.id}"
  cidr = "192.168.199.0/24"
}

resource "nhncloud_lb_loadbalancer_v2" "loadbalancer_1" {
  name = "loadbalancer_1"
  vip_subnet_id = "${nhncloud_networking_subnet_v2.subnet_1.id}"
  security_group_ids = [
    "${nhncloud_networking_secgroup_v2.secgroup_1.id}",
    "${nhncloud_networking_secgroup_v2.secgroup_2.id}"
  ]

  timeouts {
    create = "15m"
    update = "15m"
    delete = "15m"
  }
}
`

const testAccLbV2LoadBalancerSecGroupUpdate2 = `
resource "nhncloud_networking_secgroup_v2" "secgroup_1" {
  name = "secgroup_1"
  description = "secgroup_1"
}

resource "nhncloud_networking_secgroup_v2" "secgroup_2" {
  name = "secgroup_2"
  description = "secgroup_2"
}

resource "nhncloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "nhncloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  network_id = "${nhncloud_networking_network_v2.network_1.id}"
  cidr = "192.168.199.0/24"
}

resource "nhncloud_lb_loadbalancer_v2" "loadbalancer_1" {
  name = "loadbalancer_1"
  vip_subnet_id = "${nhncloud_networking_subnet_v2.subnet_1.id}"
  security_group_ids = [
    "${nhncloud_networking_secgroup_v2.secgroup_2.id}"
  ]
  depends_on = ["nhncloud_networking_secgroup_v2.secgroup_1"]

  timeouts {
    create = "15m"
    update = "15m"
    delete = "15m"
  }
}
`

const testAccLbV2LoadBalancerConfigVIPNetwork = `
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
  loadbalancer_provider = "octavia"
  vip_network_id = "${nhncloud_networking_network_v2.network_1.id}"
  depends_on = ["nhncloud_networking_subnet_v2.subnet_1"]
  timeouts {
    create = "15m"
    update = "15m"
    delete = "15m"
  }
}
`

const testAccLbV2LoadBalancerConfigVIPPortID = `
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

resource "nhncloud_networking_port_v2" "port_1" {
  name           = "port_1"
  network_id     = "${nhncloud_networking_network_v2.network_1.id}"
  admin_state_up = "true"
  depends_on = ["nhncloud_networking_subnet_v2.subnet_1"]
}

resource "nhncloud_lb_loadbalancer_v2" "loadbalancer_1" {
  name = "loadbalancer_1"
  loadbalancer_provider = "octavia"
  vip_port_id = "${nhncloud_networking_port_v2.port_1.id}"
  depends_on = ["nhncloud_networking_port_v2.port_1"]
  timeouts {
    create = "15m"
    update = "15m"
    delete = "15m"
  }
}
`
