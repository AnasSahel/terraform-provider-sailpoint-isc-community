package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOrderResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
				resource "sailpoint_transform" "test_acc" {
					name = "test-acc"
					type = "upper"
					attributes = jsonencode({
						value = "test-value"
					})
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sailpoint_transform.test_acc", "name", "test-acc"),
					resource.TestCheckResourceAttr("sailpoint_transform.test_acc", "type", "upper"),
					resource.TestCheckResourceAttr("sailpoint_transform.test_acc", "type", "upper"),
					resource.TestCheckResourceAttr("sailpoint_transform.test_acc", "internal", "false"),
					resource.TestCheckResourceAttrWith("sailpoint_transform.test_acc", "attributes", func(s string) error {
						if s != `{"value":"test-value"}` {
							return fmt.Errorf("expected attributes to be {\"value\":\"test-value\"}, got %s", s)
						}
						return nil
					}),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sailpoint_transform.test_acc",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: providerConfig + `
				resource "sailpoint_transform" "test_acc" {
					name = "test-acc"
					type = "upper"
					attributes = jsonencode({
						value = "new-test-value"
					})
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sailpoint_transform.test_acc", "name", "test-acc"),
					resource.TestCheckResourceAttr("sailpoint_transform.test_acc", "type", "upper"),
					resource.TestCheckResourceAttr("sailpoint_transform.test_acc", "type", "upper"),
					resource.TestCheckResourceAttr("sailpoint_transform.test_acc", "internal", "false"),
					resource.TestCheckResourceAttrWith("sailpoint_transform.test_acc", "attributes", func(s string) error {
						if s != `{"value":"new-test-value"}` {
							return fmt.Errorf("expected attributes to be {\"value\":\"new-test-value\"}, got %s", s)
						}
						return nil
					}),
				),
			},
		},
	})
}
