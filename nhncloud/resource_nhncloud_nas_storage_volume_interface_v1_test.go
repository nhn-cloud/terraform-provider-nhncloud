package nhncloud

import (
	"fmt"
	"slices"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nhn-cloud/nhncloud.gophercloud/nhncloud/nas/v1/volumes"
)

func TestAccNasStorageVolumeInterfaceV1_basic(t *testing.T) {
	var vInterface volumes.Interface

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckNasStorageV1(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNasStorageVolumeInterfaceV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNasStorageVolumeInterfaceV1Basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNasStorageVolumeInterfaceV1Exists("nhncloud_nas_storage_volume_interface_v1.vinterface_1", &vInterface),
					resource.TestCheckResourceAttrSet("nhncloud_nas_storage_volume_interface_v1.vinterface_1", "volume_id"),
					resource.TestCheckResourceAttr("nhncloud_nas_storage_volume_interface_v1.vinterface_1", "subnet_id", osSubnetID),
					resource.TestCheckResourceAttrSet("nhncloud_nas_storage_volume_interface_v1.vinterface_1", "id"),
					resource.TestCheckResourceAttrSet("nhncloud_nas_storage_volume_interface_v1.vinterface_1", "path"),
					resource.TestCheckResourceAttrSet("nhncloud_nas_storage_volume_interface_v1.vinterface_1", "tenant_id"),
				),
			},
		},
	})
}

func testAccCheckNasStorageVolumeInterfaceV1Destroy(s *terraform.State) error {
	config, err := testAccAuthFromEnv()
	if err != nil {
		return err
	}

	nasStorageClient, err := config.NasStorageV1Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating NHN Cloud NAS storage client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nhncloud_nas_storage_volume_interface_v1" {
			continue
		}

		volume, err := volumes.Get(nasStorageClient, rs.Primary.Attributes["volume_id"]).Extract()
		if err == nil && len(volume.Interfaces) > 0 {
			return fmt.Errorf("Volume interface still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckNasStorageVolumeInterfaceV1Exists(n string, vInterface *volumes.Interface) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		volumeID := rs.Primary.Attributes["volume_id"]

		config, err := testAccAuthFromEnv()
		if err != nil {
			return err
		}

		nasStorageClient, err := config.NasStorageV1Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating NHN Cloud NAS storage client: %s", err)
		}

		volume, err := volumes.Get(nasStorageClient, volumeID).Extract()
		if err != nil {
			return fmt.Errorf("Error getting NHN Cloud NAS storage volume %s: %s", volumeID, err)
		}

		index := slices.IndexFunc(volume.Interfaces, func(i *volumes.Interface) bool {
			return i.ID == rs.Primary.ID
		})
		if index == -1 {
			return fmt.Errorf("Volume interface not found")
		}

		*vInterface = *volume.Interfaces[index]

		return nil
	}
}

func testAccNasStorageVolumeInterfaceV1Basic() string {
	return fmt.Sprintf(`
		resource "nhncloud_nas_storage_volume_v1" "volume_1" {
			name = "volume_1"
			size_gb = 300

			mount_protocol {
				protocol = "nfs"
			}
		}

		resource "nhncloud_nas_storage_volume_interface_v1" "vinterface_1" {
			volume_id = nhncloud_nas_storage_volume_v1.volume_1.id
			subnet_id = "%s"
		}
	`, osSubnetID)
}
