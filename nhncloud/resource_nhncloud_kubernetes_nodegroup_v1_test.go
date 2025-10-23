package nhncloud

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/nhn-cloud/nhncloud.gophercloud/nhncloud/kubernetes/v1/nodegroups"
)

func TestAccKubernetesV1NodeGroup_basic(t *testing.T) {
	var nodeGroup nodegroups.Nodegroup

	resourceName := "nhncloud_kubernetes_nodegroup_v1.nodegroup_1"
	clusterName := acctest.RandomWithPrefix("tf-acc-cluster")
	keypairName := acctest.RandomWithPrefix("tf-acc-keypair")
	clusterTemplateName := acctest.RandomWithPrefix("tf-acc-clustertemplate")
	nodeGroupName := acctest.RandomWithPrefix("tf-acc-nodegroup")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckKubernetes(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckKubernetesV1NodeGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesV1NodeGroupBasic(keypairName, clusterTemplateName, clusterName, nodeGroupName, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesV1NodeGroupExists(resourceName, &nodeGroup),
					resource.TestCheckResourceAttr(resourceName, "region", osRegionName),
					resource.TestCheckResourceAttrSet(resourceName, "cluster_id"),
					resource.TestCheckResourceAttr(resourceName, "name", nodeGroupName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttr(resourceName, "docker_volume_size", strconv.Itoa(10)),
					resource.TestCheckResourceAttrSet(resourceName, "labels.%"),
					resource.TestCheckResourceAttr(resourceName, "labels.kubescheduler_options", "log-flush-frequency=1m"),
					resource.TestCheckResourceAttr(resourceName, "role", "myRole"),
					resource.TestCheckResourceAttr(resourceName, "node_count", strconv.Itoa(1)),
					resource.TestCheckResourceAttr(resourceName, "min_node_count", strconv.Itoa(1)),
					resource.TestCheckResourceAttr(resourceName, "max_node_count", strconv.Itoa(5)),
					resource.TestCheckResourceAttr(resourceName, "image_id", osMagnumImage),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", osMagnumFlavor),
				),
			},
			{
				Config: testAccKubernetesV1NodeGroupBasic(keypairName, clusterTemplateName, clusterName, nodeGroupName, 2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesV1NodeGroupExists(resourceName, &nodeGroup),
					resource.TestCheckResourceAttr(resourceName, "region", osRegionName),
					resource.TestCheckResourceAttrSet(resourceName, "cluster_id"),
					resource.TestCheckResourceAttr(resourceName, "name", nodeGroupName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttr(resourceName, "docker_volume_size", strconv.Itoa(10)),
					resource.TestCheckResourceAttrSet(resourceName, "labels.%"),
					resource.TestCheckResourceAttr(resourceName, "labels.kubescheduler_options", "log-flush-frequency=1m"),
					resource.TestCheckResourceAttr(resourceName, "role", "myRole"),
					resource.TestCheckResourceAttr(resourceName, "node_count", strconv.Itoa(2)),
					resource.TestCheckResourceAttr(resourceName, "min_node_count", strconv.Itoa(1)),
					resource.TestCheckResourceAttr(resourceName, "max_node_count", strconv.Itoa(5)),
					resource.TestCheckResourceAttr(resourceName, "image_id", osMagnumImage),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", osMagnumFlavor),
				),
			},
		},
	})
}

func TestAccKubernetesV1NodeGroup_mergeLabels(t *testing.T) {
	var nodeGroup nodegroups.Nodegroup

	resourceName := "nhncloud_kubernetes_nodegroup_v1.nodegroup_1"
	clusterName := acctest.RandomWithPrefix("tf-acc-cluster")
	keypairName := acctest.RandomWithPrefix("tf-acc-keypair")
	clusterTemplateName := acctest.RandomWithPrefix("tf-acc-clustertemplate")
	nodeGroupName := acctest.RandomWithPrefix("tf-acc-nodegroup")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckKubernetes(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckKubernetesV1NodeGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesV1NodeGroupMergeLabels(keypairName, clusterTemplateName, clusterName, nodeGroupName, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesV1NodeGroupExists(resourceName, &nodeGroup),
					resource.TestCheckResourceAttr(resourceName, "region", osRegionName),
					resource.TestCheckResourceAttrSet(resourceName, "cluster_id"),
					resource.TestCheckResourceAttr(resourceName, "name", nodeGroupName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttr(resourceName, "docker_volume_size", strconv.Itoa(10)),
					resource.TestCheckResourceAttrSet(resourceName, "labels.%"),
					resource.TestCheckResourceAttr(resourceName, "labels.boot_volume_size", "15"),
					resource.TestCheckResourceAttr(resourceName, "merge_labels", "true"),
					resource.TestCheckResourceAttr(resourceName, "role", "myRole"),
					resource.TestCheckResourceAttr(resourceName, "node_count", strconv.Itoa(1)),
					resource.TestCheckResourceAttr(resourceName, "min_node_count", strconv.Itoa(1)),
					resource.TestCheckResourceAttr(resourceName, "max_node_count", strconv.Itoa(5)),
					resource.TestCheckResourceAttr(resourceName, "image_id", osMagnumImage),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", osMagnumFlavor),
				),
			},
			{
				Config: testAccKubernetesV1NodeGroupMergeLabels(keypairName, clusterTemplateName, clusterName, nodeGroupName, 2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesV1NodeGroupExists(resourceName, &nodeGroup),
					resource.TestCheckResourceAttr(resourceName, "region", osRegionName),
					resource.TestCheckResourceAttrSet(resourceName, "cluster_id"),
					resource.TestCheckResourceAttr(resourceName, "name", nodeGroupName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttr(resourceName, "docker_volume_size", strconv.Itoa(10)),
					resource.TestCheckResourceAttrSet(resourceName, "labels.%"),
					resource.TestCheckResourceAttr(resourceName, "labels.boot_volume_size", "15"),
					resource.TestCheckResourceAttr(resourceName, "merge_labels", "true"),
					resource.TestCheckResourceAttr(resourceName, "role", "myRole"),
					resource.TestCheckResourceAttr(resourceName, "node_count", strconv.Itoa(2)),
					resource.TestCheckResourceAttr(resourceName, "min_node_count", strconv.Itoa(1)),
					resource.TestCheckResourceAttr(resourceName, "max_node_count", strconv.Itoa(5)),
					resource.TestCheckResourceAttr(resourceName, "image_id", osMagnumImage),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", osMagnumFlavor),
				),
			},
		},
	})
}

func TestAccKubernetesV1NodeGroup_overrideLabels(t *testing.T) {
	var nodeGroup nodegroups.Nodegroup

	resourceName := "nhncloud_kubernetes_nodegroup_v1.nodegroup_1"
	clusterName := acctest.RandomWithPrefix("tf-acc-cluster")
	keypairName := acctest.RandomWithPrefix("tf-acc-keypair")
	clusterTemplateName := acctest.RandomWithPrefix("tf-acc-clustertemplate")
	nodeGroupName := acctest.RandomWithPrefix("tf-acc-nodegroup")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckKubernetes(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckKubernetesV1NodeGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesV1NodeGroupOverrideLabels(keypairName, clusterTemplateName, clusterName, nodeGroupName, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesV1NodeGroupExists(resourceName, &nodeGroup),
					resource.TestCheckResourceAttr(resourceName, "region", osRegionName),
					resource.TestCheckResourceAttrSet(resourceName, "cluster_id"),
					resource.TestCheckResourceAttr(resourceName, "name", nodeGroupName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttr(resourceName, "docker_volume_size", strconv.Itoa(10)),
					resource.TestCheckResourceAttrSet(resourceName, "labels.%"),
					resource.TestCheckResourceAttr(resourceName, "labels.kubescheduler_options", "log-flush-frequency=2m"),
					resource.TestCheckResourceAttr(resourceName, "merge_labels", "false"),
					resource.TestCheckResourceAttr(resourceName, "role", "myRole"),
					resource.TestCheckResourceAttr(resourceName, "node_count", strconv.Itoa(1)),
					resource.TestCheckResourceAttr(resourceName, "min_node_count", strconv.Itoa(1)),
					resource.TestCheckResourceAttr(resourceName, "max_node_count", strconv.Itoa(5)),
					resource.TestCheckResourceAttr(resourceName, "image_id", osMagnumImage),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", osMagnumFlavor),
				),
			},
			{
				Config: testAccKubernetesV1NodeGroupOverrideLabels(keypairName, clusterTemplateName, clusterName, nodeGroupName, 2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesV1NodeGroupExists(resourceName, &nodeGroup),
					resource.TestCheckResourceAttr(resourceName, "region", osRegionName),
					resource.TestCheckResourceAttrSet(resourceName, "cluster_id"),
					resource.TestCheckResourceAttr(resourceName, "name", nodeGroupName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttr(resourceName, "docker_volume_size", strconv.Itoa(10)),
					resource.TestCheckResourceAttrSet(resourceName, "labels.%"),
					resource.TestCheckResourceAttr(resourceName, "labels.kubescheduler_options", "log-flush-frequency=2m"),
					resource.TestCheckResourceAttr(resourceName, "merge_labels", "false"),
					resource.TestCheckResourceAttr(resourceName, "role", "myRole"),
					resource.TestCheckResourceAttr(resourceName, "node_count", strconv.Itoa(2)),
					resource.TestCheckResourceAttr(resourceName, "min_node_count", strconv.Itoa(1)),
					resource.TestCheckResourceAttr(resourceName, "max_node_count", strconv.Itoa(5)),
					resource.TestCheckResourceAttr(resourceName, "image_id", osMagnumImage),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", osMagnumFlavor),
				),
			},
		},
	})
}

func testAccCheckKubernetesV1NodeGroupExists(n string, nodeGroup *nodegroups.Nodegroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		kubernetesClient, err := config.ContainerInfraV1Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating NHN Cloud container infra client: %s", err)
		}

		kubernetesClient.Microversion = kubernetesV1NodeGroupMinMicroversion
		clusterID, nodeGroupID, err := parseNodeGroupID(rs.Primary.ID)
		if err != nil {
			return err
		}
		found, err := nodegroups.Get(kubernetesClient, clusterID, nodeGroupID).Extract()
		if err != nil {
			return err
		}

		if found.UUID != nodeGroupID {
			return fmt.Errorf("Nodegroup not found")
		}

		*nodeGroup = *found

		return nil
	}
}

func testAccCheckKubernetesV1NodeGroupDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	kubernetesClient, err := config.ContainerInfraV1Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating NHN Cloud container infra client: %s", err)
	}

	kubernetesClient.Microversion = kubernetesV1NodeGroupMinMicroversion

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nhncloud_kubernetes_nodegroup_v1" {
			continue
		}
		clusterID, nodeGroupID, err := parseVolumeTypeAccessID(rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = nodegroups.Get(kubernetesClient, clusterID, nodeGroupID).Extract()
		if err == nil {
			return fmt.Errorf("node group still exists")
		}
	}

	return nil
}

func testAccKubernetesV1NodeGroupBasic(keypairName, clusterTemplateName, clusterName string, nodeGroupName string, nodeCount int) string {
	return fmt.Sprintf(`
resource "nhncloud_compute_keypair_v2" "keypair_1" {
  name = "%s"
}

resource "nhncloud_kubernetes_clustertemplate_v1" "clustertemplate_1" {
  name                  = "%s"
  image                 = "%s"
  coe                   = "kubernetes"
  master_flavor         = "%s"
  flavor                = "%s"
  floating_ip_enabled   = true
  volume_driver         = "cinder"
  docker_storage_driver = "overlay2"
  docker_volume_size    = 5
  external_network_id   = "%s"
  network_driver        = "flannel"
  http_proxy            = "%s"
  https_proxy           = "%s"
  no_proxy              = "%s"
  labels = {
    kubescheduler_options = "log-flush-frequency=1m",
	%s
  }
}

resource "nhncloud_kubernetes_cluster_v1" "cluster_1" {
  name                 = "%s"
  cluster_template_id  = "${nhncloud_kubernetes_clustertemplate_v1.clustertemplate_1.id}"
  master_count         = 1
  node_count           = 1
  keypair              = "${nhncloud_compute_keypair_v2.keypair_1.name}"
}

resource "nhncloud_kubernetes_nodegroup_v1" "nodegroup_1" {
  name                 = "%s"
  cluster_id           = "${nhncloud_kubernetes_cluster_v1.cluster_1.id}"
  node_count           = %d
  docker_volume_size   = 10
  role				   = "myRole"
  min_node_count       = 1
  max_node_count       = 5
  image_id             = "%s"
  flavor_id            = "%s"
}
`, keypairName, clusterTemplateName, osMagnumImage, osMagnumFlavor, osMagnumFlavor, osExtGwID, osMagnumHTTPProxy, osMagnumHTTPSProxy, osMagnumNoProxy, osMagnumLabels, clusterName, nodeGroupName, nodeCount, osMagnumImage, osMagnumFlavor)
}

func testAccKubernetesV1NodeGroupMergeLabels(keypairName, clusterTemplateName, clusterName string, nodeGroupName string, nodeCount int) string {
	return fmt.Sprintf(`
resource "nhncloud_compute_keypair_v2" "keypair_1" {
  name = "%s"
}

resource "nhncloud_kubernetes_clustertemplate_v1" "clustertemplate_1" {
  name                  = "%s"
  image                 = "%s"
  coe                   = "kubernetes"
  master_flavor         = "%s"
  flavor                = "%s"
  floating_ip_enabled   = true
  volume_driver         = "cinder"
  docker_storage_driver = "overlay2"
  docker_volume_size    = 5
  external_network_id   = "%s"
  network_driver        = "flannel"
  http_proxy            = "%s"
  https_proxy           = "%s"
  no_proxy              = "%s"
  labels = {
    kubescheduler_options = "log-flush-frequency=1m",
	%s
  }
}

resource "nhncloud_kubernetes_cluster_v1" "cluster_1" {
  name                 = "%s"
  cluster_template_id  = "${nhncloud_kubernetes_clustertemplate_v1.clustertemplate_1.id}"
  master_count         = 1
  node_count           = 1
  keypair              = "${nhncloud_compute_keypair_v2.keypair_1.name}"
}

resource "nhncloud_kubernetes_nodegroup_v1" "nodegroup_1" {
  name                 = "%s"
  cluster_id           = "${nhncloud_kubernetes_cluster_v1.cluster_1.id}"
  node_count           = %d
  docker_volume_size   = 10
  role				   = "myRole"
  min_node_count       = 1
  max_node_count       = 5
  image_id             = "%s"
  flavor_id            = "%s"
  merge_labels         = true
  labels = {
    boot_volume_size = 15
  }
}
`, keypairName, clusterTemplateName, osMagnumImage, osMagnumFlavor, osMagnumFlavor, osExtGwID, osMagnumHTTPProxy, osMagnumHTTPSProxy, osMagnumNoProxy, osMagnumLabels, clusterName, nodeGroupName, nodeCount, osMagnumImage, osMagnumFlavor)
}

func testAccKubernetesV1NodeGroupOverrideLabels(keypairName, clusterTemplateName, clusterName string, nodeGroupName string, nodeCount int) string {
	return fmt.Sprintf(`
resource "nhncloud_compute_keypair_v2" "keypair_1" {
  name = "%s"
}

resource "nhncloud_kubernetes_clustertemplate_v1" "clustertemplate_1" {
  name                  = "%s"
  image                 = "%s"
  coe                   = "kubernetes"
  master_flavor         = "%s"
  flavor                = "%s"
  floating_ip_enabled   = true
  volume_driver         = "cinder"
  docker_storage_driver = "overlay2"
  docker_volume_size    = 5
  external_network_id   = "%s"
  network_driver        = "flannel"
  http_proxy            = "%s"
  https_proxy           = "%s"
  no_proxy              = "%s"
  labels = {
    kubescheduler_options = "log-flush-frequency=1m",
	%s
  }
}

resource "nhncloud_kubernetes_cluster_v1" "cluster_1" {
  name                 = "%s"
  cluster_template_id  = "${nhncloud_kubernetes_clustertemplate_v1.clustertemplate_1.id}"
  master_count         = 1
  node_count           = 1
  keypair              = "${nhncloud_compute_keypair_v2.keypair_1.name}"
}

resource "nhncloud_kubernetes_nodegroup_v1" "nodegroup_1" {
  name                 = "%s"
  cluster_id           = "${nhncloud_kubernetes_cluster_v1.cluster_1.id}"
  node_count           = %d
  docker_volume_size   = 10
  role				   = "myRole"
  min_node_count       = 1
  max_node_count       = 5
  image_id             = "%s"
  flavor_id            = "%s"
  merge_labels         = false
  labels = {
	kubescheduler_options = "log-flush-frequency=2m",
	%s
  }
}
`, keypairName, clusterTemplateName, osMagnumImage, osMagnumFlavor, osMagnumFlavor, osExtGwID, osMagnumHTTPProxy, osMagnumHTTPSProxy, osMagnumNoProxy, osMagnumLabels, clusterName, nodeGroupName, nodeCount, osMagnumImage, osMagnumFlavor, osMagnumLabels)
}
