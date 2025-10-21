package nhncloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccKubernetesV1NodeGroupImport_basic(t *testing.T) {
	resourceName := "nhncloud_kubernetes_nodegroup_v1.nodegroup_1"
	clusterName := acctest.RandomWithPrefix("tf-acc-cluster")
	keypairName := acctest.RandomWithPrefix("tf-acc-keypair")
	clusterTemplateName := acctest.RandomWithPrefix("tf-acc-clustertemplate")
	nodeGroupName := acctest.RandomWithPrefix("tf-acc-cluster")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckContainerInfra(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckKubernetesV1NodeGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesV1NodeGroupBasic(keypairName, clusterTemplateName, clusterName, nodeGroupName, 1),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccKubernetesV1NodeGroupImport_mergeLabels(t *testing.T) {
	resourceName := "nhncloud_kubernetes_nodegroup_v1.nodegroup_1"
	clusterName := acctest.RandomWithPrefix("tf-acc-cluster")
	keypairName := acctest.RandomWithPrefix("tf-acc-keypair")
	clusterTemplateName := acctest.RandomWithPrefix("tf-acc-clustertemplate")
	nodeGroupName := acctest.RandomWithPrefix("tf-acc-cluster")

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
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"merge_labels", "labels"},
			},
		},
	})
}

func TestAccKubernetesV1NodeGroupImport_overrideLabels(t *testing.T) {
	resourceName := "nhncloud_kubernetes_nodegroup_v1.nodegroup_1"
	clusterName := acctest.RandomWithPrefix("tf-acc-cluster")
	keypairName := acctest.RandomWithPrefix("tf-acc-keypair")
	clusterTemplateName := acctest.RandomWithPrefix("tf-acc-clustertemplate")
	nodeGroupName := acctest.RandomWithPrefix("tf-acc-cluster")

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
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"merge_labels", "labels"},
			},
		},
	})
}
