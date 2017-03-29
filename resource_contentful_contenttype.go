package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	contentful "github.com/tolgaakyuz/contentful-go"
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
			"display_field": &schema.Schema{
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
						"items": &schema.Schema{
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
									"validations": &schema.Schema{
										Type:     schema.TypeSet,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"link_content_type": &schema.Schema{
													Type:     schema.TypeList,
													Optional: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"link_mimetype_group": &schema.Schema{
													Type:     schema.TypeList,
													Optional: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"size": &schema.Schema{
													Type:     schema.TypeMap,
													Optional: true,
													Elem:     &schema.Schema{Type: schema.TypeFloat},
												},
											},
										},
									},
									"link_type": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"required": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"localized": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"disabled": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"omitted": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"validations": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func resourceContentTypeCreate(d *schema.ResourceData, m interface{}) (err error) {
	client := m.(*contentful.Contentful)
	spaceID := d.Get("space_id").(string)

	ct := &contentful.ContentType{
		Name:         d.Get("name").(string),
		DisplayField: d.Get("display_field").(string),
		Description:  d.Get("description").(string),
		Fields:       []*contentful.Field{},
	}

	for _, rawField := range d.Get("field").(*schema.Set).List() {
		field := rawField.(map[string]interface{})

		contentfulField := &contentful.Field{
			ID:        field["id"].(string),
			Name:      field["name"].(string),
			Type:      field["type"].(string),
			Localized: field["localized"].(bool),
			Required:  field["required"].(bool),
			Disabled:  field["disabled"].(bool),
			Omitted:   field["omitted"].(bool),
		}

		if validations, ok := field["validations"].([]interface{}); ok {
			parsedValidations, err := contentful.ParseValidations(validations)
			if err != nil {
				return err
			}

			contentfulField.Validations = parsedValidations
		}

		if items := processItems(field["items"].(*schema.Set)); items != nil {
			contentfulField.Items = items
		}

		ct.Fields = append(ct.Fields, contentfulField)
	}

	if err = client.ContentTypes.Upsert(spaceID, ct); err != nil {
		return err
	}

	if err = client.ContentTypes.Activate(spaceID, ct); err != nil {
		//@TODO Maybe delete the CT ?
		return err
	}

	if err = setContentTypeProperties(d, ct); err != nil {
		return err
	}

	d.SetId(ct.Sys.ID)

	return nil
}

func resourceContentTypeRead(d *schema.ResourceData, m interface{}) (err error) {
	client := m.(*contentful.Contentful)
	spaceID := d.Get("space_id").(string)

	_, err = client.ContentTypes.Get(spaceID, d.Id())

	return err
}

func resourceContentTypeUpdate(d *schema.ResourceData, m interface{}) (err error) {
	var existingFields []*contentful.Field
	var deletedFields []*contentful.Field

	client := m.(*contentful.Contentful)
	spaceID := d.Get("space_id").(string)

	ct, err := client.ContentTypes.Get(spaceID, d.Id())
	if err != nil {
		return err
	}

	ct.Name = d.Get("name").(string)
	ct.DisplayField = d.Get("display_field").(string)
	ct.Description = d.Get("description").(string)

	// Figure out if fields were removed
	if d.HasChange("field") {
		old, new := d.GetChange("field")

		existingFields, deletedFields = checkFieldChanges(old.(*schema.Set), new.(*schema.Set))

		ct.Fields = existingFields

		if deletedFields != nil {
			ct.Fields = append(ct.Fields, deletedFields...)
		}
	}

	if err = client.ContentTypes.Upsert(spaceID, ct); err != nil {
		return err
	}

	if err = client.ContentTypes.Activate(spaceID, ct); err != nil {
		//@TODO Maybe delete the CT ?
		return err
	}

	if deletedFields != nil {
		ct.Fields = existingFields

		if err = client.ContentTypes.Upsert(spaceID, ct); err != nil {
			return err
		}

		if err = client.ContentTypes.Activate(spaceID, ct); err != nil {
			//@TODO Maybe delete the CT ?
			return err
		}
	}

	return setContentTypeProperties(d, ct)
}

func resourceContentTypeDelete(d *schema.ResourceData, m interface{}) (err error) {
	client := m.(*contentful.Contentful)
	spaceID := d.Get("space_id").(string)

	ct, err := client.ContentTypes.Get(spaceID, d.Id())
	if err != nil {
		return err
	}

	err = client.ContentTypes.Deactivate(spaceID, ct)
	if err != nil {
		return err
	}

	if err = client.ContentTypes.Delete(spaceID, ct); err != nil {
		return err
	}

	return nil
}

func setContentTypeProperties(d *schema.ResourceData, ct *contentful.ContentType) (err error) {

	if err = d.Set("version", ct.Sys.Version); err != nil {
		return err
	}

	return nil
}

func checkFieldChanges(old, new *schema.Set) ([]*contentful.Field, []*contentful.Field) {
	var contentfulField *contentful.Field
	var existingFields []*contentful.Field
	var deletedFields []*contentful.Field
	var fieldRemoved bool

	for _, oldField := range old.List() {
		fieldRemoved = true
		for _, newField := range new.List() {
			if oldField.(map[string]interface{})["id"].(string) == newField.(map[string]interface{})["id"].(string) {
				fieldRemoved = false
				break
			}
		}

		if fieldRemoved {
			deletedFields = append(deletedFields,
				&contentful.Field{
					ID:        oldField.(map[string]interface{})["id"].(string),
					Name:      oldField.(map[string]interface{})["name"].(string),
					Type:      oldField.(map[string]interface{})["type"].(string),
					Localized: oldField.(map[string]interface{})["localized"].(bool),
					Required:  oldField.(map[string]interface{})["required"].(bool),
					Disabled:  oldField.(map[string]interface{})["disabled"].(bool),
					Omitted:   true,
				})
		}
	}

	for _, field := range new.List() {

		contentfulField = &contentful.Field{
			ID:        field.(map[string]interface{})["id"].(string),
			Name:      field.(map[string]interface{})["name"].(string),
			Type:      field.(map[string]interface{})["type"].(string),
			Localized: field.(map[string]interface{})["localized"].(bool),
			Required:  field.(map[string]interface{})["required"].(bool),
			Disabled:  field.(map[string]interface{})["disabled"].(bool),
			Omitted:   field.(map[string]interface{})["omitted"].(bool),
		}

		if items := processItems(field.(map[string]interface{})["items"].(*schema.Set)); items != nil {
			contentfulField.Items = items
		}

		existingFields = append(existingFields, contentfulField)
	}

	return existingFields, deletedFields
}

func processItems(fieldItems *schema.Set) *contentful.FieldTypeArrayItem {
	var items *contentful.FieldTypeArrayItem
	for _, item := range fieldItems.List() {
		var validations []contentful.FieldValidation

		for _, validationList := range item.(map[string]interface{})["validations"].(*schema.Set).List() {

			for key, validation := range validationList.(map[string]interface{}) {

				switch key {
				case "link_content_type":
					var linkList []string
					for _, linkContentType := range validation.([]interface{}) {
						linkList = append(linkList, linkContentType.(string))
					}
					if len(linkList) > 0 {
						validations = append(validations, contentful.FieldValidationLink{LinkContentType: linkList})
					}
				case "link_mimetype_group":
					var mimeGroupList []string
					for _, mimeGroup := range validation.([]interface{}) {
						mimeGroupList = append(mimeGroupList, mimeGroup.(string))
					}
					if len(mimeGroupList) > 0 {
						validations = append(validations, contentful.FieldValidationMimeType{MimeTypes: mimeGroupList})
					}
				case "size":
					if min, ok := validation.(map[string]interface{})["min"]; ok {
						if max, ok := validation.(map[string]interface{})["max"]; ok {
							validations = append(validations, contentful.MinMax{
								Min: min.(float64),
								Max: max.(float64),
							})
						}
					}
				default:
					validations = append(validations, struct{}{})
				}

			}

		}

		items = &contentful.FieldTypeArrayItem{
			Type:        item.(map[string]interface{})["type"].(string),
			Validations: validations,
			LinkType:    item.(map[string]interface{})["link_type"].(string),
		}
	}
	return items
}
