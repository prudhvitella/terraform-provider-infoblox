package infoblox

import (
	"fmt"
	"log"
	"net/url"
	"strconv"

	infoblox "github.com/fanatic/go-infoblox"
	"github.com/hashicorp/terraform/helper/schema"
)

func infobloxRecordMX() *schema.Resource {
	return &schema.Resource{
		Create: resourceInfobloxMXRecordCreate,
		Read:   resourceInfobloxMXRecordRead,
		Update: resourceInfobloxMXRecordUpdate,
		Delete: resourceInfobloxMXRecordDelete,

		Schema: map[string]*schema.Schema{
			"exchanger": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"pref": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: false,
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

func resourceInfobloxMXRecordCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*infoblox.Client)

	record := url.Values{}
	record.Add("exchanger", d.Get("exchanger").(string))
	record.Add("name", d.Get("name").(string))
	record.Add("pref", strconv.Itoa(d.Get("pref").(int)))
	populateSharedAttributes(d, &record)

	log.Printf("[DEBUG] Creating Infoblox MX record with configuration: %#v", record)

	opts := &infoblox.Options{
		ReturnFields: []string{"exchanger", "name", "pref", "comment", "ttl", "view"},
	}

	// TODO: Add MX support to go-infoblox
	recordID, err := client.RecordMx().Create(record, opts, nil)

	if err != nil {
		return fmt.Errorf("error creating infoblox MX record: %s", err.Error())
	}

	d.SetId(recordID)
	log.Printf("[INFO] Infoblox MX record created with ID: %s", d.Id())

	return nil
}

func resourceInfobloxMXRecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*infoblox.Client)

	opts := &infoblox.Options{
		ReturnFields: []string{"exchanger", "name", "pref", "comment", "ttl", "view"},
	}
	record, err := client.GetRecordMx(d.Id(), opts)
	if err != nil {
		return handleReadError(d, "MX", err)
	}

	d.Set("exchanger", record.Exchanger)
	d.Set("name", record.Name)
	d.Set("pref", record.Pref)

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

func resourceInfobloxMXRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*infoblox.Client)

	opts := &infoblox.Options{
		ReturnFields: []string{"exchanger", "name", "pref", "comment", "ttl", "view"},
	}
	_, err := client.GetRecordMx(d.Id(), opts)
	if err != nil {
		return fmt.Errorf("error finding infoblox MX record: %s", err.Error())
	}

	record := url.Values{}
	record.Add("exchanger", d.Get("address").(string))
	record.Add("name", d.Get("name").(string))
	record.Add("pref", strconv.Itoa(d.Get("pref").(int)))
	populateSharedAttributes(d, &record)

	log.Printf("[DEBUG] Updating Infoblox MX record with configuration: %#v", record)

	recordID, err := client.RecordMxObject(d.Id()).Update(record, opts, nil)
	if err != nil {
		return fmt.Errorf("error updating Infoblox MX record: %s", err.Error())
	}

	d.SetId(recordID)
	log.Printf("[INFO] Infoblox MX record updated with ID: %s", d.Id())

	return resourceInfobloxMXRecordRead(d, meta)
}

func resourceInfobloxMXRecordDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*infoblox.Client)

	log.Printf("[DEBUG] Deleting Infoblox MX record: %s, %s", d.Get("name").(string), d.Id())
	_, err := client.GetRecordMx(d.Id(), nil)
	if err != nil {
		return fmt.Errorf("error finding Infoblox MX record: %s", err.Error())
	}

	err = client.RecordMxObject(d.Id()).Delete(nil)
	if err != nil {
		return fmt.Errorf("error deleting Infoblox MX record: %s", err.Error())
	}

	return nil
}
