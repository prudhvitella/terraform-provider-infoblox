package infoblox

import (
	"fmt"
	"log"
	"net/url"
	"strconv"

	infoblox "github.com/fanatic/go-infoblox"
	"github.com/hashicorp/terraform/helper/schema"
)

func infobloxRecordSRV() *schema.Resource {
	return &schema.Resource{
		Create: resourceInfobloxSRVRecordCreate,
		Read:   resourceInfobloxSRVRecordRead,
		Update: resourceInfobloxSRVRecordUpdate,
		Delete: resourceInfobloxSRVRecordDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"port": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: false,
			},
			"priority": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: false,
			},
			"target": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"weight": &schema.Schema{
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

func resourceInfobloxSRVRecordCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*infoblox.Client)

	record := url.Values{}
	record.Add("name", d.Get("name").(string))
	record.Add("port", strconv.Itoa(d.Get("port").(int)))
	record.Add("priority", strconv.Itoa(d.Get("priority").(int)))
	record.Add("target", d.Get("target").(string))
	record.Add("weight", strconv.Itoa(d.Get("weight").(int)))
	populateSharedAttributes(d, &record)

	log.Printf("[DEBUG] Creating Infoblox SRV record with configuration: %#v", record)

	opts := &infoblox.Options{
		ReturnFields: []string{"name", "port", "priority", "target", "weight", "comment", "ttl", "view"},
	}

	// TODO: Add SRV support to go-infoblox
	recordID, err := client.RecordSrv().Create(record, opts, nil)

	if err != nil {
		return fmt.Errorf("error creating infoblox SRV record: %s", err.Error())
	}

	d.SetId(recordID)
	log.Printf("[INFO] Infoblox SRV record created with ID: %s", d.Id())

	return nil
}

func resourceInfobloxSRVRecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*infoblox.Client)

	opts := &infoblox.Options{
		ReturnFields: []string{"name", "port", "priority", "target", "weight", "comment", "ttl", "view"},
	}
	record, err := client.GetRecordSrv(d.Id(), opts)
	if err != nil {
		return handleReadError(d, "SRV", err)
	}

	d.Set("name", record.Name)
	d.Set("port", record.Port)
	d.Set("priority", record.Priority)
	d.Set("target", record.Target)
	d.Set("weight", record.Weight)

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

func resourceInfobloxSRVRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*infoblox.Client)

	opts := &infoblox.Options{
		ReturnFields: []string{"name", "port", "priority", "target", "weight", "comment", "ttl", "view"},
	}
	_, err := client.GetRecordSrv(d.Id(), opts)
	if err != nil {
		return fmt.Errorf("error finding infoblox SRV record: %s", err.Error())
	}

	record := url.Values{}
	// name string
	// port int
	// priority int
	// target fqdn/string
	// weight int
	// shared
	record.Add("name", d.Get("name").(string))
	record.Add("port", strconv.Itoa(d.Get("port").(int)))
	record.Add("priority", strconv.Itoa(d.Get("priority").(int)))
	record.Add("target", d.Get("target").(string))
	record.Add("weight", strconv.Itoa(d.Get("weight").(int)))
	populateSharedAttributes(d, &record)

	log.Printf("[DEBUG] Updating Infoblox SRV record with configuration: %#v", record)

	recordID, err := client.RecordSrvObject(d.Id()).Update(record, opts, nil)
	if err != nil {
		return fmt.Errorf("error updating Infoblox SRV record: %s", err.Error())
	}

	d.SetId(recordID)
	log.Printf("[INFO] Infoblox SRV record updated with ID: %s", d.Id())

	return resourceInfobloxSRVRecordRead(d, meta)
}

func resourceInfobloxSRVRecordDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*infoblox.Client)

	log.Printf("[DEBUG] Deleting Infoblox SRV record: %s, %s", d.Get("name").(string), d.Id())
	_, err := client.GetRecordSrv(d.Id(), nil)
	if err != nil {
		return fmt.Errorf("error finding Infoblox SRV record: %s", err.Error())
	}

	err = client.RecordSrvObject(d.Id()).Delete(nil)
	if err != nil {
		return fmt.Errorf("error deleting Infoblox SRV record: %s", err.Error())
	}

	return nil
}
