package infoblox

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/fanatic/go-infoblox"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceInfobloxRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceInfobloxRecordCreate,
		Read:   resourceInfobloxRecordRead,
		Update: resourceInfobloxRecordUpdate,
		Delete: resourceInfobloxRecordDelete,

		Schema: map[string]*schema.Schema{
			"domain": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
				Type:     schema.TypeInt,
				Optional: true,
				//Default:  3600,
			},
		},
	}
}

func createHostJson(d *schema.ResourceData) map[string]interface{} {
	body := make(map[string]interface{})
	if attr, ok := d.GetOk("value"); ok {
		host_obj := make(map[string]string)
		host_obj["ipv4addr"] = attr.(string)
		// Map<String, String>[] var = {host_obj}
		// body["ipv4addrs"] = var
		body["ipv4addrs"] = [1]map[string]string{host_obj}
	}

	var name string
	if attr, ok := d.GetOk("name"); ok {
		name = attr.(string)
	}
	if attr, ok := d.GetOk("domain"); ok {
		name = strings.Join([]string{name, attr.(string)}, ".")
	}
	body["name"] = name

	if attr, ok := d.GetOk("ttl"); ok {
		body["ttl"] = attr.(int)
	}

	return body
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
			ReturnFields: []string{"ttl", "ipv4addr", "name"},
		}
		recID, err = client.RecordA().Create(record, opts, nil)
	case "AAAA":
		opts := &infoblox.Options{
			ReturnFields: []string{"ttl", "ipv6addr", "name"},
		}
		recID, err = client.RecordAAAA().Create(record, opts, nil)
	case "CNAME":
		opts := &infoblox.Options{
			ReturnFields: []string{"ttl", "canonical", "name"},
		}
		recID, err = client.RecordCname().Create(record, opts, nil)
	case "HOST":
		opts := &infoblox.Options{
			ReturnFields: []string{"ttl", "ipv4addrs", "name"},
		}
		recID, err = client.RecordHost().Create(url.Values{}, opts, createHostJson(d))
	default:
		return fmt.Errorf("resourceInfobloxRecordCreate: unknown type")
	}

	if err != nil {
		return fmt.Errorf("Failed to create Infblox Record: %s", err)
	}

	d.SetId(recID)

	log.Printf("[INFO] record ID: %s", d.Id())

	return resourceInfobloxRecordRead(d, meta)
}

func resourceInfobloxRecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*infoblox.Client)

	switch strings.ToUpper(d.Get("type").(string)) {
	case "A":
		opts := &infoblox.Options{
			ReturnFields: []string{"ttl", "ipv4addr", "name"},
		}
		rec, err := client.GetRecordA(d.Id(), opts)
		if err != nil {
			return fmt.Errorf("Couldn't find Infoblox A record: %s", err)
		}
		d.Set("value", rec.Ipv4Addr)
		d.Set("type", "A")
		fqdn := strings.Split(rec.Name, ".")
		d.Set("name", fqdn[0])
		d.Set("domain", strings.Join(fqdn[1:], "."))
		d.Set("ttl", rec.Ttl)

	case "AAAA":
		opts := &infoblox.Options{
			ReturnFields: []string{"ttl", "ipv6addr", "name"},
		}
		rec, err := client.GetRecordAAAA(d.Id(), opts)
		if err != nil {
			return fmt.Errorf("Couldn't find Infoblox AAAA record: %s", err)
		}
		d.Set("value", rec.Ipv6Addr)
		d.Set("type", "AAAA")
		fqdn := strings.Split(rec.Name, ".")
		d.Set("name", fqdn[0])
		d.Set("domain", strings.Join(fqdn[1:], "."))
		d.Set("ttl", rec.Ttl)

	case "CNAME":
		opts := &infoblox.Options{
			ReturnFields: []string{"ttl", "canoncial", "name"},
		}
		rec, err := client.GetRecordCname(d.Id(), opts)
		if err != nil {
			return fmt.Errorf("Couldn't find Infoblox CNAME record: %s", err)
		}
		d.Set("value", rec.Canonical)
		d.Set("type", "CNAME")
		fqdn := strings.Split(rec.Name, ".")
		d.Set("name", fqdn[0])
		d.Set("domain", strings.Join(fqdn[1:], "."))
		d.Set("ttl", rec.Ttl)
	case "HOST":
		opts := &infoblox.Options{
			ReturnFields: []string{"ttl", "ipv4addrs", "name"},
		}
		rec, err := client.GetRecordHost(d.Id(), opts)
		if err != nil {
			return fmt.Errorf("Couldn't find Infoblox Host record: %s", err)
		}
		d.Set("value", rec.Ipv4Addrs[0].Ipv4Addr)
		d.Set("type", "HOST")
		fqdn := strings.Split(rec.Name, ".")
		d.Set("name", fqdn[0])
		d.Set("domain", strings.Join(fqdn[1:], "."))
		d.Set("ttl", rec.Ttl)
	default:
		return fmt.Errorf("resourceInfobloxRecordRead: unknown type: %v", d.Get("type"))
	}

	return nil
}

func resourceInfobloxRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*infoblox.Client)
	var recID string
	var err, updateErr error
	switch strings.ToUpper(d.Get("type").(string)) {
	case "A":
		// TODO: Ensure nil works.
		// Passing nil as the return options because we aren't using the returned object, just ensuring there is no error.
		_, err = client.GetRecordA(d.Id(), nil)
	case "AAAA":
		// TODO: Ensure nil works.
		// Passing nil as the return options because we aren't using the returned object, just ensuring there is no error.
		_, err = client.GetRecordAAAA(d.Id(), nil)
	case "CNAME":
		// TODO: Ensure nil works.
		// Passing nil as the return options because we aren't using the returned object, just ensuring there is no error.
		_, err = client.GetRecordCname(d.Id(), nil)
	case "HOST":
		// TODO: Ensure nil works.
		// Passing nil as the return options because we aren't using the returned object, just ensuring there is no error.
		_, err = client.GetRecordHost(d.Id(), nil)
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
			ReturnFields: []string{"ttl", "ipv4addr", "name"},
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
	case "HOST":
		opts := &infoblox.Options{
			ReturnFields: []string{"ttl", "ipv4addrs", "name"},
		}
		recID, err = client.RecordHostObject(d.Id()).Update(url.Values{}, opts, createHostJson(d))
	default:
		return fmt.Errorf("resourceInfobloxRecordUpdate: unknown type")
	}

	if updateErr != nil {
		return fmt.Errorf("Failed to update Infblox Record: %s", err)
	}

	d.SetId(recID)

	return resourceInfobloxRecordRead(d, meta)
}

func resourceInfobloxRecordDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*infoblox.Client)

	log.Printf("[INFO] Deleting Infoblox Record: %s, %s", d.Get("name").(string), d.Id())
	switch strings.ToUpper(d.Get("type").(string)) {
	case "A":
		// TODO: Ensure nil works.
		// Passing nil as the return options because we aren't using the returned object, just ensuring there is no error.
		_, err := client.GetRecordA(d.Id(), nil)
		if err != nil {
			return fmt.Errorf("Couldn't find Infoblox A record: %s", err)
		}

		deleteErr := client.RecordAObject(d.Id()).Delete(nil)
		if deleteErr != nil {
			return fmt.Errorf("Error deleting Infoblox A Record: %s", err)
		}
	case "AAAA":
		// TODO: Ensure nil works.
		// Passing nil as the return options because we aren't using the returned object, just ensuring there is no error.
		_, err := client.GetRecordAAAA(d.Id(), nil)
		if err != nil {
			return fmt.Errorf("Couldn't find Infoblox AAAA record: %s", err)
		}

		deleteErr := client.RecordAAAAObject(d.Id()).Delete(nil)
		if deleteErr != nil {
			return fmt.Errorf("Error deleting Infoblox AAAA Record: %s", err)
		}
	case "CNAME":
		// TODO: Ensure nil works.
		// Passing nil as the return options because we aren't using the returned object, just ensuring there is no error.
		_, err := client.GetRecordCname(d.Id(), nil)
		if err != nil {
			return fmt.Errorf("Couldn't find Infoblox CNAME record: %s", err)
		}

		deleteErr := client.RecordCnameObject(d.Id()).Delete(nil)
		if deleteErr != nil {
			return fmt.Errorf("Error deleting Infoblox CNAME Record: %s", err)
		}
	case "HOST":
		// TODO: Ensure nil works.
		// Passing nil as the return options because we aren't using the returned object, just ensuring there is no error.
		_, err := client.GetRecordHost(d.Id(), nil)
		if err != nil {
			return fmt.Errorf("Couldn't find Infoblox HOST record: %s", err)
		}

		deleteErr := client.RecordHostObject(d.Id()).Delete(nil)
		if deleteErr != nil {
			return fmt.Errorf("Error deleting Infoblox HOST Record: %s", err)
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
	case "HOST":
	default:
		return fmt.Errorf("getAll: type not found")
	}

	return nil
}
