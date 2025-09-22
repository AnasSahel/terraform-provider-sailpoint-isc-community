// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource_test

import (
	"testing"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSailPointTransformResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.TestAccSailPointTransformResourceCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("sailpoint_transform.test", "id"),
					resource.TestCheckResourceAttr("sailpoint_transform.test", "name", "terraform-poc-001"),
					resource.TestCheckResourceAttr("sailpoint_transform.test", "type", "upper"),
					resource.TestCheckResourceAttr("sailpoint_transform.test", "internal", "false"),

					resource.TestCheckResourceAttrSet("sailpoint_transform.test", "attributes"),
				),
			},
			{
				ResourceName:      "sailpoint_transform.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: acctest.TestAccSailPointTransformResourceUpdate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("sailpoint_transform.test", "id"),
					resource.TestCheckResourceAttr("sailpoint_transform.test", "name", "terraform-poc-001"),
					resource.TestCheckResourceAttr("sailpoint_transform.test", "type", "upper"),
					resource.TestCheckResourceAttr("sailpoint_transform.test", "internal", "false"),

					resource.TestCheckResourceAttrSet("sailpoint_transform.test", "attributes"),
				),
			},
		},
	})
}
