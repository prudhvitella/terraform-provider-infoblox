package infoblox

import (
	"fmt"
	"github.com/fanatic/go-infoblox"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"strings"
)

func getNetworkNameFromIP(response []map[string]interface{}, err error) string {
	e(err)

	for _, v := range response {
		for k, val := range v {
			if k == "network" {
				return val.(string)
			}
		}
	}

	return ""
}

func getNetwork(client *infoblox.Client, term string) ([]map[string]interface{}, error) {
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

func getStartingIP(ip_range string) string {
	return strings.Split(ip_range, "-")[0]
}

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
