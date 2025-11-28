package nhncloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nhn-cloud/nhncloud.gophercloud/nhncloud/nas/v1/volumes"
)

func TestAccNasStorageVolumeV1_basic(t *testing.T) {
	var volume volumes.Volume

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckNasStorageV1(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNasStorageVolumeV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNasStorageVolumeV1Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNasStorageVolumeV1Exists("nhncloud_nas_storage_volume_v1.volume_1", &volume),
					resource.TestCheckResourceAttr("nhncloud_nas_storage_volume_v1.volume_1", "name", "volume_1"),
					resource.TestCheckResourceAttr("nhncloud_nas_storage_volume_v1.volume_1", "size_gb", "300"),
					resource.TestCheckResourceAttr("nhncloud_nas_storage_volume_v1.volume_1", "mount_protocol.0.protocol", "nfs"),
					resource.TestCheckResourceAttrSet("nhncloud_nas_storage_volume_v1.volume_1", "id"),
					resource.TestCheckResourceAttrSet("nhncloud_nas_storage_volume_v1.volume_1", "project_id"),
					resource.TestCheckResourceAttrSet("nhncloud_nas_storage_volume_v1.volume_1", "tenant_id"),
					resource.TestCheckResourceAttrSet("nhncloud_nas_storage_volume_v1.volume_1", "created_at"),
					resource.TestCheckResourceAttrSet("nhncloud_nas_storage_volume_v1.volume_1", "updated_at"),
				),
			},
			{
				Config: testAccNasStorageVolumeV1Update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNasStorageVolumeV1Exists("nhncloud_nas_storage_volume_v1.volume_1", &volume),
					resource.TestCheckResourceAttr("nhncloud_nas_storage_volume_v1.volume_1", "description", "test volume updated"),
					resource.TestCheckResourceAttr("nhncloud_nas_storage_volume_v1.volume_1", "size_gb", "400"),
					resource.TestCheckResourceAttr("nhncloud_nas_storage_volume_v1.volume_1", "acl.0", "0.0.0.0/0"),
					resource.TestCheckResourceAttr("nhncloud_nas_storage_volume_v1.volume_1", "mount_protocol.0.protocol", "nfs"),
					resource.TestCheckResourceAttr("nhncloud_nas_storage_volume_v1.volume_1", "snapshot_policy.0.max_scheduled_count", "10"),
					resource.TestCheckResourceAttr("nhncloud_nas_storage_volume_v1.volume_1", "snapshot_policy.0.reserve_percent", "10"),
					resource.TestCheckResourceAttr("nhncloud_nas_storage_volume_v1.volume_1", "snapshot_policy.0.schedule.0.time", "00:00"),
					resource.TestCheckResourceAttr("nhncloud_nas_storage_volume_v1.volume_1", "snapshot_policy.0.schedule.0.time_offset", "+09:00"),
				),
			},
			{
				Config: testAccNasStorageVolumeV1Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNasStorageVolumeV1Exists("nhncloud_nas_storage_volume_v1.volume_1", &volume),
					resource.TestCheckResourceAttr("nhncloud_nas_storage_volume_v1.volume_1", "name", "volume_1"),
					resource.TestCheckResourceAttr("nhncloud_nas_storage_volume_v1.volume_1", "description", ""),
					resource.TestCheckResourceAttr("nhncloud_nas_storage_volume_v1.volume_1", "size_gb", "300"),
					resource.TestCheckResourceAttr("nhncloud_nas_storage_volume_v1.volume_1", "mount_protocol.0.protocol", "nfs"),
					resource.TestCheckResourceAttr("nhncloud_nas_storage_volume_v1.volume_1", "snapshot_policy.0.max_scheduled_count", "0"),
					resource.TestCheckResourceAttr("nhncloud_nas_storage_volume_v1.volume_1", "snapshot_policy.0.reserve_percent", "5"),
					resource.TestCheckResourceAttr("nhncloud_nas_storage_volume_v1.volume_1", "snapshot_policy.0.schedule.#", "0"),
				),
			},
		},
	})
}

func testAccCheckNasStorageVolumeV1Destroy(s *terraform.State) error {
	config, err := testAccAuthFromEnv()
	if err != nil {
		return err
	}

	nasStorageClient, err := config.NasStorageV1Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating NHN Cloud NAS storage client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nhncloud_nas_storage_volume_v1" {
			continue
		}

		_, err := volumes.GetVolume(nasStorageClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Volume still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckNasStorageVolumeV1Exists(n string, volume *volumes.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config, err := testAccAuthFromEnv()
		if err != nil {
			return err
		}

		nasStorageClient, err := config.NasStorageV1Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating NHN Cloud NAS storage client: %s", err)
		}

		found, err := volumes.GetVolume(nasStorageClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Volume not found")
		}

		*volume = *found

		return nil
	}
}

const testAccNasStorageVolumeV1Basic = `
	resource "nhncloud_nas_storage_volume_v1" "volume_1" {
		name = "volume_1"
		size_gb = 300

		mount_protocol {
			protocol = "nfs"
		}
	}
`

const testAccNasStorageVolumeV1Update = `
	resource "nhncloud_nas_storage_volume_v1" "volume_1" {
		name = "volume_1"
		description = "test volume updated"
		size_gb = 400

		acl = ["0.0.0.0/0"]

		mount_protocol {
			protocol = "nfs"
		}

		snapshot_policy {
			max_scheduled_count = 10
			reserve_percent = 10
			schedule {
				time = "00:00"
				time_offset = "+09:00"
				weekdays = [1, 2, 3, 4, 5]
			}
		}
	}
`
