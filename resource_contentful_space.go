package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	contentful "github.com/tolgaakyuz/contentful-go"
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

func resourceSpaceCreate(d *schema.ResourceData, m interface{}) (err error) {
	configMap := m.(map[string]interface{})

	client := configMap["client"].(*contentful.Contentful)
	space := client.NewSpace(configMap["organization_id"].(string))
	space.Name = d.Get("name").(string)
	space.DefaultLocale = d.Get("default_locale").(string)
	err = space.Save()
	if err != nil {
		return err
	}

	err = updateSpaceProperties(d, space)
	if err != nil {
		return err
	}

	d.SetId(space.Sys.ID)

	return nil
}

func resourceSpaceRead(d *schema.ResourceData, m interface{}) error {
	configMap := m.(map[string]interface{})

	client := configMap["client"].(*contentful.Contentful)
	_, err := client.GetSpace(d.Id())

	if _, ok := err.(contentful.NotFoundError); ok {
		d.SetId("")
		return nil
	}

	return err
}

func resourceSpaceUpdate(d *schema.ResourceData, m interface{}) (err error) {
	configMap := m.(map[string]interface{})

	client := configMap["client"].(*contentful.Contentful)
	space, err := client.GetSpace(d.Id())
	if err != nil {
		return err
	}

	space.Name = d.Get("name").(string)
	err = space.Save()
	if err != nil {
		return err
	}

	return updateSpaceProperties(d, space)
}

func resourceSpaceDelete(d *schema.ResourceData, m interface{}) (err error) {
	configMap := m.(map[string]interface{})

	client := configMap["client"].(*contentful.Contentful)
	space, err := client.GetSpace(d.Id())
	if err != nil {
		return err
	}

	err = space.Delete()
	if _, ok := err.(contentful.NotFoundError); ok {
		return nil
	}

	return err
}

func updateSpaceProperties(d *schema.ResourceData, space *contentful.Space) error {
	err := d.Set("version", space.Sys.Version)
	if err != nil {
		return err
	}

	err = d.Set("name", space.Name)
	if err != nil {
		return err
	}

	return nil
}
