package main

import (
	"os"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	contentful "github.com/tolgaakyuz/contentful-go"
)

// Provider does shit
func Provider() terraform.ResourceProvider {
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
			"contentful_space":       resourceContentfulSpace(),
			"contentful_contenttype": resourceContentfulContentType(),
			"contentful_apikey":      resourceContentfulAPIKey(),
			"contentful_webhook":     resourceContentfulWebhook(),
			"contentful_locale":      resourceContentfulLocale(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	c, err := contentful.New(&contentful.Settings{
		CMAToken: d.Get("cma_token").(string),
		BaseURL:  "https://api.flinkly.com",
	})
	if err != nil {
		return nil, err
	}

	if os.Getenv("TF_LOG") != "" {
		c.Debug = true
	}

	config := map[string]interface{}{
		"cma_token":       d.Get("cma_token").(string),
		"organization_id": d.Get("organization_id").(string),
		"client":          c,
	}

	return config, nil
}
