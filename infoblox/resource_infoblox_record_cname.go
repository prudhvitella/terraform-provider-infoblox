package infoblox

import (
	"fmt"
	"log"
	"net/url"

	infoblox "github.com/fanatic/go-infoblox"
	"github.com/hashicorp/terraform/helper/schema"
)

func infobloxRecordCNAME() *schema.Resource {
	return &schema.Resource{
		Create: resourceInfobloxCNAMERecordCreate,
		Read:   resourceInfobloxCNAMERecordRead,
		Update: resourceInfobloxCNAMERecordUpdate,
		Delete: resourceInfobloxCNAMERecordDelete,

		Schema: map[string]*schema.Schema{
			"canonical": &schema.Schema{
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

func resourceInfobloxCNAMERecordCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*infoblox.Client)

	record := url.Values{}
	record.Add("canonical", d.Get("canonical").(string))
	record.Add("name", d.Get("name").(string))
	populateSharedAttributes(d, &record)

	log.Printf("[DEBUG] Creating Infoblox CNAME record with configuration: %#v", record)

	opts := &infoblox.Options{
		ReturnFields: []string{"canonical", "name", "comment", "ttl", "view"},
	}
	recordID, err := client.RecordCname().Create(record, opts, nil)
	if err != nil {
		return fmt.Errorf("error creating infoblox CNAME record: %s", err.Error())
	}

	d.SetId(recordID)
	log.Printf("[INFO] Infoblox CNAME record created with ID: %s", d.Id())

	return nil
}

func resourceInfobloxCNAMERecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*infoblox.Client)

	opts := &infoblox.Options{
		ReturnFields: []string{"canonical", "name", "comment", "ttl", "view"},
	}
	record, err := client.GetRecordCname(d.Id(), opts)
	if err != nil {
		return handleReadError(d, "CNAME", err)
	}

	d.Set("canonical", record.Canonical)
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

func resourceInfobloxCNAMERecordUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*infoblox.Client)

	opts := &infoblox.Options{
		ReturnFields: []string{"canonical", "name", "comment", "ttl", "view"},
	}
	_, err := client.GetRecordCname(d.Id(), opts)
	if err != nil {
		return fmt.Errorf("error finding infoblox CNAME record: %s", err.Error())
	}

	record := url.Values{}
	record.Add("canonical", d.Get("canonical").(string))
	record.Add("name", d.Get("name").(string))
	populateSharedAttributes(d, &record)

	log.Printf("[DEBUG] Updating Infoblox CNAME record with configuration: %#v", record)

	recordID, err := client.RecordCnameObject(d.Id()).Update(record, opts, nil)
	if err != nil {
		return fmt.Errorf("error updating Infoblox CNAME record: %s", err.Error())
	}

	d.SetId(recordID)
	log.Printf("[INFO] Infoblox CNAME record updated with ID: %s", d.Id())

	return resourceInfobloxCNAMERecordRead(d, meta)
}

func resourceInfobloxCNAMERecordDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*infoblox.Client)

	log.Printf("[DEBUG] Deleting Infoblox CNAME record: %s, %s", d.Get("name").(string), d.Id())
	_, err := client.GetRecordCname(d.Id(), nil)
	if err != nil {
		return fmt.Errorf("error finding Infoblox CNAME record: %s", err.Error())
	}

	err = client.RecordCnameObject(d.Id()).Delete(nil)
	if err != nil {
		return fmt.Errorf("error deleting Infoblox CNAME record: %s", err.Error())
	}

	return nil
}
