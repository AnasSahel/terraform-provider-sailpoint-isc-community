// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package models

import (
	"context"
	"encoding/json"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// EntitlementDirectPermission represents a permission with rights and target.
type EntitlementDirectPermission struct {
	Rights types.List   `tfsdk:"rights"`
	Target types.String `tfsdk:"target"`
}

// EntitlementManuallyUpdatedFields tracks which fields have been manually updated.
type EntitlementManuallyUpdatedFields struct {
	DisplayName types.Bool `tfsdk:"display_name"`
	Description types.Bool `tfsdk:"description"`
}

// EntitlementAccessModelMetadataValue represents a value in access model metadata.
type EntitlementAccessModelMetadataValue struct {
	Value  types.String `tfsdk:"value"`
	Name   types.String `tfsdk:"name"`
	Status types.String `tfsdk:"status"`
}

// EntitlementAccessModelMetadataAttribute represents an attribute in access model metadata.
type EntitlementAccessModelMetadataAttribute struct {
	Key         types.String `tfsdk:"key"`
	Name        types.String `tfsdk:"name"`
	Multiselect types.Bool   `tfsdk:"multiselect"`
	Status      types.String `tfsdk:"status"`
	Type        types.String `tfsdk:"type"`
	ObjectTypes types.List   `tfsdk:"object_types"`
	Description types.String `tfsdk:"description"`
	Values      types.List   `tfsdk:"values"`
}

// EntitlementAccessModelMetadata represents access model metadata for an entitlement.
type EntitlementAccessModelMetadata struct {
	Attributes types.List `tfsdk:"attributes"`
}

// Entitlement represents the Terraform model for a SailPoint Entitlement.
type Entitlement struct {
	ID                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	Created                types.String `tfsdk:"created"`
	Modified               types.String `tfsdk:"modified"`
	Description            types.String `tfsdk:"description"`
	Attribute              types.String `tfsdk:"attribute"`
	Value                  types.String `tfsdk:"value"`
	SourceSchemaObjectType types.String `tfsdk:"source_schema_object_type"`
	Privileged             types.Bool   `tfsdk:"privileged"`
	Requestable            types.Bool   `tfsdk:"requestable"`
	CloudGoverned          types.Bool   `tfsdk:"cloud_governed"`
	Source                 types.Object `tfsdk:"source"`
	Owner                  types.Object `tfsdk:"owner"`
	Attributes             types.String `tfsdk:"attributes"`
	DirectPermissions      types.List   `tfsdk:"direct_permissions"`
	Segments               types.List   `tfsdk:"segments"`
	ManuallyUpdatedFields  types.Object `tfsdk:"manually_updated_fields"`
	AccessModelMetadata    types.Object `tfsdk:"access_model_metadata"`
}

// ConvertFromSailPointForDataSource converts a SailPoint API Entitlement to the Terraform model for data sources.
func (e *Entitlement) ConvertFromSailPointForDataSource(ctx context.Context, entitlement *client.Entitlement) {
	if e == nil || entitlement == nil {
		return
	}

	e.ID = types.StringValue(entitlement.ID)
	e.Name = types.StringValue(entitlement.Name)

	// Handle optional string fields
	if entitlement.Created != nil {
		e.Created = types.StringValue(*entitlement.Created)
	}

	if entitlement.Modified != nil {
		e.Modified = types.StringValue(*entitlement.Modified)
	}

	if entitlement.Description != nil {
		e.Description = types.StringValue(*entitlement.Description)
	}

	if entitlement.Attribute != nil {
		e.Attribute = types.StringValue(*entitlement.Attribute)
	}

	if entitlement.Value != nil {
		e.Value = types.StringValue(*entitlement.Value)
	}

	if entitlement.SourceSchemaObjectType != nil {
		e.SourceSchemaObjectType = types.StringValue(*entitlement.SourceSchemaObjectType)
	}

	// Handle boolean fields
	if entitlement.Privileged != nil {
		e.Privileged = types.BoolValue(*entitlement.Privileged)
	}

	if entitlement.Requestable != nil {
		e.Requestable = types.BoolValue(*entitlement.Requestable)
	}

	if entitlement.CloudGoverned != nil {
		e.CloudGoverned = types.BoolValue(*entitlement.CloudGoverned)
	}

	// Convert source ObjectRef to types.Object
	if entitlement.Source != nil {
		sourceAttrs := map[string]attr.Value{
			"type": types.StringValue(entitlement.Source.Type),
			"id":   types.StringValue(entitlement.Source.ID),
		}
		if entitlement.Source.Name != "" {
			sourceAttrs["name"] = types.StringValue(entitlement.Source.Name)
		} else {
			sourceAttrs["name"] = types.StringValue("")
		}

		sourceObj, diag := types.ObjectValue(
			map[string]attr.Type{
				"type": types.StringType,
				"id":   types.StringType,
				"name": types.StringType,
			},
			sourceAttrs,
		)
		if !diag.HasError() {
			e.Source = sourceObj
		}
	} else {
		e.Source = types.ObjectNull(map[string]attr.Type{
			"type": types.StringType,
			"id":   types.StringType,
			"name": types.StringType,
		})
	}

	// Convert owner ObjectRef to types.Object
	if entitlement.Owner != nil {
		ownerAttrs := map[string]attr.Value{
			"type": types.StringValue(entitlement.Owner.Type),
			"id":   types.StringValue(entitlement.Owner.ID),
		}
		if entitlement.Owner.Name != "" {
			ownerAttrs["name"] = types.StringValue(entitlement.Owner.Name)
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
			e.Owner = ownerObj
		}
	} else {
		e.Owner = types.ObjectNull(map[string]attr.Type{
			"type": types.StringType,
			"id":   types.StringType,
			"name": types.StringType,
		})
	}

	// Convert attributes map to JSON string
	if len(entitlement.Attributes) > 0 {
		attrsJSON, err := json.Marshal(entitlement.Attributes)
		if err == nil {
			e.Attributes = types.StringValue(string(attrsJSON))
		}
	}

	// Convert direct permissions to list
	if len(entitlement.DirectPermissions) > 0 {
		var permElements []attr.Value
		for _, perm := range entitlement.DirectPermissions {
			// Convert rights to list
			var rightsElements []attr.Value
			for _, right := range perm.Rights {
				rightsElements = append(rightsElements, types.StringValue(right))
			}
			rightsList, _ := types.ListValue(types.StringType, rightsElements)

			// Create target string
			target := types.StringNull()
			if perm.Target != nil {
				target = types.StringValue(*perm.Target)
			}

			// Create permission object
			permObj, _ := types.ObjectValue(
				map[string]attr.Type{
					"rights": types.ListType{ElemType: types.StringType},
					"target": types.StringType,
				},
				map[string]attr.Value{
					"rights": rightsList,
					"target": target,
				},
			)
			permElements = append(permElements, permObj)
		}

		permsList, _ := types.ListValue(
			types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"rights": types.ListType{ElemType: types.StringType},
					"target": types.StringType,
				},
			},
			permElements,
		)
		e.DirectPermissions = permsList
	} else {
		e.DirectPermissions = types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"rights": types.ListType{ElemType: types.StringType},
				"target": types.StringType,
			},
		})
	}

	// Convert segments to list
	if len(entitlement.Segments) > 0 {
		var segmentElements []attr.Value
		for _, seg := range entitlement.Segments {
			segmentElements = append(segmentElements, types.StringValue(seg))
		}
		segmentsList, _ := types.ListValue(types.StringType, segmentElements)
		e.Segments = segmentsList
	} else {
		e.Segments = types.ListNull(types.StringType)
	}

	// Convert manually updated fields
	if entitlement.ManuallyUpdatedFields != nil {
		fieldsAttrs := map[string]attr.Value{}

		if entitlement.ManuallyUpdatedFields.DisplayName != nil {
			fieldsAttrs["display_name"] = types.BoolValue(*entitlement.ManuallyUpdatedFields.DisplayName)
		} else {
			fieldsAttrs["display_name"] = types.BoolNull()
		}

		if entitlement.ManuallyUpdatedFields.Description != nil {
			fieldsAttrs["description"] = types.BoolValue(*entitlement.ManuallyUpdatedFields.Description)
		} else {
			fieldsAttrs["description"] = types.BoolNull()
		}

		fieldsObj, diag := types.ObjectValue(
			map[string]attr.Type{
				"display_name": types.BoolType,
				"description":  types.BoolType,
			},
			fieldsAttrs,
		)
		if !diag.HasError() {
			e.ManuallyUpdatedFields = fieldsObj
		}
	}

	// Convert access model metadata to object with attributes list
	if entitlement.AccessModelMetadata != nil && len(entitlement.AccessModelMetadata.Attributes) > 0 {
		var attributeElements []attr.Value

		// Define the attribute object type structure
		attributeObjType := types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"key":          types.StringType,
				"name":         types.StringType,
				"multiselect":  types.BoolType,
				"status":       types.StringType,
				"type":         types.StringType,
				"object_types": types.ListType{ElemType: types.StringType},
				"description":  types.StringType,
				"values": types.ListType{
					ElemType: types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"value":  types.StringType,
							"name":   types.StringType,
							"status": types.StringType,
						},
					},
				},
			},
		}

		for _, attribute := range entitlement.AccessModelMetadata.Attributes {
			attrMap := map[string]attr.Value{}

			// Convert simple string fields
			if attribute.Key != nil {
				attrMap["key"] = types.StringValue(*attribute.Key)
			} else {
				attrMap["key"] = types.StringNull()
			}

			if attribute.Name != nil {
				attrMap["name"] = types.StringValue(*attribute.Name)
			} else {
				attrMap["name"] = types.StringNull()
			}

			if attribute.Multiselect != nil {
				attrMap["multiselect"] = types.BoolValue(*attribute.Multiselect)
			} else {
				attrMap["multiselect"] = types.BoolNull()
			}

			if attribute.Status != nil {
				attrMap["status"] = types.StringValue(*attribute.Status)
			} else {
				attrMap["status"] = types.StringNull()
			}

			if attribute.Type != nil {
				attrMap["type"] = types.StringValue(*attribute.Type)
			} else {
				attrMap["type"] = types.StringNull()
			}

			if attribute.Description != nil {
				attrMap["description"] = types.StringValue(*attribute.Description)
			} else {
				attrMap["description"] = types.StringNull()
			}

			// Convert objectTypes array
			if len(attribute.ObjectTypes) > 0 {
				var objectTypeElements []attr.Value
				for _, objType := range attribute.ObjectTypes {
					objectTypeElements = append(objectTypeElements, types.StringValue(objType))
				}
				objectTypesList, _ := types.ListValue(types.StringType, objectTypeElements)
				attrMap["object_types"] = objectTypesList
			} else {
				attrMap["object_types"] = types.ListNull(types.StringType)
			}

			// Convert values array
			if len(attribute.Values) > 0 {
				var valueElements []attr.Value
				for _, val := range attribute.Values {
					valueMap := map[string]attr.Value{}

					if val.Value != nil {
						valueMap["value"] = types.StringValue(*val.Value)
					} else {
						valueMap["value"] = types.StringNull()
					}

					if val.Name != nil {
						valueMap["name"] = types.StringValue(*val.Name)
					} else {
						valueMap["name"] = types.StringNull()
					}

					if val.Status != nil {
						valueMap["status"] = types.StringValue(*val.Status)
					} else {
						valueMap["status"] = types.StringNull()
					}

					valueObj, _ := types.ObjectValue(
						map[string]attr.Type{
							"value":  types.StringType,
							"name":   types.StringType,
							"status": types.StringType,
						},
						valueMap,
					)
					valueElements = append(valueElements, valueObj)
				}

				valuesList, _ := types.ListValue(
					types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"value":  types.StringType,
							"name":   types.StringType,
							"status": types.StringType,
						},
					},
					valueElements,
				)
				attrMap["values"] = valuesList
			} else {
				attrMap["values"] = types.ListNull(types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"value":  types.StringType,
						"name":   types.StringType,
						"status": types.StringType,
					},
				})
			}

			// Create the attribute object
			attrObj, _ := types.ObjectValue(attributeObjType.AttrTypes, attrMap)
			attributeElements = append(attributeElements, attrObj)
		}

		// Create the attributes list
		attributesList, _ := types.ListValue(attributeObjType, attributeElements)

		// Create the outer accessModelMetadata object
		metadataObj, _ := types.ObjectValue(
			map[string]attr.Type{
				"attributes": types.ListType{ElemType: attributeObjType},
			},
			map[string]attr.Value{
				"attributes": attributesList,
			},
		)
		e.AccessModelMetadata = metadataObj
	} else {
		// Set null object when no metadata
		e.AccessModelMetadata = types.ObjectNull(map[string]attr.Type{
			"attributes": types.ListType{
				ElemType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"key":          types.StringType,
						"name":         types.StringType,
						"multiselect":  types.BoolType,
						"status":       types.StringType,
						"type":         types.StringType,
						"object_types": types.ListType{ElemType: types.StringType},
						"description":  types.StringType,
						"values": types.ListType{
							ElemType: types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"value":  types.StringType,
									"name":   types.StringType,
									"status": types.StringType,
								},
							},
						},
					},
				},
			},
		})
	}
}
