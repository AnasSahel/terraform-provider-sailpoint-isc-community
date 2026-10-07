// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

// Acceptance tests for looking up and adopting a launcher's generated
// entitlement. They require TF_ACC=1 and a live SailPoint test tenant; they
// are skipped otherwise.
//
// Required environment variables:
//
//	TF_ACC=1
//	SAILPOINT_BASE_URL=https://<test-tenant>.api.identitynow.com
//	SAILPOINT_CLIENT_ID=<client-id>
//	SAILPOINT_CLIENT_SECRET=<client-secret>
//	SAILPOINT_TEST_LAUNCHER_NAME=<name of an existing INTERACTIVE_PROCESS launcher>
//
// ISC can take up to an hour to generate a launcher's entitlement, so the
// test points at a launcher that already exists rather than creating one.
package entitlement_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider"
)

var providerFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"sailpoint": providerserver.NewProtocol6WithError(provider.New("test")()),
}

func TestAccEntitlement_adoptLauncherEntitlementByName(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Set TF_ACC=1 to run acceptance tests (requires SailPoint test tenant credentials)")
	}

	launcherName := os.Getenv("SAILPOINT_TEST_LAUNCHER_NAME")
	if launcherName == "" {
		t.Fatal("SAILPOINT_TEST_LAUNCHER_NAME must be set to the name of an existing INTERACTIVE_PROCESS launcher")
	}

	config := fmt.Sprintf(`
data "sailpoint_entitlement" "launcher" {
  name = %q
}

resource "sailpoint_entitlement" "launcher" {
  id = data.sailpoint_entitlement.launcher.id
}
`, launcherName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sailpoint_entitlement.launcher", "name", launcherName),
					resource.TestCheckResourceAttrSet("data.sailpoint_entitlement.launcher", "id"),
					resource.TestCheckResourceAttrSet("data.sailpoint_entitlement.launcher", "source_id"),
					resource.TestCheckResourceAttrPair(
						"sailpoint_entitlement.launcher", "id",
						"data.sailpoint_entitlement.launcher", "id",
					),
				),
			},
			// Converged: a second plan is empty.
			{
				Config:             config,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}
