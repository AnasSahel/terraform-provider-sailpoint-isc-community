// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOrderResource(t *testing.T) {
	randomSuffix := fmt.Sprintf("%d", rand.Int63())
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "sailpoint_transform" "test_acc_%s" {
					name = "test-acc-%s"
					type = "upper"
					attributes = jsonencode({
						value = "test-value"
					})
				}
				`, randomSuffix, randomSuffix),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("sailpoint_transform.test_acc_%s", randomSuffix), "name", fmt.Sprintf("test-acc-%s", randomSuffix)),
					resource.TestCheckResourceAttr(fmt.Sprintf("sailpoint_transform.test_acc_%s", randomSuffix), "type", "upper"),
					resource.TestCheckResourceAttr(fmt.Sprintf("sailpoint_transform.test_acc_%s", randomSuffix), "internal", "false"),
					resource.TestCheckResourceAttrWith(fmt.Sprintf("sailpoint_transform.test_acc_%s", randomSuffix), "attributes", func(s string) error {
						if s != `{"value":"test-value"}` {
							return fmt.Errorf("expected attributes to be {\"value\":\"test-value\"}, got %s", s)
						}
						return nil
					}),
				),
			},
			// ImportState testing
			{
				ResourceName:      fmt.Sprintf("sailpoint_transform.test_acc_%s", randomSuffix),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "sailpoint_transform" "test_acc_%s" {
					name = "test-acc-%s"
					type = "upper"
					attributes = jsonencode({
						value = "new-test-value"
					})
				}
				`, randomSuffix, randomSuffix),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("sailpoint_transform.test_acc_%s", randomSuffix), "name", fmt.Sprintf("test-acc-%s", randomSuffix)),
					resource.TestCheckResourceAttr(fmt.Sprintf("sailpoint_transform.test_acc_%s", randomSuffix), "type", "upper"),
					resource.TestCheckResourceAttr(fmt.Sprintf("sailpoint_transform.test_acc_%s", randomSuffix), "internal", "false"),
					resource.TestCheckResourceAttrWith(fmt.Sprintf("sailpoint_transform.test_acc_%s", randomSuffix), "attributes", func(s string) error {
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
