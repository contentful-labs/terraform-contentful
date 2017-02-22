# Contentful Terraform Provider

# Overview
Supports:
- [x] Spaces
- [ ] Content Types
- [x] API Keys
- [x] Webhooks
- [ ] Locales (not tested)

# TODO
- [ ] Add proper Terraform tests using `github.com/hashicorp/terraform/helper/resource`.
  E.g https://github.com/hashicorp/terraform/blob/master/builtin/providers/bitbucket/resource_hook_test.go#L12-L45
- [ ] Possibly serialize the JSON so that Terraform can notice the entire diff in read operations.
  Current we only compare stored variables.
- [ ] Perhaps write a proper Go SDK for Contentful

# Testing
We have unit tests mocking the Contentful HTTP API using `net/http/httptest` which asserts that the
request headers/payloads we send are valid according to Contentful.

    go test

# Using the plugin
Build the binary

    $ go build -o terraform-provider-contentful

Add it to your ~/.terraformrc

    $ cat ~/.terraformrc
    providers {
        contentful = "/home/dani/repos/go_pkg/src/github.com/danihodovic/contentful-terraform/terraform-provider-contentful"
    }

Use the provider

    provider "contentful" {
      cma_token = ""
      organization_id = ""
    }
    resource "contentful_space" "test" {
      name = "my-update-space-name"
    }
