package opentelekomcloud

import (
	"fmt"
	"log"

	"github.com/huaweicloud/golangsdk/openstack/deh/v1/hosts"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceDEHV1() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDEHV1Read,

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"host_type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"host_type_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"flavor": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"availability_zone": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"limit": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"marker": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"auto_placement": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"available_vcpus": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"available_memory": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cores": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"sockets": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"instance_total": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceDEHV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	dehClient, err := config.dehV1Client(GetRegion(d, config))

	listOpts := hosts.ListOpts{
		ID:    d.Get("id").(string),
		Name:  d.Get("name").(string),
		State: d.Get("state").(string),
		Az:    d.Get("availability_zone").(string),
	}

	refinedDeh, err := hosts.List(dehClient, listOpts)
	if err != nil {
		return fmt.Errorf("Unable to retrieve dedicated hosts: %s", err)
	}

	if len(refinedDeh) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedDeh) > 1 {
		return fmt.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	Deh := refinedDeh[0]

	log.Printf("[INFO] Retrieved Deh using given filter %s: %+v", Deh.ID, Deh)
	d.SetId(Deh.ID)

	d.Set("name", Deh.Name)
	d.Set("id", Deh.ID)
	d.Set("auto_placement", Deh.AutoPlacement)
	d.Set("availability_zone", Deh.Az)
	d.Set("tenant_id", Deh.TenantId)
	d.Set("state", Deh.State)
	d.Set("available_vcpus", Deh.AvailableVcpus)
	d.Set("available_memory", Deh.AvailableMemory)
	d.Set("instance_total", Deh.InstanceTotal)
	//d.Set("instance_uuids", Deh.InstanceUuids[0])
	d.Set("host_type_name", Deh.HostProperties.HostTypeName)
	d.Set("host_type", Deh.HostProperties.HostType)
	d.Set("cores", Deh.HostProperties.Cores)
	d.Set("sockets", Deh.HostProperties.Sockets)
	d.Set("flavor", Deh.HostProperties.AvailableInstanceCapacities[0].Flavor)
	d.Set("region", GetRegion(d, config))

	return nil
}
