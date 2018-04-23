package infoblox

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"

	"github.com/fanatic/go-infoblox"
	"github.com/hashicorp/terraform/helper/schema"
)

// The comment, ttl, and, view attributes are common across all of the DNS
// record objects we deal with so far so we extract populating them into the
// url.Values object into a helper function.
func populateSharedAttributes(d *schema.ResourceData, record *url.Values) {
	if attr, ok := d.GetOk("comment"); ok {
		record.Set("comment", attr.(string))
	}

	if attr, ok := d.GetOk("ttl"); ok {
		record.Set("ttl", strconv.Itoa(attr.(int)))
	}

	if attr, ok := d.GetOk("view"); ok {
		record.Set("view", attr.(string))
	}
}

// Parses the given string as an ip address and returns "ipv4addr" if it is an
// ipv4 address and "ipv6addr" if it is an ipv6 address
func ipType(value string) (string, error) {
	ip := net.ParseIP(value)
	if ip == nil {
		return "", fmt.Errorf("value does not appear to be a valid ip address")
	}

	res := "ipv6addr"
	if ip.To4() != nil {
		res = "ipv4addr"
	}
	return res, nil
}

// Finds networks by search term, such as network CIDR.
func getNetworks(client *infoblox.Client, term string) ([]map[string]interface{}, error) {
	s := "network"
	q := []infoblox.Condition{
		infoblox.Condition{
			Field: &s,
			Value: term,
		},
	}

	network, err := client.Network().Find(q, nil)
	return network, err
}

// Builds an array of IP addresses to exclude from terraform resource data.
func buildExcludedAddressesArray(d *schema.ResourceData) []string {
	var excludedAddresses []string
	if userExcludes := d.Get("exclude"); userExcludes != nil {
		addresses := userExcludes.(*schema.Set).List()
		for _, address := range addresses {
			excludedAddresses = append(excludedAddresses, address.(string))
		}
	}
	return excludedAddresses
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

// Validates that either 'cidr' or 'ip_range' terraform argument is set.
func validateIPData(d *schema.ResourceData) error {
	_, cidrOk := d.GetOk("cidr")
	_, ipRangeOk := d.GetOk("ip_range")
	if !cidrOk && !ipRangeOk {
		return fmt.Errorf(
			"One of ['cidr', 'ip_range'] must be set to create an Infoblox IP")
	}
	return nil
}

func handleReadError(d *schema.ResourceData, recordType string, err error) error {
	if infobloxErr, ok := err.(infoblox.Error); ok {
		if infobloxErr.Code() == "Client.Ibap.Data.NotFound" {
			d.SetId("")
			return nil
		}
	}
	return fmt.Errorf("Error reading Infoblox %s record: %s", recordType, err)
}
