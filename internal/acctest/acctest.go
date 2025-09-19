// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package acctest

import (
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// This file is kept for future provider-level tests.
// Individual resource and data source tests are located in their respective service packages.

const (
	ProviderConfig = `
		provider "sailpoint" {}
	`
)

var (
	TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"sailpoint": providerserver.NewProtocol6WithError(provider.New("test")()),
	}
)

func TestAccSailPointTransformsDataSource() string {
	return ProviderConfig + `
		data "sailpoint_transforms" "test" {}
	`
}

func TestAccSailPointTransformResourceCreate() string {
	return ProviderConfig + `
		resource "sailpoint_transform" "test" {
			name = "terraform-poc-001"
			type = "upper"
			attributes = jsonencode({
				"input" = "name"
			})
		}
	`
}

func TestAccSailPointTransformResourceUpdate() string {
	return ProviderConfig + `
		resource "sailpoint_transform" "test" {
			name = "terraform-poc-001"
			type = "upper"
			attributes = jsonencode({
				"input" = "updated name"
			})
		}
	`
}
