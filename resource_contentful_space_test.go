package main

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
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
	configMap := testAccProvider.Meta().(map[string]string)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "contentful_space" {
			continue
		}

		exists, err := spaceExists(configMap["cma_token"], rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("Error checking space %s: %s", rs.Primary.ID, err)
		}

		if exists {
			return fmt.Errorf("Space %s still exists after destroy", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckContentfulSpaceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No ID is set")
		}

		configMap := testAccProvider.Meta().(map[string]string)
		_, err := spaceExists(configMap["cma_token"], rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("Error checking space %s: %s", rs.Primary.ID, err)

		}
		return nil
	}
}

var testAccCheckContentfulSpaceConfig = `
resource "contentful_space" "myspace" {
  name = "terraform test"
}

`
