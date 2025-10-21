package models

import (
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/client"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ObjectRef struct {
	Type types.String `tfsdk:"type"`
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func NewObjectRefFromAPI(apiRef *client.ObjectRef) *ObjectRef {
	if apiRef == nil {
		return nil
	}
	return &ObjectRef{
		Type: types.StringValue(apiRef.Type),
		ID:   types.StringValue(apiRef.ID),
		Name: types.StringValue(apiRef.Name),
	}
}
