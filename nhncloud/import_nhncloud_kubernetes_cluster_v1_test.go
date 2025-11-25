package nhncloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccKubernetesV1ClusterImport_basic(t *testing.T) {
	resourceName := "nhncloud_kubernetes_cluster_v1.cluster_1"
	clusterName := acctest.RandomWithPrefix("tf-acc-cluster")
	keypairName := acctest.RandomWithPrefix("tf-acc-keypair")
	clusterTemplateName := acctest.RandomWithPrefix("tf-acc-clustertemplate")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckKubernetes(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckKubernetesV1ClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesV1ClusterBasic(keypairName, clusterTemplateName, clusterName, 1),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"kubeconfig"},
			},
		},
	})
}

func TestAccKubernetesV1ClusterImport_mergeLabels(t *testing.T) {
	resourceName := "nhncloud_kubernetes_cluster_v1.cluster_1"
	clusterName := acctest.RandomWithPrefix("tf-acc-cluster")
	keypairName := acctest.RandomWithPrefix("tf-acc-keypair")
	clusterTemplateName := acctest.RandomWithPrefix("tf-acc-clustertemplate")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckKubernetes(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckKubernetesV1ClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesV1ClusterLabels(keypairName, clusterTemplateName, clusterName, 1, true),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"kubeconfig", "merge_labels", "labels"},
			},
		},
	})
}

func TestAccKubernetesV1ClusterImport_overrideLabels(t *testing.T) {
	resourceName := "nhncloud_kubernetes_cluster_v1.cluster_1"
	clusterName := acctest.RandomWithPrefix("tf-acc-cluster")
	keypairName := acctest.RandomWithPrefix("tf-acc-keypair")
	clusterTemplateName := acctest.RandomWithPrefix("tf-acc-clustertemplate")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckKubernetes(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckKubernetesV1ClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesV1ClusterLabels(keypairName, clusterTemplateName, clusterName, 1, false),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"kubeconfig", "merge_labels", "labels"},
			},
		},
	})
}
