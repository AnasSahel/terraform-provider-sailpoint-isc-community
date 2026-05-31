// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

// Package source_test holds the acceptance tests for the sailpoint_source
// resource. Acceptance tests require TF_ACC=1 and a live SailPoint test
// tenant; they are skipped otherwise.
//
// Required environment variables:
//
//	TF_ACC=1
//	SAILPOINT_BASE_URL=https://<test-tenant>.api.identitynow.com
//	SAILPOINT_CLIENT_ID=<client-id>
//	SAILPOINT_CLIENT_SECRET=<client-secret>
//	SAILPOINT_TEST_SOURCE_OWNER_ID=<identity-id in the test tenant>
//
// No real tenant names, client names, or internal IDs appear in this file.
package source_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider"
)

// providerFactories wires the in-process provider into the acceptance test
// harness.
var providerFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"sailpoint": providerserver.NewProtocol6WithError(provider.New("test")()),
}

// checkConnectorAttributesProjection is a custom resource.CheckFunc that
// verifies the partial-management contract:
//
//  1. connector_attributes contains exactly the declared keys — no server extras.
//  2. connector_attributes_all contains at least those keys (server may add more).
//  3. Values in connector_attributes match connector_attributes_all for each
//     declared key, confirming the projection re-read server values rather than
//     blindly copying prior state.
func checkConnectorAttributesProjection(resourceName string, declaredKeys []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %q not found in state", resourceName)
		}

		managedRaw, ok := rs.Primary.Attributes["connector_attributes"]
		if !ok {
			return fmt.Errorf("connector_attributes not found in state")
		}
		allRaw, ok := rs.Primary.Attributes["connector_attributes_all"]
		if !ok {
			return fmt.Errorf("connector_attributes_all not found in state")
		}

		var managed, all map[string]interface{}
		if err := json.Unmarshal([]byte(managedRaw), &managed); err != nil {
			return fmt.Errorf("connector_attributes is not valid JSON: %w", err)
		}
		if err := json.Unmarshal([]byte(allRaw), &all); err != nil {
			return fmt.Errorf("connector_attributes_all is not valid JSON: %w", err)
		}

		// connector_attributes must contain exactly the declared keys.
		for _, k := range declaredKeys {
			if _, exists := managed[k]; !exists {
				return fmt.Errorf("declared key %q missing from connector_attributes", k)
			}
		}
		if len(managed) != len(declaredKeys) {
			return fmt.Errorf("connector_attributes has %d keys, want exactly %d (declared: %v, got keys: %v)",
				len(managed), len(declaredKeys), declaredKeys, mapKeys(managed))
		}

		// connector_attributes_all must be a superset.
		for k, mv := range managed {
			av, exists := all[k]
			if !exists {
				return fmt.Errorf("key %q present in connector_attributes but missing from connector_attributes_all", k)
			}
			mvJSON, _ := json.Marshal(mv)
			avJSON, _ := json.Marshal(av)
			if string(mvJSON) != string(avJSON) {
				return fmt.Errorf("key %q: connector_attributes=%s does not match connector_attributes_all=%s — "+
					"Read projection must use server values, not stale state", k, mvJSON, avJSON)
			}
		}

		return nil
	}
}

func mapKeys(m map[string]interface{}) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

// testAccSourceConfig returns a minimal Terraform config declaring a CSV
// (delimited-file) source with only two connector_attributes keys. The owner
// ID is injected via the SAILPOINT_TEST_SOURCE_OWNER_ID environment variable
// and never hardcoded. The source name is a stable invented string with no
// client or tenant data.
func testAccSourceConfig(name, ownerID string) string {
	return fmt.Sprintf(`
resource "sailpoint_source" "test" {
  name      = %[1]q
  connector = "delimited-file"

  owner = {
    type = "IDENTITY"
    id   = %[2]q
  }

  # Only two keys declared — the server will inject additional keys
  # (status, healthy, since, cloudEnabled, etc.) into connectorAttributes.
  # The provider must NOT surface those extra keys as drift.
  connector_attributes = jsonencode({
    cloudDisplayName = %[1]q
    hasHeader        = true
  })
}
`, name, ownerID)
}

// testAccSourceConfigUpdated changes cloudDisplayName to exercise the
// drift-detection path for declared keys after an update.
func testAccSourceConfigUpdated(name, ownerID string) string {
	return fmt.Sprintf(`
resource "sailpoint_source" "test" {
  name      = %[1]q
  connector = "delimited-file"

  owner = {
    type = "IDENTITY"
    id   = %[2]q
  }

  connector_attributes = jsonencode({
    cloudDisplayName = "%[1]s-updated"
    hasHeader        = true
  })
}
`, name, ownerID)
}

// TestAccSource_connectorAttributesPartialManagement exercises the full
// partial-management lifecycle:
//
//  1. Create: connector_attributes in state has exactly the two declared keys;
//     connector_attributes_all may have more (server-injected).
//  2. No-op plan (refresh-only): zero changes — server keys must not surface
//     as drift.
//  3. Update a declared key: plan shows exactly that one change.
//  4. Import round-trip: imported state matches pre-import state immediately,
//     without requiring an intermediate apply (ModifyPlan projects to config
//     keys before the diff runs).
func TestAccSource_connectorAttributesPartialManagement(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Set TF_ACC=1 to run acceptance tests (requires SailPoint test tenant credentials)")
	}

	ownerID := os.Getenv("SAILPOINT_TEST_SOURCE_OWNER_ID")
	if ownerID == "" {
		t.Fatal("SAILPOINT_TEST_SOURCE_OWNER_ID must be set to a valid identity ID in the test tenant")
	}

	// Invented name — no client or tenant data.
	sourceName := "tf-acc-partial-mgmt-test"
	resourceName := "sailpoint_source.test"
	declaredKeys := []string{"cloudDisplayName", "hasHeader"}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			// Step 1: create with a two-key subset of connector_attributes.
			{
				Config: testAccSourceConfig(sourceName, ownerID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", sourceName),
					// connector_attributes_all must be set (server returns extras).
					resource.TestCheckResourceAttrSet(resourceName, "connector_attributes_all"),
					// Core assertion: connector_attributes = exactly declared keys,
					// values matching server (not stale config).
					checkConnectorAttributesProjection(resourceName, declaredKeys),
				),
			},
			// Step 2: plan-only refresh — must show zero changes.
			// Verifies that the Read projection suppresses server-injected keys.
			{
				Config:             testAccSourceConfig(sourceName, ownerID),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			// Step 3: update one declared key.
			{
				Config: testAccSourceConfigUpdated(sourceName, ownerID),
				Check: resource.ComposeTestCheckFunc(
					checkConnectorAttributesProjection(resourceName, declaredKeys),
				),
			},
			// Step 4: import — ImportStateVerify compares the imported state
			// against the pre-import state. The two must match, meaning that
			// ModifyPlan projected connector_attributes to the config's declared
			// keys immediately (zero-diff without an intermediate apply).
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// provision_as_csv is create-only and absent from the API response.
				ImportStateVerifyIgnore: []string{"provision_as_csv"},
			},
		},
	})
}

// TestAccSource_importZeroDiff exercises the import case in isolation:
// import an existing source, then verify that a plan against a config
// declaring only a subset of keys produces zero changes — no apply needed.
func TestAccSource_importZeroDiff(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Set TF_ACC=1 to run acceptance tests (requires SailPoint test tenant credentials)")
	}

	ownerID := os.Getenv("SAILPOINT_TEST_SOURCE_OWNER_ID")
	if ownerID == "" {
		t.Fatal("SAILPOINT_TEST_SOURCE_OWNER_ID must be set to a valid identity ID in the test tenant")
	}

	sourceName := "tf-acc-import-zero-diff-test"
	resourceName := "sailpoint_source.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			// Create the source first.
			{
				Config: testAccSourceConfig(sourceName, ownerID),
			},
			// Import and verify immediately — no intermediate apply needed.
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"provision_as_csv"},
			},
			// Belt-and-suspenders: explicit no-op plan after import.
			{
				Config:             testAccSourceConfig(sourceName, ownerID),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}
