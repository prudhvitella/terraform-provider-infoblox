package infoblox

import (
	"fmt"
	"log"
	"net/url"

	infoblox "github.com/fanatic/go-infoblox"
	"github.com/hashicorp/terraform/helper/schema"
)

func infobloxRecordA() *schema.Resource {
	return &schema.Resource{
		Create: resourceInfobloxARecordCreate,
		Read:   resourceInfobloxARecordRead,
		Update: resourceInfobloxARecordUpdate,
		Delete: resourceInfobloxARecordDelete,

		Schema: map[string]*schema.Schema{
			// TODO: validate that address is in IPv4 format.
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

func resourceInfobloxARecordCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*infoblox.Client)

	record := url.Values{}
	record.Add("ipv4addr", d.Get("address").(string))
	record.Add("name", d.Get("name").(string))
	populateSharedAttributes(d, &record)

	log.Printf("[DEBUG] Creating Infoblox A record with configuration: %#v", record)

	opts := &infoblox.Options{
		ReturnFields: []string{"ipv4addr", "name", "comment", "ttl", "view"},
	}
	recordID, err := client.RecordA().Create(record, opts, nil)

	if err != nil {
		return fmt.Errorf("error creating infoblox A record: %s", err.Error())
	}

	d.SetId(recordID)
	log.Printf("[INFO] Infoblox A record created with ID: %s", d.Id())

	return nil
}

func resourceInfobloxARecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*infoblox.Client)

	record, err := client.GetRecordA(d.Id(), nil)
	if err != nil {
		return handleReadError(d, "A", err)
	}

	d.Set("address", record.Ipv4Addr)
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

func resourceInfobloxARecordUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*infoblox.Client)

	_, err := client.GetRecordA(d.Id(), nil)
	if err != nil {
		return fmt.Errorf("error finding infoblox A record: %s", err.Error())
	}

	record := url.Values{}
	record.Add("ipv4addr", d.Get("address").(string))
	record.Add("name", d.Get("name").(string))
	populateSharedAttributes(d, &record)

	log.Printf("[DEBUG] Updating Infoblox A record with configuration: %#v", record)

	opts := &infoblox.Options{
		ReturnFields: []string{"ipv4addr", "name", "comment", "ttl", "view"},
	}
	recordID, err := client.RecordAObject(d.Id()).Update(record, opts, nil)
	if err != nil {
		return fmt.Errorf("error updating Infoblox A record: %s", err.Error())
	}

	d.SetId(recordID)
	log.Printf("[INFO] Infoblox A record updated with ID: %s", d.Id())

	return resourceInfobloxARecordRead(d, meta)
}

func resourceInfobloxARecordDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*infoblox.Client)

	log.Printf("[DEBUG] Deleting Infoblox A record: %s, %s", d.Get("name").(string), d.Id())
	_, err := client.GetRecordA(d.Id(), nil)
	if err != nil {
		return fmt.Errorf("error finding Infoblox A record: %s", err.Error())
	}

	err = client.RecordAObject(d.Id()).Delete(nil)
	if err != nil {
		return fmt.Errorf("error deleting Infoblox A record: %s", err.Error())
	}

	return nil
}
