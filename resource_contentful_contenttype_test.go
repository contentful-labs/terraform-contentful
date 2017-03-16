package main

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	contentful "github.com/tolgaakyuz/contentful.go"
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
					"contentful_space.myspace", "name", "changed-space-name"),
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

		configMap := testAccProvider.Meta().(map[string]interface{})
		client := configMap["client"].(*contentful.Contentful)

		space, err := client.GetSpace(spaceID)
		if err != nil {
			return fmt.Errorf("No space with this id: %s", rs.Primary.Attributes["space_id"])
		}

		ct, err := space.GetContentType(rs.Primary.ID)
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

		configMap := testAccProvider.Meta().(map[string]interface{})
		client := configMap["client"].(*contentful.Contentful)

		space, err := client.GetSpace(spaceID)

		if err != nil {
			if _, ok := err.(contentful.NotFoundError); ok {
				return nil
			}
			return fmt.Errorf("Error checking space_id: %s", spaceID)
		}

		_, err = space.GetContentType(rs.Primary.ID)
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
  displayField = "field1"

  field {
    id = "field1"
    name = "Field 1"
    type = "Text"
    required = true
  }

  field {
    id = "field2"
    name = "Field 2"
    type = "Number"
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
  displayField = "field1"

  field {
    id = "field1"
    name = "Field 1"
    type = "Text"
    required = true
  }

  field {
    id = "field2"
    name = "Field 2"
    type = "Number"
    required = false
  }

}
`
