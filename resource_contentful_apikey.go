package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	contentful "github.com/tolgaakyuz/contentful-go"
)

func resourceContentfulAPIKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceCreateAPIKey,
		Read:   resourceReadAPIKey,
		Update: resourceUpdateAPIKey,
		Delete: resourceDeleteAPIKey,

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"space_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			// Webhook specific props
			"delivery_api_key": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"preview_api_key": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCreateAPIKey(d *schema.ResourceData, m interface{}) (err error) {
	client := m.(*contentful.Contentful)

	space, err := client.GetSpace(d.Get("space_id").(string))
	if err != nil {
		return err
	}

	apiKey := space.NewAPIKey()
	apiKey.Name = d.Get("name").(string)
	err = apiKey.Save()
	if err != nil {
		return err
	}

	if err := setAPIKeyProperties(d, apiKey); err != nil {
		return err
	}

	d.SetId(apiKey.Sys.ID)

	return nil
}

func resourceUpdateAPIKey(d *schema.ResourceData, m interface{}) (err error) {
	client := m.(*contentful.Contentful)

	space, err := client.GetSpace(d.Get("space_id").(string))
	if err != nil {
		return err
	}

	apiKey, err := space.GetAPIKey(d.Id())
	if err != nil {
		return err
	}

	apiKey.Name = d.Get("name").(string)
	err = apiKey.Save()
	if err != nil {
		return err
	}

	if err := setAPIKeyProperties(d, apiKey); err != nil {
		return err
	}

	d.SetId(apiKey.Sys.ID)

	return nil
}

func resourceReadAPIKey(d *schema.ResourceData, m interface{}) (err error) {
	client := m.(*contentful.Contentful)

	space, err := client.GetSpace(d.Get("space_id").(string))
	if err != nil {
		return err
	}

	apiKey, err := space.GetAPIKey(d.Id())
	if _, ok := err.(contentful.NotFoundError); ok {
		d.SetId("")
		return nil
	}

	return setAPIKeyProperties(d, apiKey)
}

func resourceDeleteAPIKey(d *schema.ResourceData, m interface{}) (err error) {
	client := m.(*contentful.Contentful)

	space, err := client.GetSpace(d.Get("space_id").(string))
	if err != nil {
		return err
	}

	apiKey, err := space.GetAPIKey(d.Id())
	if err != nil {
		return err
	}

	return apiKey.Delete()
}

func setAPIKeyProperties(d *schema.ResourceData, apiKey *contentful.APIKey) error {
	if err := d.Set("space_id", apiKey.Sys.Space.Sys.ID); err != nil {
		return err
	}

	if err := d.Set("version", apiKey.Sys.Version); err != nil {
		return err
	}

	if err := d.Set("name", apiKey.Name); err != nil {
		return err
	}

	if err := d.Set("delivery_api_key", apiKey.AccessToken); err != nil {
		return err
	}

	if err := d.Set("preview_api_key", apiKey.PreviewAPIKey.Sys.ID); err != nil {
		return err
	}

	return nil
}
