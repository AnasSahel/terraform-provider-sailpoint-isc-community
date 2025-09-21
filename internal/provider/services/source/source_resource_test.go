// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package source_test

import (
	"testing"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSailPointSourceResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source: "hashicorp/random",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccSailPointSourceResourceCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("sailpoint_source.test", "id"),
					resource.TestCheckResourceAttr("sailpoint_source.test", "name", "tf-test-delimited-source"),
					resource.TestCheckResourceAttr("sailpoint_source.test", "description", "Test delimited source created by Terraform"),
					resource.TestCheckResourceAttr("sailpoint_source.test", "connector", "delimited-file"),
					resource.TestCheckResourceAttrSet("sailpoint_source.test", "owner.id"),
					resource.TestCheckResourceAttrSet("sailpoint_source.test", "configuration.%"),
					resource.TestCheckResourceAttr("sailpoint_source.test", "features", "[]"),
					resource.TestCheckResourceAttr("sailpoint_source.test", "schemas", "[]"),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				ResourceName:      "sailpoint_source.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"management_workgroup",
					"schemas",
					"features",
				},
			},
			{
				Config: testAccSailPointSourceResourceUpdate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("sailpoint_source.test", "id"),
					resource.TestCheckResourceAttr("sailpoint_source.test", "description", "Updated test delimited source"),
					resource.TestCheckResourceAttr("sailpoint_source.test", "configuration.%", "2"),
				),
			},
		},
	})
}

func testAccSailPointSourceResourceCreate() string {
	return acctest.TestAccSailPointSourceResourceCreate()
}

func testAccSailPointSourceResourceUpdate() string {
	return acctest.TestAccSailPointSourceResourceUpdate()
}
