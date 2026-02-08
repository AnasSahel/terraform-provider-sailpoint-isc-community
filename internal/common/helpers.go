// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"context"
	"fmt"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// ConfigureClient extracts the client from provider data and returns it.
// resourceType should be a descriptive name like "identity attribute resource" or "identity attribute data source".
func ConfigureClient(ctx context.Context, providerData any, resourceType string) (*client.Client, diag.Diagnostics) {
	var diagnostics diag.Diagnostics

	if providerData == nil {
		tflog.Debug(ctx, fmt.Sprintf("No provider data configured for %s", resourceType))
		return nil, diagnostics
	}

	c, ok := providerData.(*client.Client)
	if !ok {
		tflog.Debug(ctx, fmt.Sprintf("Provider data is of unexpected type for %s", resourceType))
		diagnostics.AddError(
			"Unexpected Configure Type",
			fmt.Sprintf("Expected *client.Client type for provider data but got: %T. Please report this issue to the provider developers.", providerData),
		)
		return nil, diagnostics
	}

	return c, diagnostics
}
