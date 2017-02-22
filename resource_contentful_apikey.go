package main

import (
	"github.com/hashicorp/terraform/helper/schema"
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

func resourceCreateAPIKey(d *schema.ResourceData, m interface{}) error {
	configMap := m.(map[string]string)
	cmaToken := configMap["cma_token"]
	spaceID := d.Get("space_id").(string)
	apiKeyName := d.Get("name").(string)

	key, err := createAPIKey(cmaToken, spaceID, apiKeyName)
	if err != nil {
		return err
	}

	if err := setAPIKeyProperties(d, key); err != nil {
		return err
	}

	d.SetId(key.Sys.ID)
	return nil
}

func resourceUpdateAPIKey(d *schema.ResourceData, m interface{}) error {
	configMap := m.(map[string]string)
	cmaToken := configMap["cma_token"]
	spaceID := d.Get("space_id").(string)
	ID := d.Id()
	version := d.Get("version").(int)
	newAPIKeyName := d.Get("name").(string)

	key, err := updateAPIKey(cmaToken, spaceID, ID, version, newAPIKeyName)
	if err != nil {
		return err
	}

	if err := setAPIKeyProperties(d, key); err != nil {
		return err
	}

	d.SetId(key.Sys.ID)
	return nil
}

func resourceReadAPIKey(d *schema.ResourceData, m interface{}) error {
	configMap := m.(map[string]string)
	cmaToken := configMap["cma_token"]
	spaceID := d.Get("space_id").(string)
	ID := d.Id()

	key, err := readAPIKey(cmaToken, spaceID, ID)

	if err != nil {
		d.SetId("")
		return nil
	}

	return setAPIKeyProperties(d, key)
}

func resourceDeleteAPIKey(d *schema.ResourceData, m interface{}) error {
	configMap := m.(map[string]string)
	cmaToken := configMap["cma_token"]
	spaceID := d.Get("space_id").(string)
	ID := d.Id()

	return deleteAPIKey(cmaToken, spaceID, ID)
}

func setAPIKeyProperties(d *schema.ResourceData, key *apiKeyData) error {
	if err := d.Set("space_id", key.Sys.Space.Sys.ID); err != nil {
		return err
	}

	if err := d.Set("version", key.Sys.Version); err != nil {
		return err
	}

	if err := d.Set("name", key.Name); err != nil {
		return err
	}

	if err := d.Set("delivery_api_key", key.AccessToken); err != nil {
		return err
	}

	if err := d.Set("preview_api_key", key.PreviewAPIKey.Sys.ID); err != nil {
		return err
	}

	return nil
}
