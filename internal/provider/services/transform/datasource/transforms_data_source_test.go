// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datasource_test

import (
	"testing"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSailPointTransformsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.TestAccSailPointTransformsDataSource(),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.CheckListNotEmpty("data.sailpoint_transforms.test", "transforms"),
					resource.TestCheckResourceAttrSet("data.sailpoint_transforms.test", "transforms.0.id"),
					resource.TestCheckResourceAttrSet("data.sailpoint_transforms.test", "transforms.0.name"),
					resource.TestCheckResourceAttrSet("data.sailpoint_transforms.test", "transforms.0.type"),
					resource.TestCheckResourceAttrSet("data.sailpoint_transforms.test", "transforms.0.internal"),
					resource.TestCheckResourceAttrSet("data.sailpoint_transforms.test", "transforms.0.attributes"),
				),
			},
		},
	})
}
