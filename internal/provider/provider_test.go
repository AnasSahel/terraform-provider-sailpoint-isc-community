// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// This file is kept for future provider-level tests.
// Individual resource and data source tests are located in their respective service packages.

const (
	providerConfig = `
		provider "sailpoint" {}
	`
)

var (
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"sailpoint": providerserver.NewProtocol6WithError(New("test")()),
	}
)
