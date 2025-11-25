// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package models

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/client"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AccessProfile represents the Terraform model for a SailPoint Access Profile.
type AccessProfile struct {
	ID                   types.String         `tfsdk:"id"`
	Name                 types.String         `tfsdk:"name"`
	Description          types.String         `tfsdk:"description"`
	Created              types.String         `tfsdk:"created"`
	Modified             types.String         `tfsdk:"modified"`
	Enabled              types.Bool           `tfsdk:"enabled"`
	Requestable          types.Bool           `tfsdk:"requestable"`
	Owner                *ObjectRef           `tfsdk:"owner"`
	Source               *ObjectRef           `tfsdk:"source"`
	Entitlements         types.List           `tfsdk:"entitlements"` // List of ObjectRef
	Segments             types.List           `tfsdk:"segments"`     // List of String (UUIDs)
	AccessRequestConfig  jsontypes.Normalized `tfsdk:"access_request_config"`
	RevokeRequestConfig  jsontypes.Normalized `tfsdk:"revoke_request_config"`
	ProvisioningCriteria jsontypes.Normalized `tfsdk:"provisioning_criteria"`
}

// ConvertToSailPoint converts the Terraform model to a SailPoint API AccessProfile.
func (a *AccessProfile) ConvertToSailPoint(ctx context.Context) (*client.AccessProfile, error) {
	if a == nil {
		return nil, nil
	}

	accessProfile := &client.AccessProfile{
		Name: a.Name.ValueString(),
	}

	// Convert optional scalar fields
	if !a.Description.IsNull() && !a.Description.IsUnknown() {
		desc := a.Description.ValueString()
		accessProfile.Description = &desc
	}

	if !a.Enabled.IsNull() && !a.Enabled.IsUnknown() {
		enabled := a.Enabled.ValueBool()
		accessProfile.Enabled = &enabled
	}

	if !a.Requestable.IsNull() && !a.Requestable.IsUnknown() {
		requestable := a.Requestable.ValueBool()
		accessProfile.Requestable = &requestable
	}

	// Convert owner ObjectRef
	if a.Owner != nil {
		ownerRef := a.Owner.ConvertToSailPoint(ctx)
		accessProfile.Owner = &ownerRef
	}

	// Convert source ObjectRef
	if a.Source != nil {
		sourceRef := a.Source.ConvertToSailPoint(ctx)
		accessProfile.Source = &sourceRef
	}

	// Convert entitlements list
	if !a.Entitlements.IsNull() && !a.Entitlements.IsUnknown() {
		var entitlementRefs []ObjectRef
		diags := a.Entitlements.ElementsAs(ctx, &entitlementRefs, false)
		if diags.HasError() {
			return nil, fmt.Errorf("error converting entitlements list: %v", diags)
		}

		accessProfile.Entitlements = make([]client.ObjectRef, 0, len(entitlementRefs))
		for _, ref := range entitlementRefs {
			accessProfile.Entitlements = append(accessProfile.Entitlements, ref.ConvertToSailPoint(ctx))
		}
	}

	// Convert segments list
	if !a.Segments.IsNull() && !a.Segments.IsUnknown() {
		var segments []string
		diags := a.Segments.ElementsAs(ctx, &segments, false)
		if diags.HasError() {
			return nil, fmt.Errorf("error converting segments list: %v", diags)
		}
		accessProfile.Segments = segments
	}

	// Convert JSON config objects
	if !a.AccessRequestConfig.IsNull() && !a.AccessRequestConfig.IsUnknown() {
		var config map[string]interface{}
		if err := json.Unmarshal([]byte(a.AccessRequestConfig.ValueString()), &config); err != nil {
			return nil, err
		}
		accessProfile.AccessRequestConfig = config
	}

	if !a.RevokeRequestConfig.IsNull() && !a.RevokeRequestConfig.IsUnknown() {
		var config map[string]interface{}
		if err := json.Unmarshal([]byte(a.RevokeRequestConfig.ValueString()), &config); err != nil {
			return nil, err
		}
		accessProfile.RevokeRequestConfig = config
	}

	if !a.ProvisioningCriteria.IsNull() && !a.ProvisioningCriteria.IsUnknown() {
		var criteria map[string]interface{}
		if err := json.Unmarshal([]byte(a.ProvisioningCriteria.ValueString()), &criteria); err != nil {
			return nil, err
		}
		accessProfile.ProvisioningCriteria = criteria
	}

	return accessProfile, nil
}

// ConvertFromSailPoint converts a SailPoint API AccessProfile to the Terraform model.
// For resources, set includeNull to true. For data sources, set to false.
func (a *AccessProfile) ConvertFromSailPoint(ctx context.Context, accessProfile *client.AccessProfile, includeNull bool) error {
	if a == nil || accessProfile == nil {
		return nil
	}

	a.ID = types.StringValue(accessProfile.ID)
	a.Name = types.StringValue(accessProfile.Name)

	// Convert optional scalar fields
	if accessProfile.Description != nil {
		a.Description = types.StringValue(*accessProfile.Description)
	} else if includeNull {
		a.Description = types.StringNull()
	}

	if accessProfile.Created != nil {
		a.Created = types.StringValue(*accessProfile.Created)
	} else if includeNull {
		a.Created = types.StringNull()
	}

	if accessProfile.Modified != nil {
		a.Modified = types.StringValue(*accessProfile.Modified)
	} else if includeNull {
		a.Modified = types.StringNull()
	}

	if accessProfile.Enabled != nil {
		a.Enabled = types.BoolValue(*accessProfile.Enabled)
	} else if includeNull {
		a.Enabled = types.BoolNull()
	}

	if accessProfile.Requestable != nil {
		a.Requestable = types.BoolValue(*accessProfile.Requestable)
	} else if includeNull {
		a.Requestable = types.BoolNull()
	}

	// Convert owner ObjectRef
	if accessProfile.Owner != nil {
		a.Owner = &ObjectRef{}
		if includeNull {
			a.Owner.ConvertFromSailPointForResource(ctx, accessProfile.Owner)
		} else {
			a.Owner.ConvertFromSailPointForDataSource(ctx, accessProfile.Owner)
		}
	} else if includeNull {
		a.Owner = nil
	}

	// Convert source ObjectRef
	if accessProfile.Source != nil {
		a.Source = &ObjectRef{}
		if includeNull {
			a.Source.ConvertFromSailPointForResource(ctx, accessProfile.Source)
		} else {
			a.Source.ConvertFromSailPointForDataSource(ctx, accessProfile.Source)
		}
	} else if includeNull {
		a.Source = nil
	}

	// Convert entitlements list
	if len(accessProfile.Entitlements) > 0 {
		entitlementElements := make([]attr.Value, 0, len(accessProfile.Entitlements))
		for _, entRef := range accessProfile.Entitlements {
			modelRef := &ObjectRef{}
			if includeNull {
				modelRef.ConvertFromSailPointForResource(ctx, &entRef)
			} else {
				modelRef.ConvertFromSailPointForDataSource(ctx, &entRef)
			}

			objVal, diags := types.ObjectValueFrom(ctx, map[string]attr.Type{
				"type": types.StringType,
				"id":   types.StringType,
				"name": types.StringType,
			}, modelRef)
			if diags.HasError() {
				return fmt.Errorf("error creating entitlement object: %v", diags)
			}
			entitlementElements = append(entitlementElements, objVal)
		}

		listVal, diags := types.ListValue(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"type": types.StringType,
				"id":   types.StringType,
				"name": types.StringType,
			},
		}, entitlementElements)
		if diags.HasError() {
			return fmt.Errorf("error creating entitlements list: %v", diags)
		}
		a.Entitlements = listVal
	} else if includeNull {
		a.Entitlements = types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"type": types.StringType,
				"id":   types.StringType,
				"name": types.StringType,
			},
		})
	}

	// Convert segments list
	if len(accessProfile.Segments) > 0 {
		segmentElements := make([]attr.Value, 0, len(accessProfile.Segments))
		for _, seg := range accessProfile.Segments {
			segmentElements = append(segmentElements, types.StringValue(seg))
		}

		listVal, diags := types.ListValue(types.StringType, segmentElements)
		if diags.HasError() {
			return fmt.Errorf("error creating segments list: %v", diags)
		}
		a.Segments = listVal
	} else if includeNull {
		a.Segments = types.ListNull(types.StringType)
	}

	// Convert JSON config objects
	if accessProfile.AccessRequestConfig != nil {
		configJSON, err := json.Marshal(accessProfile.AccessRequestConfig)
		if err != nil {
			return err
		}
		a.AccessRequestConfig = jsontypes.NewNormalizedValue(string(configJSON))
	} else if includeNull {
		a.AccessRequestConfig = jsontypes.NewNormalizedNull()
	}

	if accessProfile.RevokeRequestConfig != nil {
		configJSON, err := json.Marshal(accessProfile.RevokeRequestConfig)
		if err != nil {
			return err
		}
		a.RevokeRequestConfig = jsontypes.NewNormalizedValue(string(configJSON))
	} else if includeNull {
		a.RevokeRequestConfig = jsontypes.NewNormalizedNull()
	}

	if accessProfile.ProvisioningCriteria != nil {
		criteriaJSON, err := json.Marshal(accessProfile.ProvisioningCriteria)
		if err != nil {
			return err
		}
		a.ProvisioningCriteria = jsontypes.NewNormalizedValue(string(criteriaJSON))
	} else if includeNull {
		a.ProvisioningCriteria = jsontypes.NewNormalizedNull()
	}

	return nil
}

// ConvertFromSailPointForResource converts for resource operations (includes all fields).
func (a *AccessProfile) ConvertFromSailPointForResource(ctx context.Context, accessProfile *client.AccessProfile) error {
	return a.ConvertFromSailPoint(ctx, accessProfile, true)
}

// ConvertFromSailPointForDataSource converts for data source operations (preserves state).
func (a *AccessProfile) ConvertFromSailPointForDataSource(ctx context.Context, accessProfile *client.AccessProfile) error {
	return a.ConvertFromSailPoint(ctx, accessProfile, false)
}

// GeneratePatchOperations generates JSON Patch operations for updating an access profile.
func (a *AccessProfile) GeneratePatchOperations(ctx context.Context, newAccessProfile AccessProfile) ([]map[string]interface{}, error) {
	operations := []map[string]interface{}{}

	// Compare and generate patch for name
	if !a.Name.Equal(newAccessProfile.Name) {
		operations = append(operations, map[string]interface{}{
			"op":    "replace",
			"path":  "/name",
			"value": newAccessProfile.Name.ValueString(),
		})
	}

	// Compare and generate patch for description
	if !a.Description.Equal(newAccessProfile.Description) {
		if newAccessProfile.Description.IsNull() {
			operations = append(operations, map[string]interface{}{
				"op":   "remove",
				"path": "/description",
			})
		} else {
			operations = append(operations, map[string]interface{}{
				"op":    "replace",
				"path":  "/description",
				"value": newAccessProfile.Description.ValueString(),
			})
		}
	}

	// Compare and generate patch for enabled
	if !a.Enabled.Equal(newAccessProfile.Enabled) {
		if newAccessProfile.Enabled.IsNull() {
			operations = append(operations, map[string]interface{}{
				"op":   "remove",
				"path": "/enabled",
			})
		} else {
			operations = append(operations, map[string]interface{}{
				"op":    "replace",
				"path":  "/enabled",
				"value": newAccessProfile.Enabled.ValueBool(),
			})
		}
	}

	// Compare and generate patch for requestable
	if !a.Requestable.Equal(newAccessProfile.Requestable) {
		if newAccessProfile.Requestable.IsNull() {
			operations = append(operations, map[string]interface{}{
				"op":   "remove",
				"path": "/requestable",
			})
		} else {
			operations = append(operations, map[string]interface{}{
				"op":    "replace",
				"path":  "/requestable",
				"value": newAccessProfile.Requestable.ValueBool(),
			})
		}
	}

	// Compare and generate patch for owner
	if (a.Owner == nil && newAccessProfile.Owner != nil) ||
		(a.Owner != nil && newAccessProfile.Owner == nil) ||
		(a.Owner != nil && newAccessProfile.Owner != nil && !a.Owner.Equals(newAccessProfile.Owner)) {
		if newAccessProfile.Owner == nil {
			operations = append(operations, map[string]interface{}{
				"op":   "remove",
				"path": "/owner",
			})
		} else {
			operations = append(operations, map[string]interface{}{
				"op":    "replace",
				"path":  "/owner",
				"value": newAccessProfile.Owner.ConvertToSailPoint(ctx),
			})
		}
	}

	// Compare and generate patch for source
	if (a.Source == nil && newAccessProfile.Source != nil) ||
		(a.Source != nil && newAccessProfile.Source == nil) ||
		(a.Source != nil && newAccessProfile.Source != nil && !a.Source.Equals(newAccessProfile.Source)) {
		if newAccessProfile.Source == nil {
			operations = append(operations, map[string]interface{}{
				"op":   "remove",
				"path": "/source",
			})
		} else {
			operations = append(operations, map[string]interface{}{
				"op":    "replace",
				"path":  "/source",
				"value": newAccessProfile.Source.ConvertToSailPoint(ctx),
			})
		}
	}

	// Compare and generate patch for entitlements
	if !a.Entitlements.Equal(newAccessProfile.Entitlements) {
		if newAccessProfile.Entitlements.IsNull() {
			operations = append(operations, map[string]interface{}{
				"op":   "remove",
				"path": "/entitlements",
			})
		} else {
			var entitlementRefs []ObjectRef
			diags := newAccessProfile.Entitlements.ElementsAs(ctx, &entitlementRefs, false)
			if diags.HasError() {
				return nil, fmt.Errorf("error converting entitlements for patch: %v", diags)
			}

			entitlements := make([]client.ObjectRef, 0, len(entitlementRefs))
			for _, ref := range entitlementRefs {
				entitlements = append(entitlements, ref.ConvertToSailPoint(ctx))
			}

			operations = append(operations, map[string]interface{}{
				"op":    "replace",
				"path":  "/entitlements",
				"value": entitlements,
			})
		}
	}

	// Compare and generate patch for segments
	if !a.Segments.Equal(newAccessProfile.Segments) {
		if newAccessProfile.Segments.IsNull() {
			operations = append(operations, map[string]interface{}{
				"op":   "remove",
				"path": "/segments",
			})
		} else {
			var segments []string
			diags := newAccessProfile.Segments.ElementsAs(ctx, &segments, false)
			if diags.HasError() {
				return nil, fmt.Errorf("error converting segments for patch: %v", diags)
			}

			operations = append(operations, map[string]interface{}{
				"op":    "replace",
				"path":  "/segments",
				"value": segments,
			})
		}
	}

	// Compare and generate patch for accessRequestConfig
	if !a.AccessRequestConfig.Equal(newAccessProfile.AccessRequestConfig) {
		if newAccessProfile.AccessRequestConfig.IsNull() {
			operations = append(operations, map[string]interface{}{
				"op":   "remove",
				"path": "/accessRequestConfig",
			})
		} else {
			var config map[string]interface{}
			if err := json.Unmarshal([]byte(newAccessProfile.AccessRequestConfig.ValueString()), &config); err != nil {
				return nil, err
			}
			operations = append(operations, map[string]interface{}{
				"op":    "replace",
				"path":  "/accessRequestConfig",
				"value": config,
			})
		}
	}

	// Compare and generate patch for revokeRequestConfig
	if !a.RevokeRequestConfig.Equal(newAccessProfile.RevokeRequestConfig) {
		if newAccessProfile.RevokeRequestConfig.IsNull() {
			operations = append(operations, map[string]interface{}{
				"op":   "remove",
				"path": "/revokeRequestConfig",
			})
		} else {
			var config map[string]interface{}
			if err := json.Unmarshal([]byte(newAccessProfile.RevokeRequestConfig.ValueString()), &config); err != nil {
				return nil, err
			}
			operations = append(operations, map[string]interface{}{
				"op":    "replace",
				"path":  "/revokeRequestConfig",
				"value": config,
			})
		}
	}

	// Compare and generate patch for provisioningCriteria
	if !a.ProvisioningCriteria.Equal(newAccessProfile.ProvisioningCriteria) {
		if newAccessProfile.ProvisioningCriteria.IsNull() {
			operations = append(operations, map[string]interface{}{
				"op":   "remove",
				"path": "/provisioningCriteria",
			})
		} else {
			var criteria map[string]interface{}
			if err := json.Unmarshal([]byte(newAccessProfile.ProvisioningCriteria.ValueString()), &criteria); err != nil {
				return nil, err
			}
			operations = append(operations, map[string]interface{}{
				"op":    "replace",
				"path":  "/provisioningCriteria",
				"value": criteria,
			})
		}
	}

	return operations, nil
}
