package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func resourceContentfulSpace() *schema.Resource {
	return &schema.Resource{
		Create: resourceSpaceCreate,
		Read:   resourceSpaceRead,
		Update: resourceSpaceUpdate,
		Delete: resourceSpaceDelete,

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			// Space specific props
			"default_locale": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "en",
			},
		},
	}
}

func resourceSpaceCreate(d *schema.ResourceData, m interface{}) error {
	configMap := m.(map[string]string)
	name := d.Get("name").(string)
	defaultLocale := d.Get("default_locale").(string)

	spacejson, err := createSpace(
		configMap["cma_token"],
		configMap["organization_id"],
		name,
		defaultLocale,
	)

	if err != nil {
		return err
	}

	err = updateSpaceProperties(d, spacejson)
	if err != nil {
		return err
	}

	d.SetId(spacejson.Sys.ID)
	return nil
}

func resourceSpaceRead(d *schema.ResourceData, m interface{}) error {
	configMap := m.(map[string]string)
	_, err := readSpace(configMap["cma_token"], d.Id())

	log.Println("big problums", err == errorSpaceNotFound)

	if err == errorSpaceNotFound {
		d.SetId("")
		return nil
	}

	if err != nil {
		return err
	}

	return nil
}

func resourceSpaceUpdate(d *schema.ResourceData, m interface{}) error {
	configMap := m.(map[string]string)
	spaceVersion := d.Get("version").(int)
	newSpaceName := d.Get("name").(string)
	log.Println("spaceVersion", spaceVersion)

	spacejson, err := updateSpace(configMap["cma_token"], d.Id(), spaceVersion, newSpaceName)

	if err != nil {
		return err
	}

	return updateSpaceProperties(d, spacejson)
}

func resourceSpaceDelete(d *schema.ResourceData, m interface{}) error {
	configMap := m.(map[string]string)

	err := deleteSpace(configMap["cma_token"], d.Id())

	if err == errorSpaceNotFound {
		return nil
	}

	return err
}

func updateSpaceProperties(d *schema.ResourceData, spacejson *spaceData) error {
	err := d.Set("version", spacejson.Sys.Version)
	if err != nil {
		return err
	}

	err = d.Set("name", spacejson.Name)
	if err != nil {
		return err
	}

	return nil
}
