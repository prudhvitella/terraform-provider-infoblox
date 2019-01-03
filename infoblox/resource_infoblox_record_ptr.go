package infoblox

import (
	"fmt"
	"log"
	"net/url"

	infoblox "github.com/fanatic/go-infoblox"
	"github.com/hashicorp/terraform/helper/schema"
)

func infobloxRecordPTR() *schema.Resource {
	return &schema.Resource{
		Create: resourceInfobloxPTRRecordCreate,
		Read:   resourceInfobloxPTRRecordRead,
		Update: resourceInfobloxPTRRecordUpdate,
		Delete: resourceInfobloxPTRRecordDelete,

		Schema: map[string]*schema.Schema{
			"address": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"name"},
			},
			"ptrdname": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"address"},
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

func resourceInfobloxPTRRecordCreate(d *schema.ResourceData, meta interface{}) error {
	err := validatePTRFields(d)
	if err != nil {
		return err
	}

	client := meta.(*infoblox.Client)
	record := url.Values{}

	if attr, ok := d.GetOk("address"); ok {
		addressType, err := ipType(attr.(string))
		if err != nil {
			return err
		}
		record.Add(addressType, attr.(string))
	} else {
		record.Add("name", d.Get("name").(string))
	}
	record.Add("ptrdname", d.Get("ptrdname").(string))
	populateSharedAttributes(d, &record)

	log.Printf("[DEBUG] Creating Infoblox PTR record with configuration: %#v", record)

	opts := ptrOpts(d)
	recordID, err := client.RecordPtr().Create(record, opts, nil)

	if err != nil {
		return fmt.Errorf("error creating infoblox PTR record: %s", err.Error())
	}

	d.SetId(recordID)
	log.Printf("[INFO] Infoblox PTR record created with ID: %s", d.Id())

	return nil
}

func resourceInfobloxPTRRecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*infoblox.Client)

	opts := ptrOpts(d)
	record, err := client.GetRecordPtr(d.Id(), opts)
	if err != nil {
		return handleReadError(d, "PTR", err)
	}

	d.Set("ptrdname", record.PtrDname)

	if &record.Ipv4Addr != nil {
		d.Set("address", record.Ipv4Addr)
	}
	if &record.Ipv6Addr != nil {
		d.Set("address", record.Ipv6Addr)
	}
	if &record.Name != nil {
		d.Set("name", record.Name)
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

	return nil
}

func resourceInfobloxPTRRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*infoblox.Client)
	record := url.Values{}

	opts := ptrOpts(d)
	_, err := client.GetRecordPtr(d.Id(), opts)
	if err != nil {
		return fmt.Errorf("error finding infoblox PTR record: %s", err.Error())
	}

	if attr, ok := d.GetOk("address"); ok {
		addressType, err := ipType(attr.(string))
		if err != nil {
			return err
		}
		record.Add(addressType, attr.(string))
	} else {
		record.Add("name", d.Get("name").(string))
	}
	record.Add("ptrdname", d.Get("ptrdname").(string))
	populateSharedAttributes(d, &record)

	log.Printf("[DEBUG] Updating Infoblox PTR record with configuration: %#v", record)

	recordID, err := client.RecordPtrObject(d.Id()).Update(record, opts, nil)
	if err != nil {
		return fmt.Errorf("error updating Infoblox PTR record: %s", err.Error())
	}

	d.SetId(recordID)
	log.Printf("[INFO] Infoblox PTR record updated with ID: %s", d.Id())

	return resourceInfobloxPTRRecordRead(d, meta)
}

func resourceInfobloxPTRRecordDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*infoblox.Client)

	log.Printf("[DEBUG] Deleting Infoblox PTR record: %s, %s", d.Get("ptrdname").(string), d.Id())
	_, err := client.GetRecordPtr(d.Id(), nil)
	if err != nil {
		return fmt.Errorf("error finding Infoblox PTR record: %s", err.Error())
	}

	err = client.RecordPtrObject(d.Id()).Delete(nil)
	if err != nil {
		return fmt.Errorf("error deleting Infoblox PTR record: %s", err.Error())
	}

	return nil
}

// Returns an error if neither address and name are set or if both are set
func validatePTRFields(d *schema.ResourceData) error {
	_, hasAddress := d.GetOk("address")
	_, hasName := d.GetOk("name")
	if hasAddress && hasName {
		return fmt.Errorf("you must specify name or address for PTR record, not both")
	}
	if !(hasAddress || hasName) {
		return fmt.Errorf("you must specify a name of an address for PTR record")
	}

	return nil
}

// A PTR object can have either an ipv4/ipv6 address or a name, so when
// constructing our ReturnFields slice we read the Schema to see which we want
// to return.
func ptrOpts(d *schema.ResourceData) *infoblox.Options {
	opts := []string{"ptrdname", "ttl", "comment", "view"}

	if _, ok := d.GetOk("name"); ok {
		opts = append(opts, "name")
	}
	if _, ok := d.GetOk("ipv4addr"); ok {
		opts = append(opts, "ipv4addr")
	}
	if _, ok := d.GetOk("ipv6addr"); ok {
		opts = append(opts, "ipv6addr")
	}

	return &infoblox.Options{ReturnFields: opts}
}
