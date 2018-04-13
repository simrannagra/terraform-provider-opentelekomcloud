package opentelekomcloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/networking/v1/subnets"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/huaweicloud/golangsdk"
)

func resourceSubnetDNSListV1(d *schema.ResourceData) []string {
	rawDNSN := d.Get("dns_list").(*schema.Set)
	dnsn := make([]string, rawDNSN.Len())
	for i, raw := range rawDNSN.List() {
		dnsn[i] = raw.(string)
	}
	return dnsn
}

func resourceVpcSubnetV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceVpcSubnetV1Create, //providers.go
		Read:   resourceVpcSubnetV1Read,
		Update: resourceVpcSubnetV1Update,
		Delete: resourceVpcSubnetV1Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{ //request and response parameters
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     false,
				ValidateFunc: validateName,
			},
			"cidr": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateCIDR,
			},
			"dns_list": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Required: false,
				Elem:     &schema.Schema{Type: schema.TypeString, ValidateFunc: validateIP},
				Set:      schema.HashString,
				Computed: true,
			},
			"gateway_ip": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateIP,
			},
			"dhcp_enable": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: false,
			},
			"primary_dns": &schema.Schema{
				Type:         schema.TypeString,
				ForceNew:     false,
				Optional:     true,
				ValidateFunc: validateIP,
				Computed:     true,
			},
			"secondary_dns": &schema.Schema{
				Type:         schema.TypeString,
				ForceNew:     false,
				Optional:     true,
				ValidateFunc: validateIP,
				Computed:     true,
			},
			"availability_zone": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},
			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
		},
	}
}

func resourceVpcSubnetV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	subnetClient, err := config.networkingV1Client(GetRegion(d, config))

	log.Printf("[DEBUG] Value of networking Client: %#v", subnetClient)

	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
	}

	createOpts := subnets.CreateOpts{
		Name:             d.Get("name").(string),
		CIDR:             d.Get("cidr").(string),
		AvailabilityZone: d.Get("availability_zone").(string),
		GatewayIP:        d.Get("gateway_ip").(string),
		EnableDHCP:       d.Get("dhcp_enable").(bool),
		VPC_ID:           d.Get("vpc_id").(string),
		PRIMARY_DNS:      d.Get("primary_dns").(string),
		SECONDARY_DNS:    d.Get("secondary_dns").(string),
		DnsList:          resourceSubnetDNSListV1(d),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	n, err := subnets.Create(subnetClient, createOpts).Extract()
	log.Printf("[DEBUG] Create Subnet: %#v", n)
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud VPC subnet: %s", err)
	}
	log.Printf("[INFO] Vpc Subnet ID: %s", n.ID)

	log.Printf("[DEBUG] Waiting for OpenTelekomCloud Vpc Subnet(%s) to become available", n.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"CREATING"},
		Target:     []string{"ACTIVE"},
		Refresh:    waitForVpcSubnetActive(subnetClient, n.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	d.SetId(n.ID)

	return resourceVpcSubnetV1Read(d, config)

}

func resourceVpcSubnetV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	subnetClient, err := config.networkingV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
	}

	n, err := subnets.Get(subnetClient, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving OpenTelekomCloud Subnets: %s", err)
	}

	log.Printf("[DEBUG] Retrieved subnet %s: %+v", d.Id(), n)

	d.Set("name", n.Name)
	d.Set("cidr", n.CIDR)
	d.Set("dns_list", n.DnsList)
	d.Set("gateway_ip", n.GatewayIP)
	d.Set("dhcp_enable", n.EnableDHCP)
	d.Set("primary_dns", n.PRIMARY_DNS)
	d.Set("secondary_dns", n.SECONDARY_DNS)
	d.Set("availability_zone", n.AvailabilityZone)
	d.Set("vpc_id", n.VPC_ID)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceVpcSubnetV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	subnetClient, err := config.networkingV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
	}

	var update bool
	var updateOpts subnets.UpdateOpts

	//as name is mandatory while updating subnet
	updateOpts.Name = d.Get("name").(string)

	if d.HasChange("name") {
		update = true
	}
	if d.HasChange("primary_dns") {
		update = true
		updateOpts.PRIMARY_DNS = d.Get("primary_dns").(string)
	}
	if d.HasChange("secondary_dns") {
		update = true
		updateOpts.SECONDARY_DNS = d.Get("secondary_dns").(string)
	}
	if d.HasChange("dns_list") {
		update = true
		updateOpts.DnsList = resourceSubnetDNSListV1(d)
	}
	if d.HasChange("dhcp_enable") {
		update = true
		updateOpts.EnableDHCP = d.Get("dhcp_enable").(bool)

	} else if update { //maintaining dhcp to be true if it was true earlier as default update option for dhcp bool is always going to be false in golangsdk
		if d.Get("dhcp_enable").(bool) {
			updateOpts.EnableDHCP = true
		}
	}

	vpc_id := d.Get("vpc_id").(string)

	log.Printf("[DEBUG] Subnet_id %+v", d.Id())

	if update {
		log.Printf("[DEBUG] Updating subnet %s with options: %#v", d.Id(), updateOpts)
		_, err = subnets.Update(subnetClient, vpc_id, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error updating OpenTelekomCloud VPC Subnet: %s", err)
		}
	}
	return resourceVpcSubnetV1Read(d, meta)
}

func resourceVpcSubnetV1Delete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Destroy subnet: %s", d.Id())

	config := meta.(*Config)
	subnetClient, err := config.networkingV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
	}
	vpc_id := d.Get("vpc_id").(string)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForVpcSubnetDelete(subnetClient, vpc_id, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error deleting OpenTelekomCloud Subnet: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForVpcSubnetActive(subnetClient *golangsdk.ServiceClient, vpcId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := subnets.Get(subnetClient, vpcId).Extract()
		if err != nil {
			return nil, "", err
		}

		log.Printf("[DEBUG] OpenTelekomCloud Subnet Client: %+v", n)
		if n.Status == "DOWN" || n.Status == "ACTIVE" || n.Status == "ERROR" {
			return n, "ACTIVE", nil
		}

		if n.Status == "UNKNOWN" {
			return nil, "", fmt.Errorf("The CIDR of the created subnet is the same as that of an existing subnet.")
		}

		return n, "CREATING", nil
	}
}

func waitForVpcSubnetDelete(subnetClient *golangsdk.ServiceClient, vpcId string, subnetId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to delete OpenTelekomCloud subnet %s.\n", subnetId)

		r, err := subnets.Get(subnetClient, subnetId).Extract()
		log.Printf("[DEBUG] Value after extract: %#v", r)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted OpenTelekomCloud subnet %s", subnetId)
				return r, "DELETED", nil
			}
			return r, "ACTIVE", err
		}
		err = subnets.Delete(subnetClient, vpcId, subnetId).ExtractErr()
		log.Printf("[DEBUG] Value if error: %#v", err)

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted OpenTelekomCloud subnet %s", subnetId)
				return r, "DELETED", nil
			}
			if errCode, ok := err.(golangsdk.ErrUnexpectedResponseCode); ok {
				if errCode.Actual == 409 {
					return r, "ACTIVE", nil
				}
			}
			return r, "ACTIVE", err
		}

		log.Printf("[DEBUG] OpenTelekomCloud subnet %s still active.\n", subnetId)
		return r, "ACTIVE", nil
	}
}