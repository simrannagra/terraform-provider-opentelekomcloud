package opentelekomcloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/deh/v1/hosts"
	"log"
	"time"
)

func resourceDeHV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceDeHV1Create,
		Read:   resourceDeHV1Read,
		Update: resourceDeHV1Update,
		Delete: resourceDeHV1Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"dedicated_host_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"auto_placement": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"availability_zone": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"host_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"quantity": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"state": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"available_vcpus": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"available_memory": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"allocated_at": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"released_at": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_total": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_uuids": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"host_type_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vcpus": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cores": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"sockets": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"memory": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"available_instance_capacities": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"flavor": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceDeHV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	dehClient, err := config.dehV1Client(GetRegion(d, config))

	log.Printf("[DEBUG] Value of DeH Client: %#v", dehClient)

	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomcomCloud DeH Client: %s", err)
	}

	allocateOpts := hosts.AllocateOpts{
		Name:             d.Get("name").(string),
		HostType:         d.Get("host_type").(string),
		AutoPlacement:    d.Get("auto_placement").(string),
		AvailabilityZone: d.Get("availability_zone").(string),
		Quantity:         d.Get("quantity").(int),
	}

	log.Printf("[DEBUG] Create Options: %#v", allocateOpts)

	allocate, err := hosts.Allocate(dehClient, allocateOpts).Extract()

	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomcomCloud Dedicated Host: %s", err)
	}
	d.SetId(allocate.AllocatedHostIds[0])
	log.Printf("[INFO] Host ID: %s", allocate.AllocatedHostIds[0])

	log.Printf("[DEBUG] Waiting for OpenTelekomcomCloud Dedicated Host (%s) to become available", allocate.AllocatedHostIds[0])

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Creating"},
		Target:     []string{"Available"},
		Refresh:    waitForDeHActive(dehClient, allocate.AllocatedHostIds[0]),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, err = stateConf.WaitForState()

	return resourceDeHV1Read(d, meta)
}

func resourceDeHV1Read(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)
	dehClient, err := config.dehV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud DeH client: %s", err)
	}
	n, err := hosts.Get(dehClient, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving OpenTelekomCloud Dedicated Host: %s", err)
	}

	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving OpenTelekomCloud Dedicated Hosts: %s", err)
	}

	d.Set("id", n.ID)
	d.Set("name", n.Name)
	d.Set("state", n.State)
	d.Set("dedicated_host_id", n.ID)
	d.Set("auto_placement", n.AutoPlacement)
	d.Set("availability_zone", n.Az)
	d.Set("quantity", 1)
	d.Set("available_vcpus", n.AvailableVcpus)
	d.Set("available_memory", n.AvailableMemory)
	d.Set("allocated_at", n.AllocatedAt)
	d.Set("released_at", n.ReleasedAt)
	d.Set("instance_total", n.InstanceTotal)
	d.Set("instance_uuids", n.InstanceUuids)
	d.Set("host_type", n.HostProperties.HostType)
	d.Set("host_type_name", n.HostProperties.HostTypeName)
	d.Set("vcpus", n.HostProperties.Vcpus)
	d.Set("cores", n.HostProperties.Cores)
	d.Set("sockets", n.HostProperties.Sockets)
	d.Set("memory", n.HostProperties.Memory)
	d.Set("available_instance_capacities", getInstanceProperties(n))

	return nil
}

func resourceDeHV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	dehClient, err := config.dehV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud DeH Client: %s", err)
	}
	var updateOpts hosts.UpdateOpts

	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("auto_placement") {
		updateOpts.AutoPlacement = d.Get("auto_placement").(string)
	}

	_, err = hosts.Update(dehClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error updating OpenTelekomCloud Dedicated Host: %s", err)
	}
	return resourceDeHV1Read(d, meta)
}

func resourceDeHV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	dehClient, err := config.dehV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForDeHDelete(dehClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error deleting OpenTelekomCloud Dedicated Host : %s", err)
	}
	d.SetId("")
	return nil
}

func waitForDeHActive(dehClient *golangsdk.ServiceClient, dehID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := hosts.Get(dehClient, dehID).Extract()
		if err != nil {
			return nil, "", err
		}

		if n.State == "available" {
			return n, "ACTIVE", nil
		}

		if n.State == "unavailable" {
			return nil, "", fmt.Errorf("Dedicated Host state: '%s'", n.State)
		}

		return n, n.State, nil
	}
}

func waitForDeHDelete(dehClient *golangsdk.ServiceClient, dehID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		r, err := hosts.Get(dehClient, dehID).Extract()

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted OpenTelekomCloud Dedicated Host  %s", dehID)
				return r, "DELETED", nil
			}
			return r, "ACTIVE", err
		}
		result := hosts.Delete(dehClient, dehID)

		if result.Err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted OpenTelekomCloud Dedicated Host %s", dehID)
				return r, "DELETED", nil
			}
			if errCode, ok := err.(golangsdk.ErrUnexpectedResponseCode); ok {
				if errCode.Actual == 409 {
					return r, "ACTIVE", nil
				}
			}
			return r, "ACTIVE", err
		}

		return r, "ACTIVE", nil
	}
}

func getInstanceProperties(n *hosts.Host) []map[string]interface{} {
	var v []map[string]interface{}
	for _, val := range n.HostProperties.AvailableInstanceCapacities {
		mapping := map[string]interface{}{
			"flavor": val.Flavor,
		}
		v = append(v, mapping)
	}
	return v
}
