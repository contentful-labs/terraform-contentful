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
			{
				Config: testAccCheckContentfulSpaceConfig,
				Check: resource.TestCheckResourceAttr(
					"contentful_space.myspace", "name", "terraform test"),
			},
		},
	})
}

func testAccCheckContentfulSpaceDestroy(s *terraform.State) error {
	configMap := testAccProvider.Meta().(map[string]interface{})

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "contentful_space" {
			continue
		}

		client := configMap["client"].(*contentful.Contentful)
		_, err := client.GetSpace(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Space %s still exists after destroy", rs.Primary.ID)
		}

		switch t := err.(type) {
		case contentful.NotFoundError:
			return nil
		default:
			_ = t
			return fmt.Errorf("Error checking space %s: %s", rs.Primary.ID, err)
		}
	}

	return nil
}

var testAccCheckContentfulSpaceConfig = `
resource "contentful_space" "myspace" {
  name = "terraform test"
}

`
