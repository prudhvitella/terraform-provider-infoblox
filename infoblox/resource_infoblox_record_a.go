package infoblox

import (
	"fmt"
	"log"
	"net/url"

	infoblox "github.com/fanatic/go-infoblox"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func infobloxRecordA() *schema.Resource {
	return &schema.Resource{
		Create: resourceInfobloxARecordCreate,
		Read:   resourceInfobloxARecordRead,
		Update: resourceInfobloxARecordUpdate,
		Delete: resourceInfobloxARecordDelete,

		Schema: map[string]*schema.Schema{
			"address": &schema.Schema{
				Computed:      true,
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      false,
				ConflictsWith: []string{"cidr"},
				ValidateFunc:  validation.SingleIP(),

				// Because setting a cidr will implicitly fill the address in state
				// we suppress the diff for this field if the cidr attribute is set.
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if d.Get("cidr").(string) != "" && old != "" {
						return true
					}
					return false
				},
			},
			"cidr": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"address"},
				ValidateFunc:  validation.CIDRNetwork(0, 32),
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
				Default:  "default",
			},
		},
	}
}

// aObjectFromAttributes created an infoblox.RecordAObject using the attributes
// as set by terraform.
// The Infoblox WAPI does not allow updates to the "view" field on an A record,
// so we also take a skipView arg to skip setting view.
func aObjectFromAttributes(d *schema.ResourceData, skipView bool) infoblox.RecordAObject {
	aObject := infoblox.RecordAObject{}

	aObject.Name = d.Get("name").(string)
	aObject.Ipv4Addr = d.Get("address").(string)

	if attr, ok := d.GetOk("comment"); ok {
		aObject.Comment = attr.(string)
	}
	if attr, ok := d.GetOk("ttl"); ok {
		aObject.Ttl = attr.(int)
	}
	if skipView {
		return aObject
	}

	if attr, ok := d.GetOk("view"); ok {
		aObject.View = attr.(string)
	}

	return aObject
}

func resourceInfobloxARecordCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*infoblox.Client)

	record := url.Values{}
	aObject := infoblox.RecordAObject{}

	var address string
	var err error
	if v, ok := d.GetOk("address"); ok {
		address = v.(string)
	}
	if v, ok := d.GetOk("cidr"); ok {
		address, err = nextAvailableIP(client, v.(string))
	}
	if err != nil {
		return err
	}
	aObject.Ipv4Addr = address

	aObject.Name = d.Get("name").(string)
	if attr, ok := d.GetOk("comment"); ok {
		aObject.Comment = attr.(string)
	}
	if attr, ok := d.GetOk("ttl"); ok {
		aObject.Ttl = attr.(int)
	}
	if attr, ok := d.GetOk("view"); ok {
		aObject.View = attr.(string)
	}

	log.Printf("[DEBUG] Creating Infoblox A record with configuration: %#v", aObject)

	opts := &infoblox.Options{
		ReturnFields: []string{"ipv4addr", "name", "comment", "ttl", "view"},
	}
	recordID, err := client.RecordA().Create(record, opts, aObject)
	if err != nil {
		return fmt.Errorf("error creating infoblox A record: %s", err.Error())
	}

	d.SetId(recordID)
	log.Printf("[INFO] Infoblox A record created with ID: %s", d.Id())

	return resourceInfobloxARecordRead(d, meta)
}

func resourceInfobloxARecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*infoblox.Client)

	opts := &infoblox.Options{
		ReturnFields: []string{"ipv4addr", "name", "comment", "ttl", "view"},
	}
	record, err := client.GetRecordA(d.Id(), opts)
	if err != nil {
		return handleReadError(d, "A", err)
	}

	d.Set("address", record.Ipv4Addr)
	d.Set("name", record.Name)
	d.Set("comment", record.Comment)
	d.Set("ttl", record.Ttl)
	d.Set("view", record.View)

	return nil
}

func resourceInfobloxARecordUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*infoblox.Client)

	opts := &infoblox.Options{
		ReturnFields: []string{"ipv4addr", "name", "comment", "ttl", "view"},
	}
	_, err := client.GetRecordA(d.Id(), opts)
	if err != nil {
		return fmt.Errorf("error finding infoblox A record: %s", err.Error())
	}

	record := url.Values{}
	aRecordObject := aObjectFromAttributes(d, true)

	log.Printf("[DEBUG] Updating Infoblox A record with configuration: %#v", record)

	recordID, err := client.RecordAObject(d.Id()).Update(record, opts, aRecordObject)
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
