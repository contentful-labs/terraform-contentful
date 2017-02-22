package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

// Provider does shit
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"cma_token": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CONTENTFUL_MANAGEMENT_TOKEN", nil),
				Description: "The Contentful Management API token",
			},
			"organization_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CONTENTFUL_ORGANIZATION_ID", nil),
				Description: "The organization ID",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"contentful_space":   resourceContentfulSpace(),
			"contentful_apikey":  resourceContentfulAPIKey(),
			"contentful_webhook": resourceContentfulWebhook(),
			"contentful_locale":  resourceContentfulLocale(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := map[string]string{
		"cma_token":       d.Get("cma_token").(string),
		"organization_id": d.Get("organization_id").(string),
	}
	return config, nil
}
