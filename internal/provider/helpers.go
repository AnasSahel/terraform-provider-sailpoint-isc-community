package provider

// Convert map[string]interface{} → Terraform Dynamic
// Used in: Read operations across all resources
// func convertAnyToDynamic(ctx context.Context, apiMap map[string]interface{}) (types.Dynamic, diag.Diagnostics) {
// 	var diags diag.Diagnostics

// 	if apiMap == nil {
// 		return types.DynamicNull(), diags
// 	}

// 	elements := make(map[string]attr.Value, len(apiMap))
// 	for key, value := range apiMap {
// 		attrValue, err := convertAnyToAttrValue(ctx, value)
// 		if err != nil {
// 			diags.AddError(
// 				"Error Converting API Map to Dynamic",
// 				fmt.Sprintf("Error converting API map value for key %q: %v", key, err),
// 			)
// 			return types.DynamicNull(), diags
// 		}
// 		elements[key] = attrValue
// 	}

// 	if diags.HasError() {
// 		return types.DynamicNull(), diags
// 	}

// 	mapValue, mapDiags := types.MapValue(types.DynamicType, elements)
// 	diags.Append(mapDiags...)
// 	if diags.HasError() {
// 		return types.DynamicNull(), diags
// 	}
// 	return types.DynamicValue(mapValue), diags
// }

// func convertAnyToAttrValue(ctx context.Context, value interface{}) (attr.Value, error) {
// 	if value == nil {
// 		return types.DynamicNull(), nil
// 	}

// 	switch v := value.(type) {
// 	case string:
// 		return types.DynamicValue(types.StringValue(v)), nil
// 	case bool:
// 		return types.DynamicValue(types.BoolValue(v)), nil
// 	case int:
// 		return types.DynamicValue(types.Int64Value(int64(v))), nil
// 	case int32:
// 		return types.DynamicValue(types.Int64Value(int64(v))), nil
// 	case int64:
// 		return types.DynamicValue(types.Int64Value(int64(v))), nil
// 	case float32:
// 		return types.DynamicValue(types.Float64Value(float64(v))), nil
// 	case float64:
// 		return types.DynamicValue(types.Float64Value(float64(v))), nil
// 	case map[string]interface{}:
// 		elements := make(map[string]attr.Value, len(v))
// 		for key, val := range v {
// 			attrValue, err := convertAnyToAttrValue(ctx, val)
// 			if err != nil {
// 				return nil, fmt.Errorf("error in nested map key '%s': %w", key, err)
// 			}
// 			elements[key] = attrValue
// 		}
// 		mapValue, _ := types.MapValue(types.DynamicType, elements)
// 		return types.DynamicValue(mapValue), nil

// 	case []interface{}:
// 		var listElements []attr.Value
// 		for key, val := range v {
// 			attrValue, err := convertAnyToAttrValue(ctx, val)
// 			if err != nil {
// 				return nil, fmt.Errorf("error in nested list index %d: %w", key, err)
// 			}
// 			listElements = append(listElements, attrValue)
// 		}
// 		listValue, _ := types.ListValue(types.DynamicType, listElements)
// 		return types.DynamicValue(listValue), nil

// 	default:
// 		jsonBytes, err := json.Marshal(v)
// 		if err != nil {
// 			return nil, fmt.Errorf("unsupported type %T: %v", v, err)
// 		}
// 		return types.DynamicValue(types.StringValue(string(jsonBytes))), nil
// 	}
// }

// // Convert API map[string]interface{} → Terraform Dynamic
// // Used in: Read operations across all resources
// func convertAPIToDynamic(ctx context.Context, apiMap map[string]interface{}) (types.Dynamic, diag.Diagnostics) {
// 	var diags diag.Diagnostics

// 	if apiMap == nil || len(apiMap) == 0 {
// 		return types.DynamicNull(), diags
// 	}

// 	elements := make(map[string]attr.Value, len(apiMap))

// 	for key, value := range apiMap {
// 		attrValue, err := convertAnyToAttrValue(ctx, value)
// 		if err != nil {
// 			diags.AddError(
// 				"Error Converting API Map to Dynamic",
// 				fmt.Sprintf("Error converting API map value for key %q: %v", key, err),
// 			)
// 			return types.DynamicNull(), diags
// 		}
// 		elements[key] = attrValue
// 	}

// 	if diags.HasError() {
// 		return types.DynamicNull(), diags
// 	}

// 	mapValue, mapDiags := types.MapValue(types.DynamicType, elements)
// 	diags.Append(mapDiags...)
// 	if diags.HasError() {
// 		return types.DynamicNull(), diags
// 	}
// 	return types.DynamicValue(mapValue), diags

// }

// func convertAnyToAttrValue(ctx context.Context, value any) (attr.Value, error) {
// 	if value == nil {
// 		return types.DynamicNull(), nil
// 	}

// 	switch v := value.(type) {
// 	case string:
// 		return types.StringValue(v), nil
// 	case bool:
// 		return types.BoolValue(v), nil
// 	case int:
// 		return types.Int64Value(int64(v)), nil
// 	case int32:
// 		return types.Int32Value(v), nil
// 	case int64:
// 		return types.Int64Value(v), nil
// 	case float32:
// 		return types.Float64Value(float64(v)), nil
// 	case float64:
// 		return types.Float64Value(v), nil
// 	case map[string]any:
// 		elements := make(map[string]attr.Value, len(v))
// 		for key, val := range v {
// 			attrValue, err := convertAnyToAttrValue(ctx, val)
// 			if err != nil {
// 				return nil, fmt.Errorf("error in nested map key '%s': %w", key, err)
// 			}
// 			elements[key] = attrValue
// 		}
// 		mapValue, _ := types.MapValue(types.DynamicType, elements)
// 		return types.DynamicValue(mapValue), nil
// 	case []any:
// 		var listElements []attr.Value
// 		for key, val := range v {
// 			attrValue, err := convertAnyToAttrValue(ctx, val)
// 			if err != nil {
// 				return nil, fmt.Errorf("error in nested list index %d: %w", key, err)
// 			}
// 			listElements = append(listElements, attrValue)
// 		}
// 		listValue, _ := types.ListValue(types.DynamicType, listElements)
// 		return types.DynamicValue(listValue), nil
// 	default:
// 		jsonBytes, err := json.Marshal(v)
// 		if err != nil {
// 			return nil, fmt.Errorf("unsupported type %T: %v", v, err)
// 		}
// 		return types.DynamicValue(types.StringValue(string(jsonBytes))), nil
// 	}
// }

// // Convert Terraform Dynamic → API map[string]interface{}
// // Used in: Create/Update operations across all resources
// // func convertDynamicToAPI(ctx context.Context, dynValue types.Dynamic) (map[string]interface{}, error) {
// // 	if dynValue.IsNull() || dynValue.IsUnknown() {
// // 		return nil, nil
// // 	}

// // 	underlyingValue := dynValue.UnderlyingValue()

// // 	// Try to get the value as a map[string]interface{}
// // 	mapValue, ok := underlyingValue.(types.Map)
// // 	if !ok {
// // 		return nil, fmt.Errorf("expected Dynamic underlying value to be a Map, got %T", underlyingValue)
// // 	}

// // 	result := make(map[string]interface{}, len(mapValue.Elements()))

// // 	for key, value := range mapValue.Elements() {
// // 		apiValue, err := attrValueToInterface(ctx, value)
// // 		if err != nil {
// // 			return nil, fmt.Errorf("converting Dynamic map value for key %q: %w", key, err)
// // 		}
// // 		result[key] = apiValue
// // 	}
// // 	return result, nil
// // }

// // // func attrValueToInterface(ctx context.Context, value attr.Value) (interface{}, error) {
// // // 	if value.IsNull() || value.IsUnknown() {
// // // 		return nil, nil
// // // 	}

// // // 	if dynValue, ok := value.(types.Dynamic); ok {
// // // 		value = dynValue.UnderlyingValue()
// // // 	}

// // // 	switch v := value.(type) {
// // // 	case types.String:
// // // 		return v.ValueString(), nil
// // // 	case types.Bool:
// // // 		return v.ValueBool(), nil
// // // 	case types.Int64:
// // // 		return v.ValueInt64(), nil
// // // 	case types.Int32:
// // // 		return v.ValueInt32(), nil
// // // 	case types.Float64:
// // // 		return v.ValueFloat64(), nil
// // // 	case types.Float32:
// // // 		return v.ValueFloat32(), nil

// // // 	case types.Number:
// // // 		bf := v.ValueBigFloat()
// // // 		if bf.IsInt() {
// // // 			i, _ := bf.Int64()
// // // 			return i, nil
// // // 		}
// // // 		f, _ := bf.Float64()
// // // 		return f, nil

// // // 	case types.List:
// // // 		var elements []interface{}
// // // 		for _, elem := range v.Elements() {
// // // 			converted, err := attrValueToInterface(ctx, elem)
// // // 			if err != nil {
// // // 				return nil, err
// // // 			}
// // // 			elements = append(elements, converted)
// // // 		}
// // // 		return elements, nil

// // // 	case types.Map:
// // // 		result := make(map[string]interface{})
// // // 		for key, val := range v.Elements() {
// // // 			converted, err := attrValueToInterface(ctx, val)
// // // 			if err != nil {
// // // 				return nil, fmt.Errorf("error in nested map key '%s': %w", key, err)
// // // 			}
// // // 			result[key] = converted
// // // 		}
// // // 		return result, nil

// // // 	case types.Object:
// // // 		result := make(map[string]interface{})
// // // 		for key, val := range v.Attributes() {
// // // 			converted, err := attrValueToInterface(ctx, val)
// // // 			if err != nil {
// // // 				return nil, fmt.Errorf("error in object attribute '%s': %w", key, err)
// // // 			}
// // // 			result[key] = converted
// // // 		}
// // // 		return result, nil

// // // 	default:
// // // 		return nil, fmt.Errorf("unsupported type: %T", v)
// // // 	}
// // // }
