// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managed_cluster

import "github.com/hashicorp/terraform-plugin-framework/types"

type ManagedClusterRedis struct {
	RedisHost types.String `tfsdk:"redis_host"`
	RedisPort types.Int32  `tfsdk:"redis_port"`
}
