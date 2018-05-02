package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/huaweicloud/golangsdk/openstack/rts/v1/stacks"
)

// PASS
func TestAccOTCRtsStackV1_basic(t *testing.T) {
	var stacks stacks.RetrievedStack

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOTCCRtsStackV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRtsStackV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOTCRtsStackV1Exists("opentelekomcloud_sfs_stack_v1.stack_1", &stacks),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_sfs_stack_v1.stack_1", "name", "terraform_provider_stack"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_sfs_stack_v1.stack_1", "status", "CREATE_COMPLETE"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_sfs_stack_v1.stack_1", "description", "Simple template"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_sfs_stack_v1.stack_1", "status_reason", "Stack CREATE completed successfully"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_sfs_stack_v1.stack_1", "disable_rollback", "true"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_sfs_stack_v1.stack_1", "timeout_mins", "60"),


				),

			},
		},
	})
}

func TestAccOTCRtsStackV1_update(t *testing.T) {
	var stacks stacks.RetrievedStack

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOTCCRtsStackV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRtsStackV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOTCRtsStackV1Exists("opentelekomcloud_sfs_stack_v1.stack_1", &stacks),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_sfs_stack_v1.stack_1", "name", "terraform_provider_stack"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_sfs_stack_v1.stack_1", "status", "CREATE_COMPLETE"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_sfs_stack_v1.stack_1", "description", "Simple template"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_sfs_stack_v1.stack_1", "status_reason", "Stack CREATE completed successfully"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_sfs_stack_v1.stack_1", "disable_rollback", "true"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_sfs_stack_v1.stack_1", "timeout_mins", "60"),
				),
			},
			resource.TestStep{
				Config: testAccRtsStackV1_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOTCRtsStackV1Exists("opentelekomcloud_sfs_stack_v1.stack_1", &stacks),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_sfs_stack_v1.stack_1", "disable_rollback", "false"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_sfs_stack_v1.stack_1", "timeout_mins", "50"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_sfs_stack_v1.stack_1", "status", "UPDATE_COMPLETE"),

				),
			},
		},
	})
}

// PASS
func TestAccOTCRtsStackV1_timeout(t *testing.T) {
	var stacks stacks.RetrievedStack

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOTCRouteV2Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRtsStackV1_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOTCRtsStackV1Exists("opentelekomcloud_sfs_stack_v1.stack_1", &stacks),
				),
			},
		},
	})
}

func testAccCheckOTCCRtsStackV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	orchestrationClient, err := config.orchestrationV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud orchestration client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_sfs_stack_v1" {
			continue
		}

		_, err := stacks.Get(orchestrationClient,rs.Primary.Attributes["name"] ,rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Stack still exists %s",err)
		}
	}

	return nil
}

func testAccCheckOTCRtsStackV1Exists(n string, stack *stacks.RetrievedStack) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		orchestrationClient, err := config.orchestrationV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenTelekomCloud orchestration Client : %s", err)
		}

		found, err := stacks.Get(orchestrationClient, rs.Primary.Attributes["name"],rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("stack not found")
		}

		*stack = *found

		return nil
	}
}

const testAccRtsStackV1_basic = `
resource "opentelekomcloud_sfs_stack_v1" "stack_1" {
  name = "terraform_provider_stack"
  disable_rollback= true
  timeout_mins=60
  template = <<JSON
          {
			 "heat_template_version": "2013-05-23",
			 "description": "Simple template ",
			 "parameters": {
			    "image_id": {
			        "type": "string",
                    "description": "Image to be used for compute instance",
			        "label": "Image ID",
                    "default": "ea67839e-fd7a-4b99-9f81-13c4c8dc317c"
			    },
			    "net_id": {
			        "type": "string",
			        "description": "The network to be used",
			        "label": "Network UUID",
                    "default": "7eb54ab6-5cdb-446a-abbe-0dda1885c76e"
			    },
			    "instance_type": {
			        "type": "string",
			        "description": "Type of instance (flavor) to be used",
			        "label": "Instance Type",
                    "default": "s1.medium"
			      }
			  },
			 "resources": {
			    "my_instance": {
                  "type": "OS::Nova::Server",
			      "properties": {
			      "image": {
			      "get_param": "image_id"
			      },
			      "flavor": {
                  "get_param": "instance_type"
			      },
			    "networks": [
			     {
			        "network": {
			        "get_param": "net_id"
			      }
			     }
			    ]
			   }
			  }
			 }
			}
JSON

}
`

const testAccRtsStackV1_update = `
resource "opentelekomcloud_sfs_stack_v1" "stack_1" {
  name = "terraform_provider_stack"
  disable_rollback= false
  timeout_mins=50
  template = <<JSON
          {
			 "heat_template_version": "2013-05-23",
			 "description": "Simple template ",
			 "parameters": {
			    "image_id": {
			        "type": "string",
                    "description": "Image to be used for compute instance",
			        "label": "Image ID",
                    "default": "ea67839e-fd7a-4b99-9f81-13c4c8dc317c"
			    },
			    "net_id": {
			        "type": "string",
			        "description": "The network to be used",
			        "label": "Network UUID",
                    "default": "7eb54ab6-5cdb-446a-abbe-0dda1885c76e"
			    },
			    "instance_type": {
			        "type": "string",
			        "description": "Type of instance (flavor) to be used",
			        "label": "Instance Type",
                    "default": "s1.medium"
			      }
			  },
			 "resources": {
			    "my_instance": {
                  "type": "OS::Nova::Server",
			      "properties": {
			      "image": {
			      "get_param": "image_id"
			      },
			      "flavor": {
                  "get_param": "instance_type"
			      },
			    "networks": [
			     {
			        "network": {
			        "get_param": "net_id"
			      }
			     }
			    ]
			   }
			  }
			 }
			}
JSON

}
`
const testAccRtsStackV1_timeout = `
resource "opentelekomcloud_sfs_stack_v1" "stack_1" {
  name = "terraform_provider_stack"
  disable_rollback= true
  timeout_mins=60

  template = <<JSON
          {
			 "heat_template_version": "2013-05-23",
			 "description": "Simple template ",
			 "parameters": {
			    "image_id": {
			        "type": "string",
                    "description": "Image to be used for compute instance",
			        "label": "Image ID",
                    "default": "ea67839e-fd7a-4b99-9f81-13c4c8dc317c"
			    },
			    "net_id": {
			        "type": "string",
			        "description": "The network to be used",
			        "label": "Network UUID",
                    "default": "7eb54ab6-5cdb-446a-abbe-0dda1885c76e"
			    },
			    "instance_type": {
			        "type": "string",
			        "description": "Type of instance (flavor) to be used",
			        "label": "Instance Type",
                    "default": "s1.medium"
			      }
			  },
			 "resources": {
			    "my_instance": {
                  "type": "OS::Nova::Server",
			      "properties": {
			      "image": {
			      "get_param": "image_id"
			      },
			      "flavor": {
                  "get_param": "instance_type"
			      },
			    "networks": [
			     {
			        "network": {
			        "get_param": "net_id"
			      }
			     }
			    ]
			   }
			  }
			 }
			}
JSON

}

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`
