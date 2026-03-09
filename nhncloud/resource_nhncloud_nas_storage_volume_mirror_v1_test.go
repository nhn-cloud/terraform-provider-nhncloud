package nhncloud

import (
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nhn-cloud/nhncloud.gophercloud/nhncloud/nas/v1/volumes"
)

func TestAccNasStorageVolumeMirrorV1_basic(t *testing.T) {
	var mirror volumes.Mirror

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckNasStorageV1(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNasStorageVolumeMirrorV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNasStorageVolumeMirrorV1Basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNasStorageVolumeMirrorV1Exists("nhncloud_nas_storage_volume_mirror_v1.mirror_1", &mirror),
					resource.TestCheckResourceAttr("nhncloud_nas_storage_volume_mirror_v1.mirror_1", "role", "SOURCE"),
					resource.TestCheckResourceAttr("nhncloud_nas_storage_volume_mirror_v1.mirror_1", "direction", "FORWARD"),
					resource.TestCheckResourceAttrSet("nhncloud_nas_storage_volume_mirror_v1.mirror_1", "dst_project_id"),
					resource.TestCheckResourceAttr("nhncloud_nas_storage_volume_mirror_v1.mirror_1", "dst_region", osNasStorageDstRegionName),
					resource.TestCheckResourceAttr("nhncloud_nas_storage_volume_mirror_v1.mirror_1", "dst_tenant_id", osNasStorageDstTenantID),
					resource.TestCheckResourceAttrSet("nhncloud_nas_storage_volume_mirror_v1.mirror_1", "dst_volume_id"),
					resource.TestCheckResourceAttr("nhncloud_nas_storage_volume_mirror_v1.mirror_1", "dst_volume_name", "volume_mirror_1"),
					resource.TestCheckResourceAttrSet("nhncloud_nas_storage_volume_mirror_v1.mirror_1", "src_project_id"),
					resource.TestCheckResourceAttr("nhncloud_nas_storage_volume_mirror_v1.mirror_1", "src_region", osRegionName),
					resource.TestCheckResourceAttrSet("nhncloud_nas_storage_volume_mirror_v1.mirror_1", "src_tenant_id"),
					resource.TestCheckResourceAttr("nhncloud_nas_storage_volume_mirror_v1.mirror_1", "src_volume_name", "volume_1"),
					resource.TestCheckResourceAttrSet("nhncloud_nas_storage_volume_mirror_v1.mirror_1", "created_at"),
					resource.TestCheckResourceAttr("nhncloud_nas_storage_volume_mirror_v1.mirror_1", "dst_volume.0.name", "volume_mirror_1"),
					resource.TestCheckResourceAttr("nhncloud_nas_storage_volume_mirror_v1.mirror_1", "dst_volume.0.size_gb", "300"),
					resource.TestCheckResourceAttr("nhncloud_nas_storage_volume_mirror_v1.mirror_1", "dst_volume.0.mount_protocol.0.protocol", "nfs"),
				),
			},
			{
				Config: testAccNasStorageVolumeMirrorV1BasicUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNasStorageVolumeMirrorV1Exists("nhncloud_nas_storage_volume_mirror_v1.mirror_1", &mirror),
					resource.TestCheckResourceAttr("nhncloud_nas_storage_volume_mirror_v1.mirror_1", "dst_volume.0.description", "volume_mirror_1_description"),
				),
			},
		},
	})

	// The volume created by mirror must be deleted manually.
	err := removeVolumeCreatedByMirror(mirror)
	if err != nil {
		t.Logf("[WARNING] failed to remove the volume created by mirror: %s", err)
	}
}

func testAccCheckNasStorageVolumeMirrorV1Destroy(s *terraform.State) error {
	config, err := testAccAuthFromEnv()
	if err != nil {
		return err
	}

	nasStorageClient, err := config.NasStorageV1Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating NHN Cloud NAS storage client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nhncloud_nas_storage_volume_mirror_v1" {
			continue
		}

		volume, err := volumes.Get(nasStorageClient, rs.Primary.Attributes["volume_id"]).Extract()
		if err == nil && len(volume.Mirrors) > 0 {
			return fmt.Errorf("Volume mirror still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckNasStorageVolumeMirrorV1Exists(n string, mirror *volumes.Mirror) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		volumeID := rs.Primary.Attributes["src_volume_id"]

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

		index := slices.IndexFunc(volume.Mirrors, func(m *volumes.Mirror) bool {
			return m.ID == rs.Primary.ID
		})
		if index == -1 {
			return fmt.Errorf("Volume mirror not found")
		}

		*mirror = *volume.Mirrors[index]
		return nil
	}
}

func removeVolumeCreatedByMirror(mirror volumes.Mirror) error {
	config, err := testAccAuthFromEnv()
	if err != nil {
		return err
	}

	nasStorageClient, err := config.NasStorageV1Client(mirror.DstRegion)
	if err != nil {
		return fmt.Errorf("Error creating NHN Cloud NAS storage client: %s", err)
	}

	err = volumes.Delete(nasStorageClient, mirror.DstVolumeID).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting NHN Cloud NAS storage volume %s: %s", mirror.DstVolumeID, err)
	}

	timer := time.NewTimer(3 * time.Minute)
	defer timer.Stop()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

Loop:
	for {
		select {
		case <-ticker.C:
			_, err := volumes.Get(nasStorageClient, mirror.DstVolumeID).Extract()
			if err == nil {
				continue
			}

			if _, ok := err.(gophercloud.ErrDefault404); ok {
				break Loop
			}

			return fmt.Errorf("Error getting NHN Cloud NAS storage volume %s: %s", mirror.DstVolumeID, err)
		case <-timer.C:
			return fmt.Errorf("Deletion is taking too long, proceeding without waiting")
		}

	}
	return nil
}

func testAccNasStorageVolumeMirrorV1Basic() string {
	return fmt.Sprintf(`
		resource "nhncloud_nas_storage_volume_v1" "volume_1" {
			name = "volume_1"
			size_gb = 300

			mount_protocol {
				protocol = "nfs"
			}
		}

		resource "nhncloud_nas_storage_volume_mirror_v1" "mirror_1" {
			src_region = "%s"
			src_volume_id = nhncloud_nas_storage_volume_v1.volume_1.id
			dst_region = "%s"
			dst_tenant_id = "%s"

			dst_volume {
				name = "volume_mirror_1"
				size_gb = 300

				mount_protocol {
					protocol = "nfs"
				}
			}
		}
	`, osRegionName, osNasStorageDstRegionName, osNasStorageDstTenantID)
}

func testAccNasStorageVolumeMirrorV1BasicUpdate() string {
	return fmt.Sprintf(`
		resource "nhncloud_nas_storage_volume_v1" "volume_1" {
			name = "volume_1"
			size_gb = 300

			mount_protocol {
				protocol = "nfs"
			}
		}

		resource "nhncloud_nas_storage_volume_mirror_v1" "mirror_1" {
			src_region = "%s"
			src_volume_id = nhncloud_nas_storage_volume_v1.volume_1.id
			dst_region = "%s"
			dst_tenant_id = "%s"

			dst_volume {
				name = "volume_mirror_1"
				description = "volume_mirror_1_description"
				size_gb = 300

				mount_protocol {
					protocol = "nfs"
				}
			}
		}
	`, osRegionName, osNasStorageDstRegionName, osNasStorageDstTenantID)
}
