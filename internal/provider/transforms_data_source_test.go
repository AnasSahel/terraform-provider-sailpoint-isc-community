package provider

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCoffeesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "sailpoint_transforms" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of transforms returned
					// resource.TestCheckResourceAttr("data.sailpoint_transforms.test", "transforms.#", ">0"),
					resource.TestCheckResourceAttrWith("data.sailpoint_transforms.test", "transforms.#", func(value string) error {
						valueInt, err := strconv.Atoi(value)
						if err != nil {
							return err
						}
						if valueInt <= 0 {
							return fmt.Errorf("expected transforms.# to be > 0, got %d", valueInt)
						}
						return nil
					}),
					// // Verify the first coffee to ensure all attributes are set
					resource.TestCheckResourceAttrWith("data.sailpoint_transforms.test", "transforms.0.id", func(value string) error {
						if value == "" {
							return fmt.Errorf("expected transforms.0.id to be set, got empty string")
						}
						return nil
					}),
					resource.TestCheckResourceAttrWith("data.sailpoint_transforms.test", "transforms.0.internal", func(value string) error {
						if value != "true" && value != "false" {
							return fmt.Errorf("expected transforms.0.internal to be 'true' or 'false', got %s", value)
						}
						return nil
					}),
					resource.TestCheckResourceAttrWith("data.sailpoint_transforms.test", "transforms.0.name", func(value string) error {
						if value == "" {
							return fmt.Errorf("expected transforms.0.name to be set, got empty string")
						}
						return nil
					}),
					resource.TestCheckResourceAttrWith("data.sailpoint_transforms.test", "transforms.0.attributes", func(value string) error {
						if value == "" {
							return fmt.Errorf("expected transforms.0.attributes to be set, got empty string")
						}
						return nil
					}),
				),
			},
		},
	})
}
