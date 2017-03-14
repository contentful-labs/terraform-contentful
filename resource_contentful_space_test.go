package main

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	contentful "github.com/tolgaakyuz/contentful.go"
)

func TestAccContentfulSpace_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContentfulSpaceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccContentfulSpaceConfig,
				Check: resource.TestCheckResourceAttr(
					"contentful_space.myspace", "name", "space-name"),
			},
			resource.TestStep{
				Config: testAccContentfulSpaceUpdateConfig,
				Check: resource.TestCheckResourceAttr(
					"contentful_space.myspace", "name", "changed-space-name"),
			},
		},
	})
}

func testAccCheckContentfulSpaceDestroy(s *terraform.State) error {
	configMap := testAccProvider.Meta().(map[string]interface{})
	client := configMap["client"].(*contentful.Contentful)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "contentful_space" {
			continue
		}

		space, err := client.GetSpace(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Space %s still exists after destroy", space.ID())
		}
	}

	return nil
}

var testAccContentfulSpaceConfig = `
resource "contentful_space" "myspace" {
  name = "space-name"
}
`

var testAccContentfulSpaceUpdateConfig = `
resource "contentful_space" "myspace" {
  name = "changed-space-name"
}
`
