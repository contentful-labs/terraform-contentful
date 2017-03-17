package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	contentful "github.com/tolgaakyuz/contentful.go"
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

func resourceCreateWebhook(d *schema.ResourceData, m interface{}) (err error) {
	configMap := m.(map[string]interface{})
	client := configMap["client"].(*contentful.Contentful)

	space, err := client.GetSpace(d.Get("space_id").(string))
	if err != nil {
		return err
	}

	webhook := space.NewWebhook()
	webhook.Name = d.Get("name").(string)
	webhook.URL = d.Get("url").(string)
	webhook.Topics = transformTopicsToContentfulFormat(d.Get("topics").([]interface{}))
	webhook.Headers = transformHeadersToContentfulFormat(d.Get("headers"))
	webhook.HTTPBasicUsername = d.Get("http_basic_auth_username").(string)
	webhook.HTTPBasicPassword = d.Get("http_basic_auth_password").(string)

	err = webhook.Save()
	if err != nil {
		return err
	}

	err = setWebhookProperties(d, webhook)
	if err != nil {
		return err
	}

	d.SetId(webhook.Sys.ID)

	return nil
}

func resourceUpdateWebhook(d *schema.ResourceData, m interface{}) (err error) {
	configMap := m.(map[string]interface{})
	client := configMap["client"].(*contentful.Contentful)

	space, err := client.GetSpace(d.Get("space_id").(string))
	if err != nil {
		return err
	}

	webhook, err := space.GetWebhook(d.Id())
	if err != nil {
		return err
	}

	webhook.Name = d.Get("name").(string)
	webhook.URL = d.Get("url").(string)
	webhook.Topics = transformTopicsToContentfulFormat(d.Get("topics").([]interface{}))
	webhook.Headers = transformHeadersToContentfulFormat(d.Get("headers"))
	webhook.HTTPBasicUsername = d.Get("http_basic_auth_username").(string)
	webhook.HTTPBasicPassword = d.Get("http_basic_auth_password").(string)

	err = webhook.Save()
	if err != nil {
		return err
	}

	err = setWebhookProperties(d, webhook)
	if err != nil {
		return err
	}

	d.SetId(webhook.Sys.ID)

	return nil
}

func resourceReadWebhook(d *schema.ResourceData, m interface{}) error {
	configMap := m.(map[string]interface{})
	client := configMap["client"].(*contentful.Contentful)

	space, err := client.GetSpace(d.Get("space_id").(string))
	if err != nil {
		return err
	}

	webhook, err := space.GetWebhook(d.Id())
	if _, ok := err.(contentful.NotFoundError); ok {
		d.SetId("")
		return nil
	}

	if err != nil {
		return err
	}

	return setWebhookProperties(d, webhook)
}

func resourceDeleteWebhook(d *schema.ResourceData, m interface{}) (err error) {
	configMap := m.(map[string]interface{})
	client := configMap["client"].(*contentful.Contentful)

	space, err := client.GetSpace(d.Get("space_id").(string))
	if err != nil {
		return err
	}

	webhook, err := space.GetWebhook(d.Id())
	if err != nil {
		return err
	}

	err = webhook.Delete()
	if _, ok := err.(contentful.NotFoundError); ok {
		return nil
	}

	return err
}

func setWebhookProperties(d *schema.ResourceData, webhook *contentful.Webhook) (err error) {
	headers := make(map[string]string)
	for _, entry := range webhook.Headers {
		headers[entry.Key] = entry.Value
	}

	err = d.Set("headers", headers)
	if err != nil {
		return err
	}

	err = d.Set("space_id", webhook.Sys.Space.Sys.ID)
	if err != nil {
		return err
	}

	err = d.Set("version", webhook.Sys.Version)
	if err != nil {
		return err
	}

	err = d.Set("name", webhook.Name)
	if err != nil {
		return err
	}

	err = d.Set("url", webhook.URL)
	if err != nil {
		return err
	}

	err = d.Set("http_basic_auth_username", webhook.HTTPBasicUsername)
	if err != nil {
		return err
	}

	err = d.Set("topics", webhook.Topics)
	if err != nil {
		return err
	}

	return nil
}

func transformHeadersToContentfulFormat(headersTerraform interface{}) []*contentful.WebhookHeader {
	headers := []*contentful.WebhookHeader{}

	for k, v := range headersTerraform.(map[string]interface{}) {
		headers = append(headers, &contentful.WebhookHeader{
			Key:   k,
			Value: v.(string),
		})
	}

	return headers
}

func transformTopicsToContentfulFormat(topicsTerraform []interface{}) []string {
	var topics []string

	for _, v := range topicsTerraform {
		topics = append(topics, v.(string))
	}

	return topics
}
