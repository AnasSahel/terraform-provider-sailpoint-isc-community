// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managed_cluster

import "github.com/hashicorp/terraform-plugin-framework/types"

type ManagedClusterKeyPair struct {
	PublicKey            types.String `tfsdk:"public_key"`
	PublicKeyThumbprint  types.String `tfsdk:"public_key_thumbprint"`
	PublicKeyCertificate types.String `tfsdk:"public_key_certificate"`
}
