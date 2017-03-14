package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	contentful "github.com/tolgaakyuz/contentful.go"
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

func resourceCreateLocale(d *schema.ResourceData, m interface{}) (err error) {
	configMap := m.(map[string]interface{})
	client := configMap["client"].(*contentful.Contentful)

	space, err := client.GetSpace(d.Get("space_id").(string))
	if err != nil {
		return err
	}

	locale := space.NewLocale()
	locale.Name = d.Get("name").(string)
	locale.Code = d.Get("code").(string)
	locale.FallbackCode = d.Get("fallback_code").(string)
	locale.Optional = d.Get("optional").(bool)

	err = locale.Save()
	if err != nil {
		return err
	}

	err = setLocaleProperties(d, locale)
	if err != nil {
		return err
	}

	d.SetId(locale.Sys.ID)

	return nil
}

func resourceReadLocale(d *schema.ResourceData, m interface{}) error {
	configMap := m.(map[string]interface{})
	client := configMap["client"].(*contentful.Contentful)

	space, err := client.GetSpace(d.Get("space_id").(string))
	if err != nil {
		return err
	}

	locale, err := space.GetLocale(d.Id())
	if _, ok := err.(*contentful.NotFoundError); ok {
		d.SetId("")
		return nil
	}

	if err != nil {
		return err
	}

	return setLocaleProperties(d, locale)
}

func resourceUpdateLocale(d *schema.ResourceData, m interface{}) (err error) {
	configMap := m.(map[string]interface{})
	client := configMap["client"].(*contentful.Contentful)

	space, err := client.GetSpace(d.Get("space_id").(string))
	if err != nil {
		return err
	}

	locale, err := space.GetLocale(d.Id())
	if err != nil {
		return err
	}

	locale.Name = d.Get("name").(string)
	locale.Code = d.Get("code").(string)
	locale.FallbackCode = d.Get("fallback_code").(string)
	locale.Optional = d.Get("optional").(bool)
	err = locale.Save()
	if err != nil {
		return err
	}

	err = setLocaleProperties(d, locale)
	if err != nil {
		return err
	}

	return nil
}

func resourceDeleteLocale(d *schema.ResourceData, m interface{}) (err error) {
	configMap := m.(map[string]interface{})
	client := configMap["client"].(*contentful.Contentful)

	space, err := client.GetSpace(d.Get("space_id").(string))
	if err != nil {
		return err
	}

	locale, err := space.GetLocale(d.Id())
	if err != nil {
		return err
	}

	err = locale.Delete()
	if _, ok := err.(*contentful.NotFoundError); ok {
		return nil
	}

	if err != nil {
		return err
	}

	return nil
}

func setLocaleProperties(d *schema.ResourceData, locale *contentful.Locale) error {
	err := d.Set("name", locale.Name)
	if err != nil {
		return err
	}

	err = d.Set("code", locale.Code)
	if err != nil {
		return err
	}

	err = d.Set("fallbackCode", locale.FallbackCode)
	if err != nil {
		return err
	}

	err = d.Set("optional", locale.Optional)
	if err != nil {
		return err
	}

	return nil
}
