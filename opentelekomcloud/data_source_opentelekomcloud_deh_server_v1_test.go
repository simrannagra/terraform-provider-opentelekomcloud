package opentelekomcloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

// PASS
func TestAccOTCDedicatedHostServerV1DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOTCDedicatedHostServerV1DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDedicatedHostServerV1DataSourceID("data.opentelekomcloud_deh_server_v1.servers"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_deh_server_v1.servers", "server_name", "ecs-Deh-c2c"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_deh_server_v1.servers", "status", "ACTIVE"),
				),
			},
		},
	})
}

func testAccCheckDedicatedHostServerV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find servers data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("server data source ID not set ")
		}

		return nil
	}
}

var testAccOTCDedicatedHostServerV1DataSource_basic = `
data "opentelekomcloud_deh_server_v1" "servers" {
  id = "dcaa399d-0d72-42a0-8a55-0030d25828cc"
  server_id = "eb18f8a6-c51c-4509-a1a7-e298246a7352"
}
`
