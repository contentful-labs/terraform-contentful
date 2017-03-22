[![Build Status](https://img.shields.io/circleci/project/github/contentful-labs/terraform-contentful.svg?branch=master)](https://circleci.com/gh/contentful-labs/terraform-contentful)
[![license](https://img.shields.io/github/license/contentful-labs/terraform-contentful.svg)](https://github.com/contentful-labs/terraform-contentful/blob/master/LICENSE)

Terraform Provider for [Contentful's](https://www.contentful.com) Content Management API

# About

[Contentful](https://www.contentful.com) is a content management platform for web applications, mobile apps and connected devices. It allows you to create, edit & manage content in the cloud and publish it anywhere via a powerful API. Contentful offers tools for managing editorial teams and enabling cooperation between organizations.

[Terraform](https://www.terraform.io) is a tool for building, changing, and versioning infrastructure safely and efficiently. Terraform can manage existing and popular service providers as well as custom in-house solutions.

# Features

Create, update and delete Contentful resources such as:
- [x] Spaces
- [ ] Content Types
- [x] API Keys
- [x] Webhooks
- [ ] Locales

# Getting started

Download [go](https://golang.org/dl) for your platform.

Follow the [Install the Go tools](https://golang.org/doc/install#install) instructions.

Download [terraform](https://www.terraform.io/downloads.html) for your platform.

Follow the [Installing Terraform](https://www.terraform.io/intro/getting-started/install.html) instructions.

Create a directory where your terraform files and states will be placed. Although not mandatory this should be placed under a version control software such as [git](https://git-scm.com).

Make sure you have your Content Management API Token and the organization ID before starting. As an alternative to configuring the provider in the terraform file you can also set environment variables.

```sh
    # For Linux/Mac OS
    export CONTENTFUL_ORGANIZATION_ID=<your organization ID>
    export CONTENTFUL_MANAGEMENT_TOKEN=<your CMA Token>
```

```
    REM For Windows
    setx CONTENTFUL_ORGANIZATION_ID "<your organization ID>"
    setx CONTENTFUL_MANAGEMENT_TOKEN "<your CMA Token>"
```

# Using the provider
Build the binary

    $ go build -o terraform-provider-contentful

Add it to your ~/.terraformrc (or %APPDATA%/terraform.rc for Windows)

    $ cat ~/.terraformrc
    providers {
        contentful = "<path to the go binary>/terraform-provider-contentful"
    }

Use the provider by creating a main.tf file with:

    provider "contentful" {
      cma_token = "<your CMA Token>"
      organization_id = "<your organization ID>"
    }
    resource "contentful_space" "test" {
      name = "my-update-space-name"
    }

Run the terraform plan

    terraform plan -out=contentful.plan

Check the changes
```
The refreshed state will be used to calculate this plan, but will not be
persisted to local or remote state storage.

The Terraform execution plan has been generated and is shown below.
Resources are shown in alphabetical order for quick scanning. Green resources
will be created (or destroyed and then created if an existing resource
exists), yellow resources are being changed in-place, and red resources
will be destroyed. Cyan entries are data sources to be read.

Your plan was also saved to the path below. Call the "apply" subcommand
with this plan file and Terraform will exactly execute this execution
plan.

Path: contentful.plan

+ contentful_space.test
    default_locale: "en"
    name:           "my-update-space-name"
    version:        "<computed>"


Plan: 1 to add, 0 to change, 0 to destroy.
```

Apply the plan
```
contentful_space.test: Creating...
  default_locale: "" => "en"
  name:           "" => "my-update-space-name"
  version:        "" => "<computed>"
contentful_space.test: Creation complete (ID: yculypygam9h)

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

The state of your infrastructure has been saved to the path
below. This state is required to modify and destroy your
infrastructure, so keep it safe. To inspect the complete state
use the `terraform show` command.

State path:
```

## Testing

    TF_ACC=1 go test -v

To enable higher verbose mode

    TF_LOG=debug TF_ACC=1 go test -v

## Documentation/References

### Hashicorp
[Terraform plugins](https://www.terraform.io/docs/plugins/basics.html)

[Writing custom terraform providers](https://www.hashicorp.com/blog/writing-custom-terraform-providers)

### Other references
Julien Fabre: [Writing a Terraform provider](http://blog.jfabre.net/2017/01/22/writing-terraform-provider)

## Support

If you have a problem with this provider, please file an [issue](https://github.com/contentful-labs/terraform-contentful/issues/new) here on Github.

## License

MIT
