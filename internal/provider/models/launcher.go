// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package models

import (
	"context"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/client"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// LauncherReference represents a reference to a workflow or other resource.
type LauncherReference struct {
	Type types.String `tfsdk:"type"`
	ID   types.String `tfsdk:"id"`
}

// Launcher represents the Terraform model for a SailPoint Launcher.
type Launcher struct {
	ID          types.String         `tfsdk:"id"`
	Name        types.String         `tfsdk:"name"`
	Description types.String         `tfsdk:"description"`
	Type        types.String         `tfsdk:"type"`
	Disabled    types.Bool           `tfsdk:"disabled"`
	Reference   *LauncherReference   `tfsdk:"reference"`
	Config      jsontypes.Normalized `tfsdk:"config"`
	Owner       types.Object         `tfsdk:"owner"`
	Created     types.String         `tfsdk:"created"`
	Modified    types.String         `tfsdk:"modified"`
}

// ConvertToSailPoint converts the Terraform LauncherReference to a SailPoint API LauncherReference.
func (r *LauncherReference) ConvertToSailPoint(ctx context.Context) *client.LauncherReference {
	if r == nil {
		return nil
	}

	return &client.LauncherReference{
		Type: r.Type.ValueString(),
		ID:   r.ID.ValueString(),
	}
}

// ConvertFromSailPointForResource converts a SailPoint API LauncherReference to the Terraform model for resources.
func (r *LauncherReference) ConvertFromSailPointForResource(ctx context.Context, ref *client.LauncherReference) {
	if r == nil || ref == nil {
		return
	}

	r.Type = types.StringValue(ref.Type)
	r.ID = types.StringValue(ref.ID)
}

// ConvertFromSailPointForDataSource converts a SailPoint API LauncherReference to the Terraform model for data sources.
func (r *LauncherReference) ConvertFromSailPointForDataSource(ctx context.Context, ref *client.LauncherReference) {
	if r == nil || ref == nil {
		return
	}

	if ref.Type != "" {
		r.Type = types.StringValue(ref.Type)
	}
	if ref.ID != "" {
		r.ID = types.StringValue(ref.ID)
	}
}

// ConvertToSailPoint converts the Terraform model to a SailPoint API Launcher.
func (l *Launcher) ConvertToSailPoint(ctx context.Context) *client.Launcher {
	if l == nil {
		return nil
	}

	launcher := &client.Launcher{
		Name:        l.Name.ValueString(),
		Description: l.Description.ValueString(),
		Type:        l.Type.ValueString(),
		Disabled:    l.Disabled.ValueBool(),
		Config:      l.Config.ValueString(),
	}

	// Convert reference
	if l.Reference != nil {
		launcher.Reference = l.Reference.ConvertToSailPoint(ctx)
	}

	return launcher
}

// ConvertFromSailPoint converts a SailPoint API Launcher to the Terraform model.
// For resources, set includeNull to true. For data sources, set to false.
func (l *Launcher) ConvertFromSailPoint(ctx context.Context, launcher *client.Launcher, includeNull bool) {
	if l == nil || launcher == nil {
		return
	}

	l.ID = types.StringValue(launcher.ID)
	l.Name = types.StringValue(launcher.Name)
	l.Description = types.StringValue(launcher.Description)
	l.Type = types.StringValue(launcher.Type)
	l.Disabled = types.BoolValue(launcher.Disabled)
	l.Config = jsontypes.NewNormalizedValue(launcher.Config)

	// Convert reference
	if launcher.Reference != nil {
		l.Reference = &LauncherReference{}
		if includeNull {
			l.Reference.ConvertFromSailPointForResource(ctx, launcher.Reference)
		} else {
			l.Reference.ConvertFromSailPointForDataSource(ctx, launcher.Reference)
		}
	} else if includeNull {
		l.Reference = nil
	}

	// Convert owner ObjectRef to types.Object
	if launcher.Owner != nil {
		ownerAttrs := map[string]attr.Value{
			"type": types.StringValue(launcher.Owner.Type),
			"id":   types.StringValue(launcher.Owner.ID),
		}
		if launcher.Owner.Name != "" {
			ownerAttrs["name"] = types.StringValue(launcher.Owner.Name)
		} else if includeNull {
			ownerAttrs["name"] = types.StringNull()
		} else {
			ownerAttrs["name"] = types.StringValue("")
		}

		ownerObj, diag := types.ObjectValue(
			map[string]attr.Type{
				"type": types.StringType,
				"id":   types.StringType,
				"name": types.StringType,
			},
			ownerAttrs,
		)
		if !diag.HasError() {
			l.Owner = ownerObj
		}
	} else if includeNull {
		l.Owner = types.ObjectNull(map[string]attr.Type{
			"type": types.StringType,
			"id":   types.StringType,
			"name": types.StringType,
		})
	}

	// Handle computed fields
	if launcher.Created != nil {
		l.Created = types.StringValue(*launcher.Created)
	} else if includeNull {
		l.Created = types.StringNull()
	}

	if launcher.Modified != nil {
		l.Modified = types.StringValue(*launcher.Modified)
	} else if includeNull {
		l.Modified = types.StringNull()
	}
}

// ConvertFromSailPointForResource converts for resource operations (includes all fields).
func (l *Launcher) ConvertFromSailPointForResource(ctx context.Context, launcher *client.Launcher) {
	l.ConvertFromSailPoint(ctx, launcher, true)
}

// ConvertFromSailPointForDataSource converts for data source operations (preserves state).
func (l *Launcher) ConvertFromSailPointForDataSource(ctx context.Context, launcher *client.Launcher) {
	l.ConvertFromSailPoint(ctx, launcher, false)
}
