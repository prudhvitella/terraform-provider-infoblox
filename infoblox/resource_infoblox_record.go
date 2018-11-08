package infoblox

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/fanatic/go-infoblox"
	"github.com/hashicorp/terraform/helper/schema"
)

var deprecated = `The entire 'infoblox_record'
resource is deprecated and will no longer see active development. It is
recommended you use the dedicated infoblox_record_* resources instead.`

func resourceInfobloxRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceInfobloxRecordCreate,
		Read:   resourceInfobloxRecordRead,
		Update: resourceInfobloxRecordUpdate,
		Delete: resourceInfobloxRecordDelete,

		Schema: map[string]*schema.Schema{
			"domain": &schema.Schema{
				Type:       schema.TypeString,
				Required:   true,
				ForceNew:   true,
				Deprecated: deprecated,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"value": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"ttl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "3600",
			},

			"view": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "default",
			},
		},
	}
}

func resourceInfobloxRecordCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*infoblox.Client)

	record := url.Values{}
	if err := getAll(d, record); err != nil {
		return err
	}

	log.Printf("[DEBUG] Infoblox Record create configuration: %#v", record)

	var recID string
	var err error

	switch strings.ToUpper(d.Get("type").(string)) {
	case "A":
		opts := &infoblox.Options{
			ReturnFields: []string{"ttl", "ipv4addr", "name", "view"},
		}
		recID, err = client.RecordA().Create(record, opts, nil)
	case "AAAA":
		opts := &infoblox.Options{
			ReturnFields: []string{"ttl", "ipv6addr", "name", "view"},
		}
		recID, err = client.RecordAAAA().Create(record, opts, nil)
	case "CNAME":
		opts := &infoblox.Options{
			ReturnFields: []string{"ttl", "canonical", "name", "view"},
		}
		recID, err = client.RecordCname().Create(record, opts, nil)
	default:
		return fmt.Errorf("resourceInfobloxRecordCreate: unknown type")
	}

	if err != nil {
		return fmt.Errorf("Failed to create Infoblox Record: %s", err.Error())
	}

	d.SetId(recID)
	log.Printf("[INFO] record ID: %s", d.Id())

	return resourceInfobloxRecordRead(d, meta)
}

func resourceInfobloxRecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*infoblox.Client)

	switch strings.ToUpper(d.Get("type").(string)) {
	case "A":
		rec, err := client.GetRecordA(d.Id(), nil)
		if err != nil {
			return handleReadError(d, "A", err)
		}

		d.Set("value", rec.Ipv4Addr)
		d.Set("type", "A")
		fqdn := strings.Split(rec.Name, ".")
		d.Set("name", fqdn[0])
		d.Set("domain", strings.Join(fqdn[1:], "."))
		d.Set("ttl", rec.Ttl)
		d.Set("view", rec.View)

	case "AAAA":
		rec, err := client.GetRecordAAAA(d.Id(), nil)
		if err != nil {
			return handleReadError(d, "AAAA", err)
		}
		d.Set("value", rec.Ipv6Addr)
		d.Set("type", "AAAA")
		fqdn := strings.Split(rec.Name, ".")
		d.Set("name", fqdn[0])
		d.Set("domain", strings.Join(fqdn[1:], "."))
		d.Set("ttl", rec.Ttl)
		d.Set("view", rec.View)

	case "CNAME":
		rec, err := client.GetRecordCname(d.Id(), nil)
		if err != nil {
			return handleReadError(d, "CNAME", err)
		}
		d.Set("value", rec.Canonical)
		d.Set("type", "CNAME")
		fqdn := strings.Split(rec.Name, ".")
		d.Set("name", fqdn[0])
		d.Set("domain", strings.Join(fqdn[1:], "."))
		d.Set("ttl", rec.Ttl)
		d.Set("view", rec.View)
	default:
		return fmt.Errorf("resourceInfobloxRecordRead: unknown type")
	}

	return nil
}

func resourceInfobloxRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*infoblox.Client)
	var recID string
	var err, updateErr error
	switch strings.ToUpper(d.Get("type").(string)) {
	case "A":
		_, err = client.GetRecordA(d.Id(), nil)
	case "AAAA":
		_, err = client.GetRecordAAAA(d.Id(), nil)
	case "CNAME":
		_, err = client.GetRecordCname(d.Id(), nil)
	default:
		return fmt.Errorf("resourceInfobloxRecordUpdate: unknown type")
	}

	if err != nil {
		return fmt.Errorf("Couldn't find Infoblox record: %s", err)
	}

	record := url.Values{}
	if err := getAll(d, record); err != nil {
		return err
	}

	log.Printf("[DEBUG] Infoblox Record update configuration: %#v", record)

	switch strings.ToUpper(d.Get("type").(string)) {
	case "A":
		opts := &infoblox.Options{
			ReturnFields: []string{"ttl", "ipv4addr", "name", "view"},
		}
		recID, updateErr = client.RecordAObject(d.Id()).Update(record, opts, nil)
	case "AAAA":
		opts := &infoblox.Options{
			ReturnFields: []string{"ttl", "ipv6addr", "name"},
		}
		recID, updateErr = client.RecordAAAAObject(d.Id()).Update(record, opts, nil)
	case "CNAME":
		opts := &infoblox.Options{
			ReturnFields: []string{"ttl", "canonical", "name"},
		}
		recID, updateErr = client.RecordCnameObject(d.Id()).Update(record, opts, nil)
	default:
		return fmt.Errorf("resourceInfobloxRecordUpdate: unknown type")
	}

	if updateErr != nil {
		return fmt.Errorf("Failed to update Infoblox Record: %s", err.Error())
	}

	d.SetId(recID)

	return resourceInfobloxRecordRead(d, meta)
}

func resourceInfobloxRecordDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*infoblox.Client)

	log.Printf("[INFO] Deleting Infoblox Record: %s, %s", d.Get("name").(string), d.Id())
	switch strings.ToUpper(d.Get("type").(string)) {
	case "A":
		_, err := client.GetRecordA(d.Id(), nil)
		if err != nil {
			return fmt.Errorf("Couldn't find Infoblox A record: %s", err)
		}

		deleteErr := client.RecordAObject(d.Id()).Delete(nil)
		if deleteErr != nil {
			return fmt.Errorf("Error deleting Infoblox A Record: %s", deleteErr)
		}
	case "AAAA":
		_, err := client.GetRecordAAAA(d.Id(), nil)
		if err != nil {
			return fmt.Errorf("Couldn't find Infoblox AAAA record: %s", err)
		}

		deleteErr := client.RecordAAAAObject(d.Id()).Delete(nil)
		if deleteErr != nil {
			return fmt.Errorf("Error deleting Infoblox AAAA Record: %s", deleteErr)
		}
	case "CNAME":
		_, err := client.GetRecordCname(d.Id(), nil)
		if err != nil {
			return fmt.Errorf("Couldn't find Infoblox CNAME record: %s", err)
		}

		deleteErr := client.RecordCnameObject(d.Id()).Delete(nil)
		if deleteErr != nil {
			return fmt.Errorf("Error deleting Infoblox CNAME Record: %s", deleteErr)
		}
	default:
		return fmt.Errorf("resourceInfobloxRecordDelete: unknown type")
	}
	return nil
}

func getAll(d *schema.ResourceData, record url.Values) error {
	if attr, ok := d.GetOk("name"); ok {
		record.Set("name", attr.(string))
	}

	if attr, ok := d.GetOk("domain"); ok {
		record.Set("name", strings.Join([]string{record.Get("name"), attr.(string)}, "."))
	}

	if attr, ok := d.GetOk("ttl"); ok {
		record.Set("ttl", attr.(string))
	}

	if attr, ok := d.GetOk("view"); ok {
		record.Set("view", attr.(string))
	}

	var value string
	if attr, ok := d.GetOk("value"); ok {
		value = attr.(string)
	}

	switch strings.ToUpper(d.Get("type").(string)) {
	case "A":
		record.Set("ipv4addr", value)
	case "AAAA":
		record.Set("ipv6addr", value)
	case "CNAME":
		record.Set("canonical", value)
	default:
		return fmt.Errorf("getAll: type not found")
	}

	return nil
}
