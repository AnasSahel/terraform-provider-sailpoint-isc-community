// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package connector_test

import (
	"testing"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSailPointConnectorsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSailPointConnectorsDataSourceBasic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Check that the data source returns results
					acctest.CheckListNotEmpty("data.sailpoint_connectors.test", "connectors"),
					// Check that first connector has required attributes
					resource.TestCheckResourceAttrSet("data.sailpoint_connectors.test", "connectors.0.id"),
					resource.TestCheckResourceAttrSet("data.sailpoint_connectors.test", "connectors.0.name"),
					resource.TestCheckResourceAttrSet("data.sailpoint_connectors.test", "connectors.0.type"),
					resource.TestCheckResourceAttrSet("data.sailpoint_connectors.test", "connectors.0.script_name"),
				),
			},
		},
	})
}

func TestAccSailPointConnectorsDataSource_withFilters(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSailPointConnectorsDataSourceWithFilters(),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Check that the data source returns results
					resource.TestCheckResourceAttrSet("data.sailpoint_connectors.filtered", "id"),
					// Check that connectors list exists (may be empty due to filters)
					resource.TestCheckResourceAttr("data.sailpoint_connectors.filtered", "filters", `name sw "Active"`),
				),
			},
		},
	})
}

func TestAccSailPointConnectorsDataSource_withPagination(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSailPointConnectorsDataSourceWithPagination(),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Check that pagination parameters are set
					resource.TestCheckResourceAttr("data.sailpoint_connectors.paginated", "limit", "10"),
					resource.TestCheckResourceAttr("data.sailpoint_connectors.paginated", "offset", "0"),
					resource.TestCheckResourceAttr("data.sailpoint_connectors.paginated", "include_count", "true"),
				),
			},
		},
	})
}

func TestAccSailPointConnectorsDataSource_withLocale(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSailPointConnectorsDataSourceWithLocale(),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Check that locale parameter is set
					resource.TestCheckResourceAttr("data.sailpoint_connectors.localized", "locale", "de"),
					resource.TestCheckResourceAttrSet("data.sailpoint_connectors.localized", "id"),
				),
			},
		},
	})
}

func testAccSailPointConnectorsDataSourceBasic() string {
	return acctest.ProviderConfig + `
		data "sailpoint_connectors" "test" {}
	`
}

func testAccSailPointConnectorsDataSourceWithFilters() string {
	return acctest.ProviderConfig + `
		data "sailpoint_connectors" "filtered" {
			filters = "name sw \"Active\""
		}
	`
}

func testAccSailPointConnectorsDataSourceWithPagination() string {
	return acctest.ProviderConfig + `
		data "sailpoint_connectors" "paginated" {
			limit         = 10
			offset        = 0
			include_count = true
		}
	`
}

func testAccSailPointConnectorsDataSourceWithLocale() string {
	return acctest.ProviderConfig + `
		data "sailpoint_connectors" "localized" {
			locale = "de"
		}
	`
}

func TestAccSailPointConnectorsDataSource_withPaginateAll(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSailPointConnectorsDataSourceWithPaginateAll(),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Check that pagination parameters are set
					resource.TestCheckResourceAttr("data.sailpoint_connectors.all_paginated", "paginate_all", "true"),
					resource.TestCheckResourceAttrSet("data.sailpoint_connectors.all_paginated", "id"),
					// Should return connectors
					acctest.CheckListNotEmpty("data.sailpoint_connectors.all_paginated", "connectors"),
				),
			},
		},
	})
}

func TestAccSailPointConnectorsDataSource_withCustomPagination(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSailPointConnectorsDataSourceWithCustomPagination(),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Check that custom pagination parameters are set
					resource.TestCheckResourceAttr("data.sailpoint_connectors.custom_paginated", "paginate_all", "true"),
					resource.TestCheckResourceAttr("data.sailpoint_connectors.custom_paginated", "max_results", "5000"),
					resource.TestCheckResourceAttr("data.sailpoint_connectors.custom_paginated", "page_size", "100"),
					resource.TestCheckResourceAttrSet("data.sailpoint_connectors.custom_paginated", "id"),
				),
			},
		},
	})
}

func testAccSailPointConnectorsDataSourceWithPaginateAll() string {
	return acctest.ProviderConfig + `
		data "sailpoint_connectors" "all_paginated" {
			paginate_all = true
		}
	`
}

func testAccSailPointConnectorsDataSourceWithCustomPagination() string {
	return acctest.ProviderConfig + `
		data "sailpoint_connectors" "custom_paginated" {
			paginate_all = true
			max_results  = 5000
			page_size    = 100
		}
	`
}
