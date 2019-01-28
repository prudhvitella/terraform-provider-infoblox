package infoblox

import (
	"fmt"
	"log"
	"net/url"

	infoblox "github.com/fanatic/go-infoblox"
	"github.com/hashicorp/terraform/helper/schema"
)

func infobloxRecordAAAA() *schema.Resource {
	return &schema.Resource{
		Create: resourceInfobloxAAAARecordCreate,
		Read:   resourceInfobloxAAAARecordRead,
		Update: resourceInfobloxAAAARecordUpdate,
		Delete: resourceInfobloxAAAARecordDelete,

		Schema: map[string]*schema.Schema{
			// TODO: validate that address is in IPv6 format.
			"address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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

func resourceInfobloxAAAARecordCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*infoblox.Client)

	record := url.Values{}
	record.Add("ipv6addr", d.Get("address").(string))
	record.Add("name", d.Get("name").(string))
	populateSharedAttributes(d, &record)

	log.Printf("[DEBUG] Creating Infoblox AAAA record with configuration: %#v", record)

	opts := &infoblox.Options{
		ReturnFields: []string{"ipv6addr", "name", "comment", "ttl", "view"},
	}
	recordID, err := client.RecordAAAA().Create(record, opts, nil)

	if err != nil {
		return fmt.Errorf("error creating infoblox AAAA record: %s", err.Error())
	}

	d.SetId(recordID)
	log.Printf("[INFO] Infoblox AAAA record created with ID: %s", d.Id())

	return nil
}

func resourceInfobloxAAAARecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*infoblox.Client)

	opts := &infoblox.Options{
		ReturnFields: []string{"ipv6addr", "name", "comment", "ttl", "view"},
	}
	record, err := client.GetRecordAAAA(d.Id(), opts)
	if err != nil {
		return handleReadError(d, "AAAA", err)
	}

	d.Set("address", record.Ipv6Addr)
	d.Set("name", record.Name)

	if &record.Comment != nil {
		d.Set("comment", record.Comment)
	}
	if &record.Ttl != nil {
		d.Set("ttl", record.Ttl)
	}
	if &record.View != nil {
		d.Set("view", record.View)
	}

	return nil
}

func resourceInfobloxAAAARecordUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*infoblox.Client)

	opts := &infoblox.Options{
		ReturnFields: []string{"ipv6addr", "name", "comment", "ttl", "view"},
	}
	_, err := client.GetRecordAAAA(d.Id(), opts)
	if err != nil {
		return fmt.Errorf("error finding infoblox AAAA record: %s", err.Error())
	}

	record := url.Values{}
	record.Add("ipv6addr", d.Get("address").(string))
	record.Add("name", d.Get("name").(string))
	populateSharedAttributes(d, &record)

	log.Printf("[DEBUG] Updating Infoblox AAAA record with configuration: %#v", record)

	recordID, err := client.RecordAAAAObject(d.Id()).Update(record, opts, nil)
	if err != nil {
		return fmt.Errorf("error updating Infoblox AAAA record: %s", err.Error())
	}

	d.SetId(recordID)
	log.Printf("[INFO] Infoblox AAAA record updated with ID: %s", d.Id())

	return resourceInfobloxAAAARecordRead(d, meta)
}

func resourceInfobloxAAAARecordDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*infoblox.Client)

	log.Printf("[DEBUG] Deleting Infoblox AAAA record: %s, %s", d.Get("name").(string), d.Id())
	_, err := client.GetRecordAAAA(d.Id(), nil)
	if err != nil {
		return fmt.Errorf("error finding Infoblox AAAA record: %s", err.Error())
	}

	err = client.RecordAAAAObject(d.Id()).Delete(nil)
	if err != nil {
		return fmt.Errorf("error deleting Infoblox AAAA record: %s", err.Error())
	}

	return nil
}
