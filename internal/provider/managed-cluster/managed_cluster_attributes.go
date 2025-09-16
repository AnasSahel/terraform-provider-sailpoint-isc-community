// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managed_cluster

import "github.com/hashicorp/terraform-plugin-framework/types"

type ManagedClusterAttributes struct {
	Queue    ManagedClusterQueue `tfsdk:"queue"`
	Keystore types.String        `tfsdk:"keystore"`
}
