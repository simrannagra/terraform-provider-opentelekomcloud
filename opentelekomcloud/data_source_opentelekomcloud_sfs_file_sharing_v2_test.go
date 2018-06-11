package opentelekomcloud

import (
	"fmt"
	"testing"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

// PASS
func TestAccOpenTelekomCloudSFSFileSharingV2DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOpenTelekomCloudSFSFileSharingV2DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSFileSharingV2DataSourceID("data.opentelekomcloud_sfs_file_sharing_v2.shares"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_sfs_file_sharing_v2.shares", "name", "sfs-c2c"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_sfs_file_sharing_v2.shares", "status", "available"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_sfs_file_sharing_v2.shares", "size", "1"),
				),
			},
		},
	})
}

func testAccCheckSFSFileSharingV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find share file data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("share file data source ID not set ")
		}

		return nil
	}
}

var testAccOpenTelekomCloudSFSFileSharingV2DataSource_basic = `
data "opentelekomcloud_sfs_file_sharing_v2" "shares" {
  id = "051ab809-e43a-4db5-a1d6-d85c298054d8"
}
`

