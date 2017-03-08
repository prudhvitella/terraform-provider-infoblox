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
	"log"
	"strings"
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
				Required: true,
				ForceNew: false,
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

func resourceInfobloxIPCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*infoblox.Client)

	log.Print("[TRACE] inside resourceInfobloxIPCreate.")

	ntwork := d.Get("cidr")

	log.Printf("[TRACE] CIDR from terraform file: %s", ntwork.(string))

	s := "network"
	q := []infoblox.Condition{
		infoblox.Condition{
			Field: &s,
			Value: ntwork.(string),
		},
	}

	log.Print("[TRACE] invoking client.Network().find")

	out, err := client.Network().Find(q, nil)

	if err != nil {
		log.Printf("[ERROR] Unable to invoke find on cidr: %s, %s", ntwork, err)
		return err
	}

	if len(out) == 0 {
		return fmt.Errorf("Empty response from client.Network().find. Is %s a valid network?", ntwork)
	}

	printList(out, nil)

	log.Print("[TRACE] invoking client.NetworkObject().NextAvailableIP")

	var excludedAddresses []string
	if userExcludes := d.Get("exclude"); userExcludes != nil {
		addresses := userExcludes.(*schema.Set).List()
		for _, address := range addresses {
			excludedAddresses = append(excludedAddresses, address.(string))
		}
	}

	log.Printf("[TRACE] Excluding Addresses = %v", excludedAddresses)

	ou, err := client.NetworkObject(out[0]["_ref"].(string)).NextAvailableIP(1, excludedAddresses)

	if err != nil {
		log.Printf("[ERROR] Unable to allocate NextAvailableIP: %s", err)
		return err
	}

	printObject(ou, nil)

	log.Print("[TRACE] Walking NextAvailableIP output to get ip")

	res := getMapValueAsString(ou, "ips")

	if res == "" {
		log.Print("Error: unable to determine IP address from response \n", err)
		return nil
	}

	log.Printf("[TRACE] returned value in ips structure: %s", res)

	log.Print("[TRACE] Setting ID, locking provisioned IP in terraform")

	d.SetId(res)

	log.Print("[TRACE] Setting output variable 'ipaddress'")

	d.Set("ipaddress", res)

	log.Print("[TRACE] exiting resourceInfobloxIPCreate.")

	return nil
}

// TODO: I'm positive there's a better way to do this, but this works for now
func getMapValueAsString(mymap map[string]interface{}, val string) string {
	for k, v := range mymap {
		if k == val {
			vout := fmt.Sprintf("%q", v)
			vout = strings.Replace(vout, "[", "", -1)
			vout = strings.Replace(vout, "]", "", -1)
			vout = strings.Replace(vout, "\"", "", -1)
			return vout
		}
	}

	return ""
}

func printList(out []map[string]interface{}, err error) {
	e(err)
	for i, v := range out {
		log.Printf("[%d]\n", i)
		printObject(v, nil)
	}
}

func printObject(out map[string]interface{}, err error) {
	e(err)
	for k, v := range out {
		log.Printf("  %s: %q\n", k, v)
	}
	log.Printf("\n")
}

func e(err error) {
	if err != nil {
		log.Printf("Error: %v\n", err)
	}
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
