// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managed_cluster

import "github.com/hashicorp/terraform-plugin-framework/types"

type ManagedClusterQueue struct {
	Name   types.String `tfsdk:"name"`
	Region types.String `tfsdk:"region"`
}
