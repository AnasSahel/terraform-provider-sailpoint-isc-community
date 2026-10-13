// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

// Acceptance tests for the approval config resource. They require TF_ACC=1 and a
// live SailPoint test tenant; they are skipped otherwise.
//
// Required environment variables:
//
//	TF_ACC=1
//	SAILPOINT_BASE_URL=https://<test-tenant>.api.identitynow.com
//	SAILPOINT_CLIENT_ID=<client-id>
//	SAILPOINT_CLIENT_SECRET=<client-secret>
//	SAILPOINT_TEST_OWNER_ID=<ID of an identity that can own a role>
package approval_config_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider"
)

var providerFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"sailpoint": providerserver.NewProtocol6WithError(provider.New("test")()),
}

func accConfig(ownerID, suffix string, maxReminders int) string {
	return fmt.Sprintf(`
resource "sailpoint_role" "check" {
  name  = "tf-acc-approval-config-%[2]s"
  owner = { type = "IDENTITY", id = %[1]q }
}

resource "sailpoint_approval_config" "check" {
  scope     = "ROLE"
  object_id = sailpoint_role.check.id

  reminder_config = {
    enabled                   = true
    days_until_first_reminder = 2
    reminder_cron_schedule    = "@every 48h"
    max_reminders             = %[3]d
  }

  escalation_config = {
    escalation_chain = [{ identity_type = "MANAGER_OF" }]
  }
}
`, ownerID, suffix, maxReminders)
}

func TestAccApprovalConfig_role(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Set TF_ACC=1 to run acceptance tests (requires SailPoint test tenant credentials)")
	}
	ownerID := os.Getenv("SAILPOINT_TEST_OWNER_ID")
	if ownerID == "" {
		t.Skip("Set SAILPOINT_TEST_OWNER_ID to the ID of an identity that can own a role")
	}
	suffix := acctest.RandString(6)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: accConfig(ownerID, suffix, 3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith("sailpoint_approval_config.check", "id", func(v string) error {
						if !strings.HasPrefix(v, "ROLE/") {
							return fmt.Errorf("id %q does not start with ROLE/", v)
						}
						return nil
					}),
					resource.TestCheckResourceAttr("sailpoint_approval_config.check", "scope", "ROLE"),
					resource.TestCheckResourceAttr("sailpoint_approval_config.check", "reminder_config.max_reminders", "3"),
					resource.TestCheckResourceAttr("sailpoint_approval_config.check", "escalation_config.escalation_chain.0.identity_type", "MANAGER_OF"),
					resource.TestCheckResourceAttrSet("sailpoint_approval_config.check", "timeout_config.days_until_timeout"),
				),
			},
			// Converged: a second plan is empty.
			{
				Config:   accConfig(ownerID, suffix, 3),
				PlanOnly: true,
			},
			{
				Config: accConfig(ownerID, suffix, 5),
				Check:  resource.TestCheckResourceAttr("sailpoint_approval_config.check", "reminder_config.max_reminders", "5"),
			},
			{
				Config:            accConfig(ownerID, suffix, 5),
				ResourceName:      "sailpoint_approval_config.check",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["sailpoint_role.check"]
					if !ok {
						return "", fmt.Errorf("sailpoint_role.check not found in state")
					}
					return "ROLE/" + rs.Primary.ID, nil
				},
			},
		},
	})
}
