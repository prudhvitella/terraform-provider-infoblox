package infoblox

/*
   This resource interface is basically a "helper" for the most common use case
   when using Infoblox -- you want to allocate an IP address from a particular
   network, and you want to get the next available IP address to use in creating
   an A record.  This resource will locate the network you want by CIDR (or other
   infoblox supported keys -- still specified as a "CIDR" in your terraform file), and
   then invoke NextAvailableIP against it, and return the result in a variable called
   "ipaddress".

   Note: this entire resource should probably be deprecated if someone
   implements a full Network resource (though the complexity of the
   API for such a resource might make it advisable to leave this
   around as a simple alternative for this common use case.


   Usage in Terraform file:


provider "infoblox" {
    username="whazzup"
    password="nuttin"
    host="https://infoblox.mydomain.com"
    sslverify="false"
    usecookies="false"
}

#this is the resource exposed by resource_infoblox_ip.go
#it will create a variable called "ipaddress"
resource "infoblox_ip" "theIPAddress" {
	cidr = "10.0.0.0/24"
}

#notice how the requested IP address is passed from the previous resource
#to this one through the "ipaddress" variable
resource "infoblox_record" "foobar" {
    value = "${infoblox_ip.theIPAddress.ipaddress}"
    name = "terraform"
    domain = "mydomain.com"
    type = "A"
    ttl = 3600
}


*/

import (
	"fmt"
	"github.com/fanatic/go-infoblox"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceInfobloxIP() *schema.Resource {
	return &schema.Resource{
		Create: resourceInfobloxIPCreate,
		Read:   resourceInfobloxIPRead,
		Update: resourceInfobloxIPUpdate,
		Delete: resourceInfobloxIPDelete,

		Schema: map[string]*schema.Schema{
			"cidr": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},

			"ip_range": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      false,
				ConflictsWith: []string{"cidr"},
			},

			"ipaddress": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Required: false,
			},

			"exclude": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
		},
	}
}

func getNextAvailableIPFromCIDR(client *infoblox.Client, cidr string, excludedAddresses []string) (string, error) {
	var (
		err    error
		result string
		ou     map[string]interface{}
	)

	network, err := getNetwork(client, cidr)

	if len(network) == 0 {
		err = fmt.Errorf("[ERROR] Empty response from client.Network().find. Is %s a valid network?", cidr)
	}

	if err == nil {
		ou, err = client.NetworkObject(network[0]["_ref"].(string)).NextAvailableIP(1, excludedAddresses)
		result = getMapValueAsString(ou, "ips")
		if result == "" {
			err = fmt.Errorf("Error: unable to determine IP address from response.\n")
		}
	}

	return result, err
}

func getNextAvailableIPFromRange(client *infoblox.Client, ip_range string, excludedAddresses []string) (string, error) {
	var (
		err    error
		result string
		ou     map[string]interface{}
	)

	s := "ip_address"
	q := []infoblox.Condition{
		infoblox.Condition{
			Field: &s,
			Value: getStartingIP(ip_range),
		},
	}

	out, err := client.Ipv4address().Find(q, nil)
	networkName := getNetworkNameFromIP(out, err)
	network, err := getNetwork(client, networkName)

	if err == nil {
		ou, err = client.NetworkObject(network[0]["_ref"].(string)).NextAvailableIP(1, excludedAddresses)
		result = getMapValueAsString(ou, "ips")
	}

	return result, err
}

func resourceInfobloxIPCreate(d *schema.ResourceData, meta interface{}) error {
	if err := validateIPData(d); err != nil {
		return err
	}

	var (
		result string
		err    error
	)

	client := meta.(*infoblox.Client)
	excludedAddresses := buildExcludedAddressesArray(d)

	if cidr, ok := d.GetOk("cidr"); ok {
		result, err = getNextAvailableIPFromCIDR(client, cidr.(string), excludedAddresses)
	} else if ip_range, ok := d.GetOk("ip_range"); ok {
		result, err = getNextAvailableIPFromRange(client, ip_range.(string), excludedAddresses)
	}

	if err != nil {
		return err
	}

	d.SetId(result)
	d.Set("ipaddress", result)

	return nil
}

// Validates that either 'cidr' or 'ip_range' is set
func validateIPData(d *schema.ResourceData) error {
	_, cidrOk := d.GetOk("cidr")
	_, ipRangeOk := d.GetOk("ip_range")
	if !cidrOk && !ipRangeOk {
		return fmt.Errorf(
			"One of ['cidr', 'ip_range'] must be set to create an Infoblox IP")
	}
	return nil
}

func resourceInfobloxIPRead(d *schema.ResourceData, meta interface{}) error {

	// since the infoblox network object's NextAvailableIP function isn't exactly
	// a resource (you don't really allocate an IP address until you use the record:a or
	// record:host object), we don't actually implement READ, UPDATE, or DELETE

	return nil
}

func resourceInfobloxIPUpdate(d *schema.ResourceData, meta interface{}) error {

	// since the infoblox network object's NextAvailableIP function isn't exactly
	// a resource (you don't really allocate an IP address until you use the record:a or
	// record:host object), we don't actually implement READ, UPDATE, or DELETE

	return nil
}

func resourceInfobloxIPDelete(d *schema.ResourceData, meta interface{}) error {

	// since the infoblox network object's NextAvailableIP function isn't exactly
	// a resource (you don't really allocate an IP address until you use the record:a or
	// record:host object), we don't actually implement READ, UPDATE, or DELETE

	return nil
}
