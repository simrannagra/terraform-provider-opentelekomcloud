package opentelekomcloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

// PASS
func TestAccOTCDedicatedHostV1DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOTCDedicatedHostV1DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDedicatedHostV1DataSourceID("data.opentelekomcloud_deh_v1.hosts"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_deh_v1.hosts", "name", "c2c-test"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_deh_v1.hosts", "availability_zone", "eu-de-01"),
				),
			},
		},
	})
}

func testAccCheckDedicatedHostV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find deh data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("deh data source ID not set ")
		}

		return nil
	}
}

var testAccOTCDedicatedHostV1DataSource_basic = `
data "opentelekomcloud_deh_v1" "hosts" {
  id = "95baaeb7-d933-440e-ad43-62b59b723aee"
}
`
