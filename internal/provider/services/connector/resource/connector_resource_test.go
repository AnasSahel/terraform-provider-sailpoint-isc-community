// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package connector_resource_test

import (
	"fmt"
	"testing"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSailPointConnectorResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccSailPointConnectorResourceConfig("TestConnector", "custom"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sailpoint_connector.test", "name", "TestConnector"),
					resource.TestCheckResourceAttr("sailpoint_connector.test", "type", "custom"),
					resource.TestCheckResourceAttrSet("sailpoint_connector.test", "id"),
					resource.TestCheckResourceAttrSet("sailpoint_connector.test", "script_name"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sailpoint_connector.test",
				ImportState:       true,
				ImportStateVerify: true,
				// Ignore computed fields that might not match exactly
				ImportStateVerifyIgnore: []string{
					"application_xml", "source_config_xml", "correlation_config_xml",
					"file_upload", "s3_location", "source_config", "translation_properties",
				},
			},
			// Update and Read testing
			{
				Config: testAccSailPointConnectorResourceConfig("TestConnectorUpdated", "custom-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sailpoint_connector.test", "name", "TestConnectorUpdated"),
					resource.TestCheckResourceAttr("sailpoint_connector.test", "type", "custom-updated"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccSailPointConnectorResource_withOptionalFields(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with all optional fields
			{
				Config: testAccSailPointConnectorResourceConfigComplete(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sailpoint_connector.test", "name", "TestCompleteConnector"),
					resource.TestCheckResourceAttr("sailpoint_connector.test", "type", "custom-complete"),
					resource.TestCheckResourceAttr("sailpoint_connector.test", "class_name", "sailpoint.connector.CustomAdapter"),
					resource.TestCheckResourceAttr("sailpoint_connector.test", "direct_connect", "true"),
					resource.TestCheckResourceAttrSet("sailpoint_connector.test", "id"),
					resource.TestCheckResourceAttrSet("sailpoint_connector.test", "script_name"),
				),
			},
		},
	})
}

func testAccSailPointConnectorResourceConfig(name, connectorType string) string {
	return fmt.Sprintf(`
%s

resource "sailpoint_connector" "test" {
  name = "%s"
  type = "%s"
}
`, acctest.ProviderConfig, name, connectorType)
}

func testAccSailPointConnectorResourceConfigComplete() string {
	return fmt.Sprintf(`
%s

resource "sailpoint_connector" "test" {
  name           = "TestCompleteConnector"
  type           = "custom-complete"
  class_name     = "sailpoint.connector.CustomAdapter"
  direct_connect = true
}
`, acctest.ProviderConfig)
}
