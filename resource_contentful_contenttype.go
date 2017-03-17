package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	contentful "github.com/tolgaakyuz/contentful.go"
)

func resourceContentfulContentType() *schema.Resource {
	return &schema.Resource{
		Create: resourceContentTypeCreate,
		Read:   resourceContentTypeRead,
		Update: resourceContentTypeUpdate,
		Delete: resourceContentTypeDelete,

		Schema: map[string]*schema.Schema{
			"space_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
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
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"displayField": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"field": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						//@TODO Add ValidateFunc to validate field type
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"required": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						// "localized": &schema.Schema{
						// 	Type:     schema.TypeBool,
						// 	Optional: true,
						// 	Default:  false,
						// },
						"disabled": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
		},
	}
}

func resourceContentTypeCreate(d *schema.ResourceData, m interface{}) (err error) {
	configMap := m.(map[string]interface{})
	client := configMap["client"].(*contentful.Contentful)

	space, err := client.GetSpace(d.Get("space_id").(string))
	if err != nil {
		return err
	}

	ct := space.NewContentType()
	ct.Name = d.Get("name").(string)
	ct.DisplayField = d.Get("displayField").(string)
	ct.Description = d.Get("description").(string)

	for _, field := range d.Get("field").(*schema.Set).List() {
		ct.Fields = append(ct.Fields, &contentful.Field{
			ID:   field.(map[string]interface{})["id"].(string),
			Name: field.(map[string]interface{})["name"].(string),
			Type: field.(map[string]interface{})["type"].(string),
			// Localized: field.(map[string]interface{})["localized"].(bool),
			Required: field.(map[string]interface{})["required"].(bool),
			Disabled: field.(map[string]interface{})["disabled"].(bool),
		})
	}

	if err = ct.Save(); err != nil {
		return err
	}

	if err = ct.Activate(); err != nil {
		//@TODO Maybe delete the CT ?
		return err
	}

	if err = d.Set("version", ct.Sys.Version); err != nil {
		return err
	}

	d.SetId(ct.Sys.ID)

	return nil
}

func resourceContentTypeRead(d *schema.ResourceData, m interface{}) (err error) {
	configMap := m.(map[string]interface{})

	client := configMap["client"].(*contentful.Contentful)
	space, err := client.GetSpace(d.Get("space_id").(string))
	if err != nil {
		return err
	}

	_, err = space.GetContentType(d.Id())

	return err
}

func resourceContentTypeUpdate(d *schema.ResourceData, m interface{}) (err error) {
	configMap := m.(map[string]interface{})

	client := configMap["client"].(*contentful.Contentful)
	space, err := client.GetSpace(d.Get("space_id").(string))
	if err != nil {
		return err
	}

	ct, err := space.GetContentType(d.Id())
	if err != nil {
		return err
	}

	ct.Name = d.Get("name").(string)
	ct.DisplayField = d.Get("displayField").(string)
	ct.Description = d.Get("description").(string)

	var fields []*contentful.Field
	for _, field := range d.Get("field").(*schema.Set).List() {
		fields = append(fields, &contentful.Field{
			ID:   field.(map[string]interface{})["id"].(string),
			Name: field.(map[string]interface{})["name"].(string),
			Type: field.(map[string]interface{})["type"].(string),
			// Localized: field.(map[string]interface{})["localized"].(bool),
			Required: field.(map[string]interface{})["required"].(bool),
			Disabled: field.(map[string]interface{})["disabled"].(bool),
		})
	}

	ct.Fields = fields

	if err = ct.Save(); err != nil {
		return err
	}

	if err = ct.Activate(); err != nil {
		//@TODO Maybe delete the CT ?
		return err
	}

	if err = setContentTypeProperties(d, ct); err != nil {
		return err
	}

	return nil
}

func resourceContentTypeDelete(d *schema.ResourceData, m interface{}) (err error) {
	configMap := m.(map[string]interface{})

	client := configMap["client"].(*contentful.Contentful)
	space, err := client.GetSpace(d.Get("space_id").(string))
	if err != nil {
		return err
	}

	ct, err := space.GetContentType(d.Id())
	if err != nil {
		return err
	}

	if err = ct.Delete(); err != nil {
		return err
	}

	return nil
}

func setContentTypeProperties(d *schema.ResourceData, ct *contentful.ContentType) (err error) {
	if err = d.Set("space_id", ct.Sys.Space.ID()); err != nil {
		return err
	}

	if err = d.Set("version", ct.Sys.Version); err != nil {
		return err
	}

	if err = d.Set("name", ct.Name); err != nil {
		return err
	}

	if err = d.Set("description", ct.Description); err != nil {
		return err
	}

	if err = d.Set("displayField", ct.DisplayField); err != nil {
		return err
	}

	return nil
}
