// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managed_cluster

import "github.com/hashicorp/terraform-plugin-framework/types"

type ManagedClusterDataSourceModel struct {
	ID                   types.String             `tfsdk:"id"`
	Name                 types.String             `tfsdk:"name"`
	Pod                  types.String             `tfsdk:"pod"`
	Org                  types.String             `tfsdk:"org"`
	Type                 types.String             `tfsdk:"type"`
	Configuration        types.Map                `tfsdk:"configuration"`
	KeyPair              ManagedClusterKeyPair    `tfsdk:"key_pair"`
	Attributes           ManagedClusterAttributes `tfsdk:"attributes"`
	Description          types.String             `tfsdk:"description"`
	Redis                ManagedClusterRedis      `tfsdk:"redis"`
	ClientType           types.String             `tfsdk:"client_type"`
	CcgVersion           types.String             `tfsdk:"ccg_version"`
	PinnedConfig         types.Bool               `tfsdk:"pinned_config"`
	Operational          types.Bool               `tfsdk:"operational"`
	Status               types.String             `tfsdk:"status"`
	PublicKeyCertificate types.String             `tfsdk:"public_key_certificate"`
	PublicKeyThumbprint  types.String             `tfsdk:"public_key_thumbprint"`
	PublicKey            types.String             `tfsdk:"public_key_type"`
	AlertKey             types.String             `tfsdk:"alert_key"`
	ClientIds            types.List               `tfsdk:"client_ids"`
	ServiceCount         types.Int32              `tfsdk:"service_count"`
	CcId                 types.String             `tfsdk:"cc_id"`
	CreatedAt            types.String             `tfsdk:"created_at"`
	UpdatedAt            types.String             `tfsdk:"updated_at"`
}
