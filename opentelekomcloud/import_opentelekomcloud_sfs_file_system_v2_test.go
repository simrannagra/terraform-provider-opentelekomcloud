package opentelekomcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

// PASS
func TestAccOTCSFSFileSharingV2_importBasic(t *testing.T) {
	resourceName := "opentelekomcloud_sfs_file_sharing_v2.sfs_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOTCSFSFileSharingV2Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSFSFileSharingV2_basic,
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}