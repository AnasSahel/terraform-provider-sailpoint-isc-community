// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package models

import (
	"context"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/client"
)

type ModelWithSourceTerraformConversionMethods[T any] interface {
	ConvertFromSailPointForResource(ctx context.Context, source *T)
	ConvertFromSailPointForDataSource(ctx context.Context, source *T)
}

type ModelWithSailPointConversionMethods[T any] interface {
	ConvertToSailPoint(ctx context.Context) T
}

type ModelWithSailPointPatchMethods[T any] interface {
	BuildPatchOptions(ctx context.Context, desired *T) []client.JSONPatchOperation
}
