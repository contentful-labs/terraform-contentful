package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceContentfulLocale() *schema.Resource {
	return &schema.Resource{
		Create: resourceCreateLocale,
		Read:   resourceReadLocale,
		Update: resourceUpdateLocale,
		Delete: resourceDeleteLocale,

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
			// Locale specific props
			"code": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"fallback_code": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "en-US",
			},
			"optional": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceCreateLocale(d *schema.ResourceData, m interface{}) error {
	configMap := m.(map[string]string)
	cmaToken := configMap["cma_token"]
	spaceID := d.Get("space_id").(string)

	localeProps := localeProperties{
		Name:         d.Get("name").(string),
		Code:         d.Get("code").(string),
		FallbackCode: d.Get("fallback_code").(string),
		Optional:     d.Get("optional").(bool),
	}

	loc, err := createLocale(
		cmaToken,
		spaceID,
		localeProps,
	)
	if err != nil {
		return err
	}

	err = setLocaleProperties(d, loc)
	if err != nil {
		return err
	}

	d.SetId(loc.Sys.ID)
	return nil
}

func resourceReadLocale(d *schema.ResourceData, m interface{}) error {
	configMap := m.(map[string]string)
	cmaToken := configMap["cma_token"]
	spaceID := d.Get("space_id").(string)
	localeID := d.Id()

	loc, err := readLocale(cmaToken, spaceID, localeID)

	if err == errorLocaleNotFound {
		d.SetId("")
		return nil
	}

	return setLocaleProperties(d, loc)
}

func resourceUpdateLocale(d *schema.ResourceData, m interface{}) error {
	configMap := m.(map[string]string)
	cmaToken := configMap["cma_token"]
	spaceID := d.Get("space_id").(string)
	localeID := d.Id()

	localeProps := localeProperties{
		Name:         d.Get("name").(string),
		Code:         d.Get("code").(string),
		FallbackCode: d.Get("fallback_code").(string),
		Optional:     d.Get("optional").(bool),
	}

	loc, err := updateLocale(
		cmaToken,
		spaceID,
		localeID,
		localeProps,
	)
	if err != nil {
		return err
	}

	err = setLocaleProperties(d, loc)
	if err != nil {
		return err
	}

	return nil
}

func resourceDeleteLocale(d *schema.ResourceData, m interface{}) error {
	configMap := m.(map[string]string)
	cmaToken := configMap["cma_token"]
	spaceID := d.Get("space_id").(string)
	localeID := d.Id()

	err := deleteLocale(cmaToken, spaceID, localeID)

	if err == errorLocaleNotFound {
		return nil
	}

	return err
}

func setLocaleProperties(d *schema.ResourceData, loc *locale) error {
	err := d.Set("name", loc.Name)
	if err != nil {
		return err
	}

	err = d.Set("code", loc.Code)
	if err != nil {
		return err
	}

	err = d.Set("fallbackCode", loc.FallbackCode)
	if err != nil {
		return err
	}

	err = d.Set("optional", loc.Optional)
	if err != nil {
		return err
	}

	return nil
}
