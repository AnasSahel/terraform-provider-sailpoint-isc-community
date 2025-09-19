// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedcluster_test

import (
	"testing"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSailPointManagedClusterDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source: "hashicorp/random",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: acctest.TestAccSailPointManagedClusterDataSource(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.sailpoint_managed_cluster.test", "id"),
					resource.TestCheckResourceAttrSet("data.sailpoint_managed_cluster.test", "name"),
					resource.TestCheckResourceAttrSet("data.sailpoint_managed_cluster.test", "type"),
					resource.TestCheckResourceAttrSet("data.sailpoint_managed_cluster.test", "description"),
					resource.TestCheckResourceAttrSet("data.sailpoint_managed_cluster.test", "configuration.%"),
					resource.TestCheckResourceAttrSet("data.sailpoint_managed_cluster.test", "status"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
