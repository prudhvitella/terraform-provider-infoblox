package infoblox

import (
	"fmt"
	"log"
	"net/url"

	infoblox "github.com/fanatic/go-infoblox"
	"github.com/hashicorp/terraform/helper/schema"
)

func infobloxRecordTXT() *schema.Resource {
	return &schema.Resource{
		Create: resourceInfobloxTXTRecordCreate,
		Read:   resourceInfobloxTXTRecordRead,
		Update: resourceInfobloxTXTRecordUpdate,
		Delete: resourceInfobloxTXTRecordDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"text": &schema.Schema{
				Type:     schema.TypeString,
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

func resourceInfobloxTXTRecordCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*infoblox.Client)

	record := url.Values{}
	record.Add("name", d.Get("name").(string))
	record.Add("text", d.Get("text").(string))
	populateSharedAttributes(d, &record)

	log.Printf("[DEBUG] Creating Infoblox TXT record with configuration: %#v", record)

	opts := &infoblox.Options{
		ReturnFields: []string{"name", "text", "comment", "ttl", "view"},
	}

	recordID, err := client.RecordTxt().Create(record, opts, nil)

	if err != nil {
		return fmt.Errorf("error creating infoblox TXT record: %s", err.Error())
	}

	d.SetId(recordID)
	log.Printf("[INFO] Infoblox TXT record created with ID: %s", d.Id())

	return nil
}

func resourceInfobloxTXTRecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*infoblox.Client)

	opts := &infoblox.Options{
		ReturnFields: []string{"name", "text", "comment", "ttl", "view"},
	}
	record, err := client.GetRecordTxt(d.Id(), opts)
	if err != nil {
		return handleReadError(d, "TXT", err)
	}

	d.Set("name", record.Name)
	d.Set("text", record.Text)

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

func resourceInfobloxTXTRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*infoblox.Client)

	opts := &infoblox.Options{
		ReturnFields: []string{"name", "text", "comment", "ttl", "view"},
	}
	_, err := client.GetRecordTxt(d.Id(), opts)
	if err != nil {
		return fmt.Errorf("error finding infoblox TXT record: %s", err.Error())
	}

	record := url.Values{}
	record.Add("name", d.Get("name").(string))
	record.Add("text", d.Get("text").(string))
	populateSharedAttributes(d, &record)

	log.Printf("[DEBUG] Updating Infoblox TXT record with configuration: %#v", record)

	recordID, err := client.RecordTxtObject(d.Id()).Update(record, opts, nil)
	if err != nil {
		return fmt.Errorf("error updating Infoblox TXT record: %s", err.Error())
	}

	d.SetId(recordID)
	log.Printf("[INFO] Infoblox TXT record updated with ID: %s", d.Id())

	return resourceInfobloxTXTRecordRead(d, meta)
}

func resourceInfobloxTXTRecordDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*infoblox.Client)

	log.Printf("[DEBUG] Deleting Infoblox TXT record: %s, %s", d.Get("name").(string), d.Id())
	_, err := client.GetRecordTxt(d.Id(), nil)
	if err != nil {
		return fmt.Errorf("error finding Infoblox TXT record: %s", err.Error())
	}

	err = client.RecordTxtObject(d.Id()).Delete(nil)
	if err != nil {
		return fmt.Errorf("error deleting Infoblox TXT record: %s", err.Error())
	}

	return nil
}
