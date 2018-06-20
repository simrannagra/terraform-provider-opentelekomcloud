package opentelekomcloud

import (
	"fmt"
	"log"

	"github.com/huaweicloud/golangsdk/openstack/deh/v1/hosts"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceDEHServersV1() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDEHServersV1Read,

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
			"server_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"server_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"server_flavor": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},
			"metadata": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"limit": &schema.Schema{
				Type:     schema.TypeInt,
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
		},
	}
}

func dataSourceDEHServersV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	dehClient, err := config.dehV1Client(GetRegion(d, config))

	listServerOpts := hosts.ListServerOpts{
		Limit:  d.Get("limit").(int),
		Marker: d.Get("marker").(string),
		ID:     d.Get("server_id").(string),
		Name:   d.Get("server_name").(string),
		Status: d.Get("status").(string),
		UserID: d.Get("user_id").(string),
	}
	pages, err := hosts.ListServer(dehClient, d.Get("id").(string), listServerOpts)

	if err != nil {
		return fmt.Errorf("Unable to retrieve deh server: %s", err)
	}

	if len(pages) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}
	if len(pages) > 1 {
		return fmt.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	DehServer := pages[0]

	log.Printf("[INFO] Retrieved Deh using given filter %s: %+v", DehServer.ID, DehServer)
	d.SetId(DehServer.ID)

	d.Set("server_id", DehServer.ID)
	d.Set("user_id", DehServer.UserID)
	d.Set("server_name", DehServer.Name)
	d.Set("status", DehServer.Status)
	d.Set("server_flavor", DehServer.Flavor)
	d.Set("metadata", DehServer.Metadata)
	d.Set("region", GetRegion(d, config))

	return nil

}
