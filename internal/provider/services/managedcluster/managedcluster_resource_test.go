// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedcluster_test

import (
	"testing"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSailPointManagedClusterResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source: "hashicorp/random",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: acctest.TestAccSailPointManagedClusterResourceCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("sailpoint_managed_cluster.test", "id"),
					resource.TestCheckResourceAttrSet("sailpoint_managed_cluster.test", "name"),
					resource.TestCheckResourceAttr("sailpoint_managed_cluster.test", "type", "idn"),
					resource.TestCheckResourceAttr("sailpoint_managed_cluster.test", "description", "Test managed cluster created by Terraform"),
					resource.TestCheckResourceAttrSet("sailpoint_managed_cluster.test", "configuration.%"),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				ResourceName:      "sailpoint_managed_cluster.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"alert_key",
					"ccg_version",
					"client_type",
					"status",
					"configuration.%",
					"configuration.cluster_external_id",
					"configuration.cluster_type",
					"configuration.restart_threshold_in_hours",
				},
			},
			{
				Config: acctest.TestAccSailPointManagedClusterResourceUpdate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("sailpoint_managed_cluster.test", "id"),
					resource.TestCheckResourceAttrSet("sailpoint_managed_cluster.test", "name"),
					resource.TestCheckResourceAttr("sailpoint_managed_cluster.test", "type", "idn"),
					resource.TestCheckResourceAttr("sailpoint_managed_cluster.test", "description", "Updated test managed cluster description"),
					resource.TestCheckResourceAttrSet("sailpoint_managed_cluster.test", "configuration.%"),
				),
			},
		},
	})
}
