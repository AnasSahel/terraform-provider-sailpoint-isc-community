// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package connector

import (
	"context"
	"testing"

	api_v2025 "github.com/sailpoint-oss/golang-sdk/v2/api_v2025"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnectorSummaryFromAPI(t *testing.T) {
	ctx := context.Background()

	t.Run("complete connector", func(t *testing.T) {
		input := api_v2025.V3ConnectorDto{
			Name:          stringPtr("Active Directory"),
			Type:          stringPtr("active-directory"),
			ScriptName:    stringPtr("active-directory"),
			DirectConnect: boolPtr(true),
			Status:        stringPtr("RELEASED"),
			Features:      []string{"PROVISIONING", "SYNC_PROVISIONING"},
		}

		var model ConnectorSummaryModel
		err := model.FromSailPointV3ConnectorDto(ctx, &input)
		require.NoError(t, err)

		assert.Equal(t, "active-directory", model.ID.ValueString())
		assert.Equal(t, "Active Directory", model.Name.ValueString())
		assert.Equal(t, "active-directory", model.Type.ValueString())
		assert.Equal(t, "active-directory", model.ScriptName.ValueString())
		assert.Equal(t, true, model.DirectConnect.ValueBool())
		assert.Equal(t, "RELEASED", model.Status.ValueString())
		assert.True(t, model.Category.IsNull())
		assert.True(t, model.Labels.IsNull())
		assert.False(t, model.Features.IsNull())
	})

	t.Run("minimal connector", func(t *testing.T) {
		input := api_v2025.V3ConnectorDto{
			Name:       stringPtr("Test Connector"),
			ScriptName: stringPtr("test-connector"),
		}

		var model ConnectorSummaryModel
		err := model.FromSailPointV3ConnectorDto(ctx, &input)
		require.NoError(t, err)

		assert.Equal(t, "test-connector", model.ID.ValueString())
		assert.Equal(t, "Test Connector", model.Name.ValueString())
		assert.True(t, model.Type.IsNull())
		assert.Equal(t, "test-connector", model.ScriptName.ValueString())
		assert.True(t, model.DirectConnect.IsNull())
		assert.True(t, model.Status.IsNull())
	})

	t.Run("connector without script name", func(t *testing.T) {
		input := api_v2025.V3ConnectorDto{
			Name: stringPtr("Test Connector"),
			Type: stringPtr("test"),
		}

		var model ConnectorSummaryModel
		err := model.FromSailPointV3ConnectorDto(ctx, &input)
		require.NoError(t, err)

		// Should have generated ID from name hash
		assert.False(t, model.ID.IsNull())
		assert.Equal(t, "Test Connector", model.Name.ValueString())
		assert.True(t, model.ScriptName.IsNull())
	})
}

// Helper functions to create pointers
func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}
