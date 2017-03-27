package main

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	contentful "github.com/tolgaakyuz/contentful-go"
)

func TestAccContentfulContentType_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContentfulContentTypeDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccContentfulContentTypeConfig,
				Check: resource.TestCheckResourceAttr(
					"contentful_contenttype.mycontenttype", "name", "Terraform"),
			},
			resource.TestStep{
				Config: testAccContentfulContentTypeUpdateConfig,
				Check: resource.TestCheckResourceAttr(
					"contentful_contenttype.mycontenttype", "name", "Terraform name change"),
			},
			resource.TestStep{
				Config: testAccContentfulContentTypeLinkConfig,
				Check: resource.TestCheckResourceAttr(
					"contentful_contenttype.linked", "name", "Terraform Links"),
			},
		},
	})
}

func testAccCheckContentfulContentTypeExists(n string, contentType *contentful.ContentType) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No content type ID is set")
		}

		spaceID := rs.Primary.Attributes["space_id"]
		if spaceID == "" {
			return fmt.Errorf("No space_id is set")
		}

		client := testAccProvider.Meta().(*contentful.Contentful)

		ct, err := client.ContentTypes.Get(spaceID, rs.Primary.ID)
		if err != nil {
			return err
		}

		*contentType = *ct

		return nil
	}
}

func testAccCheckContentfulContentTypeDestroy(s *terraform.State) (err error) {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "contentful_contenttype" {
			continue
		}

		spaceID := rs.Primary.Attributes["space_id"]
		if spaceID == "" {
			return fmt.Errorf("No space_id is set")
		}

		client := testAccProvider.Meta().(*contentful.Contentful)

		_, err := client.ContentTypes.Get(spaceID, rs.Primary.ID)
		if _, ok := err.(contentful.NotFoundError); ok {
			return nil
		}

		return fmt.Errorf("Content Type still exists with id: %s", rs.Primary.ID)
	}

	return nil
}

var testAccContentfulContentTypeConfig = `
resource "contentful_space" "myspace" {
  name = "Terraform Space"
}

resource "contentful_contenttype" "mycontenttype" {
  space_id = "${contentful_space.myspace.id}"
  depends_on = ["contentful_space.myspace"]

  name = "Terraform"
  description = "Terraform Content Type"
  display_field = "field1"

  field {
    id = "field1"
    name = "Field 1"
    type = "Text"
    required = true
  }

  field {
    id = "field2"
    name = "Field 2"
    type = "Integer"
    required = false
  }
}
`

var testAccContentfulContentTypeUpdateConfig = `
resource "contentful_space" "myspace" {
  name = "Terraform Space"
}

resource "contentful_contenttype" "mycontenttype" {
  space_id = "${contentful_space.myspace.id}"
  depends_on = ["contentful_space.myspace"]

  name = "Terraform name change"
  description = "Terraform Content Type"
  display_field = "field1"

  field {
    id = "field1"
    name = "New field name"
    type = "Text"
    required = true
  }
}
`
var testAccContentfulContentTypeLinkConfig = `
resource "contentful_space" "myspace" {
  name = "Terraform Space"
}

resource "contentful_contenttype" "mycontenttype" {
  space_id = "${contentful_space.myspace.id}"
  depends_on = ["contentful_space.myspace"]

  name = "Terraform name change"
  description = "Terraform Content Type"
  display_field = "field1"

  field {
    id = "field1"
    name = "New field name"
    type = "Text"
    required = true
  }
}

resource "contentful_contenttype" "linked" {
  space_id = "${contentful_space.myspace.id}"
  depends_on = ["contentful_space.myspace"]
  depends_on = ["contentful_contenttype.mycontenttype"]

  name = "Terraform Links"
  description = "Terraform Content Type with links"
  display_field = "image"

  field {
    id = "image"
    name = "Image"
    type = "Array"
	items {
		type = "Link"
		link_type = "Asset"
	}
    required = false
  }

  field {
    id = "ctlink"
    name = "CT Link"
    type = "Array"
	items {
		type = "Link"
		validations {
			linkContentType = ["${contentful_contenttype.mycontenttype.id}"]
		}
		link_type = "Entry"
	}
    required = false
  }

}

`
