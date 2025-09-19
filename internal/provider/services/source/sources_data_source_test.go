// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package source_test

import (
	"testing"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSailPointSourcesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSailPointSourcesDataSource(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.sailpoint_sources.test", "sources.#"),
					resource.TestCheckResourceAttrSet("data.sailpoint_sources.test", "sources.0.id"),
					resource.TestCheckResourceAttrSet("data.sailpoint_sources.test", "sources.0.name"),
					resource.TestCheckResourceAttrSet("data.sailpoint_sources.test", "sources.0.connector"),
				),
			},
		},
	})
}

func TestAccSailPointSourceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source: "hashicorp/random",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccSailPointSourceDataSource(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.sailpoint_source.test", "id"),
					resource.TestCheckResourceAttrSet("data.sailpoint_source.test", "name"),
					resource.TestCheckResourceAttrSet("data.sailpoint_source.test", "connector"),
					resource.TestCheckResourceAttrSet("data.sailpoint_source.test", "owner.id"),
				),
			},
		},
	})
}

func testAccSailPointSourcesDataSource() string {
	return acctest.TestAccSailPointSourcesDataSource()
}

func testAccSailPointSourceDataSource() string {
	return acctest.TestAccSailPointSourceDataSource()
}
