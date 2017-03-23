package infoblox

import (
	"fmt"
	"github.com/fanatic/go-infoblox"
	"github.com/hashicorp/terraform/helper/schema"
	"strings"
)

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
