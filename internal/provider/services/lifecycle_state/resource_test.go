// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package lifecycle_state_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccSailPointLifecycleStateResource tests the lifecycle state resource.
// This test requires an existing identity profile ID to be provided
// via environment variable: TF_VAR_identity_profile_id.
// If this is not set, the test will be skipped.
func TestAccSailPointLifecycleStateResource(t *testing.T) {
	// Skip test if required environment variable is not set
	acctest.RequireEnvVar(t, "TF_VAR_identity_profile_id")

	identityProfileID := os.Getenv("TF_VAR_identity_profile_id")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: acctest.LifecycleStateResourceCreateConfigWithVars(),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify resource attributes
					resource.TestCheckResourceAttr("sailpoint_lifecycle_state.test", "identity_profile_id", identityProfileID),
					resource.TestCheckResourceAttr("sailpoint_lifecycle_state.test", "name", "Terraform Test State"),
					resource.TestCheckResourceAttr("sailpoint_lifecycle_state.test", "technical_name", "terraform_test_state"),
					resource.TestCheckResourceAttr("sailpoint_lifecycle_state.test", "description", "A lifecycle state created by Terraform for testing"),
					resource.TestCheckResourceAttr("sailpoint_lifecycle_state.test", "enabled", "true"),

					// Verify computed attributes are set
					resource.TestCheckResourceAttrSet("sailpoint_lifecycle_state.test", "id"),
					resource.TestCheckResourceAttrSet("sailpoint_lifecycle_state.test", "created"),
					resource.TestCheckResourceAttrSet("sailpoint_lifecycle_state.test", "modified"),
					resource.TestCheckResourceAttrSet("sailpoint_lifecycle_state.test", "identity_count"),
					resource.TestCheckResourceAttrSet("sailpoint_lifecycle_state.test", "priority"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sailpoint_lifecycle_state.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccLifecycleStateImportStateIdFunc,
				// Import verification may fail due to API response differences
				ImportStateVerifyIgnore: []string{"modified"},
			},
			// Update and Read testing
			{
				Config: acctest.LifecycleStateResourceUpdateConfigWithVars(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sailpoint_lifecycle_state.test", "identity_profile_id", identityProfileID),
					resource.TestCheckResourceAttr("sailpoint_lifecycle_state.test", "name", "Terraform Test State"),
					resource.TestCheckResourceAttr("sailpoint_lifecycle_state.test", "technical_name", "terraform_test_state"),
					resource.TestCheckResourceAttr("sailpoint_lifecycle_state.test", "description", "Updated description for Terraform testing"),
					resource.TestCheckResourceAttr("sailpoint_lifecycle_state.test", "enabled", "false"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

// testAccLifecycleStateImportStateIdFunc returns the import ID for the lifecycle state resource.
func testAccLifecycleStateImportStateIdFunc(s *terraform.State) (string, error) {
	rs, ok := s.RootModule().Resources["sailpoint_lifecycle_state.test"]
	if !ok {
		return "", fmt.Errorf("resource not found: sailpoint_lifecycle_state.test")
	}

	identityProfileID := rs.Primary.Attributes["identity_profile_id"]
	lifecycleStateID := rs.Primary.Attributes["id"]

	return fmt.Sprintf("%s:%s", identityProfileID, lifecycleStateID), nil
}
