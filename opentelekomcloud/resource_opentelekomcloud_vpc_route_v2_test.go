package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/huaweicloud/golangsdk/openstack/networking/v2/routes"
)

// PASS
func TestAccOTCVpcRouteV2_basic(t *testing.T) {
	var route routes.Route

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOTCRouteV2Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRouteV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOTCRouteV2Exists("opentelekomcloud_vpc_route_v2.route_1", &route),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vpc_route_v2.route_1", "destination", "192.168.0.0/16"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vpc_route_v2.route_1", "type", "peering"),
				),
			},
		},
	})
}

// PASS
func TestAccOTCVpcRouteV2_timeout(t *testing.T) {
	var route routes.Route

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOTCRouteV2Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRouteV2_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOTCRouteV2Exists("opentelekomcloud_vpc_route_v2.route_1", &route),
				),
			},
		},
	})
}

func testAccCheckOTCRouteV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	routeClient, err := config.networkingHwV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud route client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_vpc_route_v2" {
			continue
		}

		_, err := routes.Get(routeClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Route still exists")
		}
	}

	return nil
}

func testAccCheckOTCRouteV2Exists(n string, route *routes.Route) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		routeClient, err := config.networkingHwV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenTelekomCloud route client: %s", err)
		}

		found, err := routes.Get(routeClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.RouteID != rs.Primary.ID {
			return fmt.Errorf("route not found")
		}

		*route = *found

		return nil
	}
}

const testAccRouteV2_basic = `
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name = "vpc_test"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_v1" "vpc_2" {
  name = "vpc_test1"
  cidr = "192.168.0.0/16"
}
resource "opentelekomcloud_vpc_peering_connection_v2" "peering_1" {
  name = "opentelekomcloud_peering"
  vpc_id = "${opentelekomcloud_vpc_v1.vpc_1.id}"
  peer_vpc_id = "${opentelekomcloud_vpc_v1.vpc_2.id}"
}
resource "opentelekomcloud_vpc_route_v2" "route_1" {
  type = "peering"
  nexthop = "${opentelekomcloud_vpc_peering_connection_v2.peering_1.id}"
  destination = "192.168.0.0/16"
  vpc_id ="${opentelekomcloud_vpc_v1.vpc_1.id}"

}
`

const testAccRouteV2_timeout = `
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name = "vpc_test"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_v1" "vpc_2" {
  name = "vpc_test1"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_peering_connection_v2" "peering_1" {
  name = "opentelekomcloud_peering"
  vpc_id = "${opentelekomcloud_vpc_v1.vpc_1.id}"
  peer_vpc_id = "${opentelekomcloud_vpc_v1.vpc_2.id}"
}

resource "opentelekomcloud_vpc_route_v2" "route_1" {
   type = "peering"
  nexthop = "${opentelekomcloud_vpc_peering_connection_v2.peering_1.id}"
  destination = "192.168.0.0/16"
  vpc_id ="${opentelekomcloud_vpc_v1.vpc_1.id}"

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`
