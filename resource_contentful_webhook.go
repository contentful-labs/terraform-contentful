package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceContentfulWebhook() *schema.Resource {
	return &schema.Resource{
		Create: resourceCreateWebhook,
		Read:   resourceReadWebhook,
		Update: resourceUpdateWebhook,
		Delete: resourceDeleteWebhook,

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
			"url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"http_basic_auth_username": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"http_basic_auth_password": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"headers": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
			"topics": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				MinItems: 1,
				Required: true,
			},
		},
	}
}

func resourceCreateWebhook(d *schema.ResourceData, m interface{}) error {
	configMap := m.(map[string]string)
	cmaToken := configMap["cma_token"]
	spaceID := d.Get("space_id").(string)

	headers := transformHeadersToContentfulFormat(d.Get("headers"))
	topics := transformTopicsToContentfulFormat(d.Get("topics"))

	webhookprops := webhookProperties{
		Name:              d.Get("name").(string),
		URL:               d.Get("url").(string),
		Topics:            topics,
		Headers:           headers,
		HTTPBasicUsername: d.Get("http_basic_auth_username").(string),
		HTTPBasicPassword: d.Get("http_basic_auth_password").(string),
	}

	webhookjson, err := createWebhook(
		cmaToken,
		spaceID,
		webhookprops,
	)
	if err != nil {
		return err
	}

	err = setWebhookProperties(d, webhookjson)
	if err != nil {
		return err
	}

	d.SetId(webhookjson.Sys.ID)
	return nil
}

func resourceUpdateWebhook(d *schema.ResourceData, m interface{}) error {
	configMap := m.(map[string]string)
	cmaToken := configMap["cma_token"]
	spaceID := d.Get("space_id").(string)
	webhookID := d.Id()
	version := d.Get("version").(int)

	headers := transformHeadersToContentfulFormat(d.Get("headers"))
	topics := transformTopicsToContentfulFormat(d.Get("topics"))

	webhookprops := webhookProperties{
		Name:              d.Get("name").(string),
		URL:               d.Get("url").(string),
		Topics:            topics,
		Headers:           headers,
		HTTPBasicUsername: d.Get("http_basic_auth_username").(string),
		HTTPBasicPassword: d.Get("http_basic_auth_password").(string),
	}

	wh, err := updateWebhook(cmaToken, spaceID, webhookID, version, webhookprops)
	if err != nil {
		return err
	}

	err = setWebhookProperties(d, wh)
	if err != nil {
		return err
	}

	d.SetId(wh.Sys.ID)
	return nil
}

func resourceReadWebhook(d *schema.ResourceData, m interface{}) error {
	configMap := m.(map[string]string)
	cmaToken := configMap["cma_token"]
	spaceID := d.Get("space_id").(string)
	webhookID := d.Id()
	wh, err := readWebhook(cmaToken, spaceID, webhookID)

	if err == errorWebhookNotFound {
		d.SetId("")
		return nil
	}

	return setWebhookProperties(d, wh)
}

func resourceDeleteWebhook(d *schema.ResourceData, m interface{}) error {
	configMap := m.(map[string]string)
	cmaToken := configMap["cma_token"]
	spaceID := d.Get("space_id").(string)
	webhookID := d.Id()

	err := deleteWebhook(cmaToken, spaceID, webhookID)

	if err == errorSpaceNotFound {
		return nil
	}

	return err
}

func setWebhookProperties(d *schema.ResourceData, webhookjson *webhookData) error {
	headers := make(map[string]string)
	for _, entry := range webhookjson.Headers {
		headers[entry.Key] = entry.Value
	}
	err := d.Set("headers", headers)
	if err != nil {
		return err
	}

	err = d.Set("space_id", webhookjson.Sys.Space.Sys.ID)
	if err != nil {
		return err
	}

	err = d.Set("version", webhookjson.Sys.Version)
	if err != nil {
		return err
	}

	err = d.Set("name", webhookjson.Name)
	if err != nil {
		return err
	}

	err = d.Set("url", webhookjson.URL)
	if err != nil {
		return err
	}

	err = d.Set("http_basic_auth_username", webhookjson.HTTPBasicUsername)
	if err != nil {
		return err
	}

	err = d.Set("http_basic_auth_password", webhookjson.HTTPBasicPassword)
	if err != nil {
		return err
	}

	err = d.Set("topics", webhookjson.Topics)
	if err != nil {
		return err
	}
	return nil
}

func transformHeadersToContentfulFormat(headersTerraform interface{}) []headerKeyValue {
	headers := headersTerraform.(map[string]interface{})
	headerList := []headerKeyValue{}
	for k, v := range headers {
		val := v.(string)
		headerList = append(headerList, headerKeyValue{Key: k, Value: val})
	}
	return headerList
}

func transformTopicsToContentfulFormat(topicsTerraform interface{}) []string {
	topicList := []string{}
	for _, v := range topicsTerraform.([]interface{}) {
		topicList = append(topicList, v.(string))
	}
	return topicList
}
