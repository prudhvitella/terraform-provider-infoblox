package infoblox

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"

	infoblox "github.com/fanatic/go-infoblox"
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

func nextAvailableIP(client *infoblox.Client, cidr string) (string, error) {
	var result string
	var err error
	var ou map[string]interface{}

	network, err := getNetworks(client, cidr)

	if err != nil {
		if strings.Contains(err.Error(), "Authorization Required") {
			return "", fmt.Errorf("[ERROR] Authentication Error, Please check your username/password ")
		}
	}

	if len(network) == 0 {
		err = fmt.Errorf("[ERROR] Empty response from client.Network().find. Is %s a valid network?", cidr)
	}

	if err == nil {
		ou, err = client.NetworkObject(network[0]["_ref"].(string)).NextAvailableIP(1, nil)
		result = getMapValueAsString(ou, "ips")
		if result == "" {
			err = fmt.Errorf("[ERROR] Unable to determine IP address from response")
		}
	}

	return result, err
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
