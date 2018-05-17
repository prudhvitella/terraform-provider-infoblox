package infoblox

import (
	"fmt"
	"log"
	"net/url"

	infoblox "github.com/defilan/go-infoblox"
	"github.com/hashicorp/terraform/helper/schema"
)

// hostIPv4Schema represents the schema for the host IPv4 sub-resource
func hostIPv4Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"address": {
			Type:     schema.TypeString,
			Required: true,
		},
		"configure_for_dhcp": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"host": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"mac": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
}

// hostIPv6Schema represents the schema for the host IPv4 sub-resource
func hostIPv6Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"address": {
			Type:     schema.TypeString,
			Required: true,
		},
		"configure_for_dhcp": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"host": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"mac": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
}

func infobloxRecordHost() *schema.Resource {
	return &schema.Resource{
		Create: resourceInfobloxHostRecordCreate,
		Read:   resourceInfobloxHostRecordRead,
		Update: resourceInfobloxHostRecordUpdate,
		Delete: resourceInfobloxHostRecordDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ipv4addr": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Resource{Schema: hostIPv4Schema()},
			},
			"ipv6addr": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Resource{Schema: hostIPv6Schema()},
			},
			"configure_for_dns": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"comment": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"ttl": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"view": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func ipv4sFromList(ipv4s []interface{}) []infoblox.HostIpv4Addr {
	result := make([]infoblox.HostIpv4Addr, 0, len(ipv4s))
	for _, v := range ipv4s {
		ipMap := v.(map[string]interface{})
		i := infoblox.HostIpv4Addr{}

		i.Ipv4Addr = ipMap["address"].(string)

		if val, ok := ipMap["configure_for_dhcp"]; ok {
			i.ConfigureForDHCP = val.(bool)
		}
		if val, ok := ipMap["host"]; ok {
			i.Host = val.(string)
		}
		if val, ok := ipMap["mac"]; ok {
			i.MAC = val.(string)
		}

		result = append(result, i)
	}
	return result
}

func ipv6sFromList(ipv6s []interface{}) []infoblox.HostIpv6Addr {
	result := make([]infoblox.HostIpv6Addr, 0, len(ipv6s))
	for _, v := range ipv6s {
		ip := v.(*schema.ResourceData)
		i := infoblox.HostIpv6Addr{}

		if attr, ok := ip.GetOk("address"); ok {
			i.Ipv6Addr = attr.(string)
		}
		if attr, ok := ip.GetOk("configure_for_dhcp"); ok {
			i.ConfigureForDHCP = attr.(bool)
		}
		if attr, ok := ip.GetOk("host"); ok {
			i.Host = attr.(string)
		}
		if attr, ok := ip.GetOk("mac"); ok {
			i.MAC = attr.(string)
		}
		result = append(result, i)
	}
	return result
}

func hostObjectFromAttributes(d *schema.ResourceData) infoblox.RecordHostObject {
	hostObject := infoblox.RecordHostObject{}
	if attr, ok := d.GetOk("name"); ok {
		hostObject.Name = attr.(string)
	}
	if attr, ok := d.GetOk("configure_for_dns"); ok {
		log.Printf("[DEBUG] FOUND CONFIGURE_FOR_DNS")
		hostObject.ConfigureForDNS = attr.(bool)
	}
	if attr, ok := d.GetOk("comment"); ok {
		hostObject.Comment = attr.(string)
	}
	if attr, ok := d.GetOk("ttl"); ok {
		hostObject.Ttl = attr.(int)
	}
	if attr, ok := d.GetOk("view"); ok {
		hostObject.View = attr.(string)
	}
	if attr, ok := d.GetOk("ipv4addr"); ok {
		hostObject.Ipv4Addrs = ipv4sFromList(attr.([]interface{}))
	}
	if attr, ok := d.GetOk("ipv6addr"); ok {
		hostObject.Ipv6Addrs = ipv6sFromList(attr.([]interface{}))
	}

	return hostObject
}

func resourceInfobloxHostRecordCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*infoblox.Client)

	record := url.Values{}
	hostObject := hostObjectFromAttributes(d)

	log.Printf("[DEBUG] Creating Infoblox Host record with configuration: %#v", hostObject)
	opts := &infoblox.Options{
		ReturnFields: []string{"name", "ipv4addr", "ipv6addr", "configure_for_dns", "comment", "ttl", "view"},
	}
	recordID, err := client.RecordHost().Create(record, opts, hostObject)
	if err != nil {
		return fmt.Errorf("error creating infoblox Host record: %s", err.Error())
	}

	d.SetId(recordID)
	log.Printf("[INFO] Infoblox Host record created with ID: %s", d.Id())
	return nil
}

func resourceInfobloxHostRecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*infoblox.Client)

	record, err := client.GetRecordHost(d.Id(), nil)
	if err != nil {
		return handleReadError(d, "Host", err)
	}

	d.Set("name", record.Name)
	if &record.ConfigureForDNS != nil {
		d.Set("configure_for_dns", record.ConfigureForDNS)
	}
	if &record.Comment != nil {
		d.Set("comment", record.Comment)
	}
	if &record.Ttl != nil {
		d.Set("ttl", record.Ttl)
	}
	if &record.View != nil {
		d.Set("view", record.View)
	}
	if &record.Ipv4Addrs != nil {
		result := make([]map[string]interface{}, 1)
		for _, v := range record.Ipv4Addrs {
			i := make(map[string]interface{})

			i["address"] = v.Ipv4Addr
			if &v.ConfigureForDHCP != nil {
				i["configure_for_dhcp"] = v.ConfigureForDHCP
			}
			if &v.Host != nil {
				i["host"] = v.Host
			}
			if &v.MAC != nil {
				i["mac"] = v.MAC
			}

			result = append(result, i)
		}
		d.Set("ipv4addr", result)
	}
	if &record.Ipv6Addrs != nil {
		result := make([]map[string]interface{}, 1)
		for _, v := range record.Ipv6Addrs {
			i := make(map[string]interface{})

			i["address"] = v.Ipv6Addr
			if &v.ConfigureForDHCP != nil {
				i["configure_for_dhcp"] = v.ConfigureForDHCP
			}
			if &v.Host != nil {
				i["host"] = v.Host
			}
			if &v.MAC != nil {
				i["mac"] = v.MAC
			}

			result = append(result, i)
		}
		d.Set("ipv6addr", result)
	}

	return nil
}

func resourceInfobloxHostRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*infoblox.Client)

	_, err := client.GetRecordHost(d.Id(), nil)
	if err != nil {
		return fmt.Errorf("error finding infoblox Host record: %s", err.Error())
	}

	record := url.Values{}
	hostObject := hostObjectFromAttributes(d)

	log.Printf("[DEBUG] Updating Infoblox Host record with configuration: %#v", hostObject)

	opts := &infoblox.Options{
		ReturnFields: []string{"name", "ipv4addr", "ipv6addr", "configure_for_dns", "comment", "ttl", "view"},
	}
	recordID, err := client.RecordHostObject(d.Id()).Update(record, opts, hostObject)
	if err != nil {
		return fmt.Errorf("error updating Infoblox Host record: %s", err.Error())
	}

	d.SetId(recordID)
	log.Printf("[INFO] Infoblox Host record updated with ID: %s", d.Id())

	return resourceInfobloxHostRecordRead(d, meta)
}

func resourceInfobloxHostRecordDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*infoblox.Client)

	log.Printf("[DEBUG] Deleting Infoblox Host record: %s, %s", d.Get("name").(string), d.Id())
	_, err := client.GetRecordHost(d.Id(), nil)
	if err != nil {
		return fmt.Errorf("error finding Infoblox Host record: %s", err.Error())
	}

	err = client.RecordHostObject(d.Id()).Delete(nil)
	if err != nil {
		return fmt.Errorf("error deleting Infoblox Host record: %s", err.Error())
	}

	return nil
}
