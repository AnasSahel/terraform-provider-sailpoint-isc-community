package lifecycle_state

import (
	"context"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sailpoint-oss/golang-sdk/v2/api_v2025"
)

func FromCreatePlanToSailPointLifecycleState(ctx context.Context, plan *LifecycleStateResourceModel) api_v2025.LifecycleState {
	lifecycleState := api_v2025.NewLifecycleStateWithDefaults()

	// Required fields
	lifecycleState.SetName(plan.Name.ValueString())
	lifecycleState.SetTechnicalName(plan.TechnicalName.ValueString())

	// Optional fields using helper functions
	utils.SetStringIfPresent(plan.Description, lifecycleState.SetDescription)
	utils.SetBoolIfPresent(plan.Enabled, lifecycleState.SetEnabled)
	utils.SetStringIfPresent(plan.IdentityState, lifecycleState.SetIdentityState)
	utils.SetInt32IfPresent(plan.Priority, lifecycleState.SetPriority)
	utils.SetStringListIfPresent(ctx, plan.AccessProfileIds, lifecycleState.SetAccessProfileIds)

	// Handle nested configuration
	if plan.AccessActionConfiguration != nil {
		accessActionConfiguration := api_v2025.NewAccessActionConfigurationWithDefaults()
		utils.SetBoolIfPresent(plan.AccessActionConfiguration.RemoveAllAccessEnabled, accessActionConfiguration.SetRemoveAllAccessEnabled)
		lifecycleState.SetAccessActionConfiguration(*accessActionConfiguration)
	}

	// Handle email notification option
	if plan.EmailNotificationOption != nil {
		emailNotificationOption := api_v2025.NewEmailNotificationOptionWithDefaults()
		utils.SetBoolIfPresent(plan.EmailNotificationOption.NotifyManagers, emailNotificationOption.SetNotifyManagers)
		utils.SetBoolIfPresent(plan.EmailNotificationOption.NotifyAllAdmins, emailNotificationOption.SetNotifyAllAdmins)
		utils.SetBoolIfPresent(plan.EmailNotificationOption.NotifySpecificUsers, emailNotificationOption.SetNotifySpecificUsers)
		utils.SetStringListIfPresent(ctx, plan.EmailNotificationOption.EmailAddressList, emailNotificationOption.SetEmailAddressList)
		lifecycleState.SetEmailNotificationOption(*emailNotificationOption)
	}

	// Handle account actions
	if len(plan.AccountActions) > 0 {
		accountActionsList := make([]api_v2025.AccountAction, 0, len(plan.AccountActions))
		for _, actionModel := range plan.AccountActions {
			accountAction := api_v2025.NewAccountActionWithDefaults()

			// Set action (required field)
			if !actionModel.Action.IsNull() && !actionModel.Action.IsUnknown() {
				accountAction.SetAction(actionModel.Action.ValueString())
			}

			// Set optional fields using helper functions
			utils.SetBoolIfPresent(actionModel.AllSources, accountAction.SetAllSources)
			utils.SetStringListIfPresent(ctx, actionModel.SourceIds, accountAction.SetSourceIds)
			utils.SetStringListIfPresent(ctx, actionModel.ExcludeSourceIds, accountAction.SetExcludeSourceIds)

			accountActionsList = append(accountActionsList, *accountAction)
		}
		lifecycleState.SetAccountActions(accountActionsList)
	}

	return *lifecycleState
}

func ToTerraformResource(ctx context.Context, plan *LifecycleStateResourceModel, lifecycleState *api_v2025.LifecycleState) LifecycleStateResourceModel {
	// Always populate all fields from the API response

	newState := LifecycleStateResourceModel{
		LifecycleStateModel: LifecycleStateModel{
			Id:               types.StringValue(lifecycleState.GetId()),
			Name:             types.StringValue(lifecycleState.GetName()),
			TechnicalName:    types.StringValue(lifecycleState.GetTechnicalName()),
			Enabled:          types.BoolValue(lifecycleState.GetEnabled()),
			IdentityCount:    types.Int32Value(lifecycleState.GetIdentityCount()),
			AccessProfileIds: types.ListNull(types.StringType),
			Priority:         types.Int32Value(lifecycleState.GetPriority()),
			Created:          types.StringValue(lifecycleState.GetCreated().String()),
			Modified:         types.StringValue(lifecycleState.GetModified().String()),
		},
	}

	if plan != nil {
		newState.IdentityProfileId = plan.IdentityProfileId

		// Set optional fields only if they were present in the plan
		newState.Description = utils.GetStringOrNull(plan.Description, lifecycleState.GetDescription())
		newState.IdentityState = utils.GetStringOrNull(plan.IdentityState, lifecycleState.GetIdentityState())
		newState.AccessProfileIds = utils.GetListOrNull(ctx, plan.AccessProfileIds, lifecycleState.GetAccessProfileIds())

		// Only populate AccessActionConfiguration if it was specified in the plan
		if plan.AccessActionConfiguration != nil {
			if accessActionConfig, ok := lifecycleState.GetAccessActionConfigurationOk(); ok {
				newState.AccessActionConfiguration = &AccessActionConfigurationModel{
					RemoveAllAccessEnabled: types.BoolValue(accessActionConfig.GetRemoveAllAccessEnabled()),
				}
			}
		}

		if plan.EmailNotificationOption != nil {
			if emailNotificationOption, ok := lifecycleState.GetEmailNotificationOptionOk(); ok {
				emailAddressList, _ := types.ListValueFrom(ctx, types.StringType, emailNotificationOption.GetEmailAddressList())

				newState.EmailNotificationOption = &EmailNotificationOptionModel{
					NotifyManagers:      types.BoolValue(emailNotificationOption.GetNotifyManagers()),
					NotifyAllAdmins:     types.BoolValue(emailNotificationOption.GetNotifyAllAdmins()),
					NotifySpecificUsers: types.BoolValue(emailNotificationOption.GetNotifySpecificUsers()),
					EmailAddressList:    emailAddressList,
				}
			}
		}

		// Handle account actions if specified in the plan
		if len(plan.AccountActions) > 0 {
			if accountActions, ok := lifecycleState.GetAccountActionsOk(); ok {
				newState.AccountActions = make([]AccountActionModel, 0, len(accountActions))
				for _, action := range accountActions {
					sourceIds, _ := types.ListValueFrom(ctx, types.StringType, action.GetSourceIds())
					excludeSourceIds, _ := types.ListValueFrom(ctx, types.StringType, action.GetExcludeSourceIds())

					actionModel := AccountActionModel{
						Action:           types.StringValue(action.GetAction()),
						SourceIds:        sourceIds,
						ExcludeSourceIds: excludeSourceIds,
						AllSources:       types.BoolValue(action.GetAllSources()),
					}
					newState.AccountActions = append(newState.AccountActions, actionModel)
				}
			}
		}
	}

	return newState
}

// Datasource part
func ToTerraformDataSource(ctx context.Context, lifecycleState api_v2025.LifecycleState) LifecycleStateModel {
	accessProfileIds, _ := types.ListValueFrom(ctx, types.StringType, lifecycleState.GetAccessProfileIds())

	m := LifecycleStateModel{
		Id:               types.StringValue(lifecycleState.GetId()),
		Name:             types.StringValue(lifecycleState.GetName()),
		Enabled:          types.BoolValue(lifecycleState.GetEnabled()),
		TechnicalName:    types.StringValue(lifecycleState.GetTechnicalName()),
		Description:      types.StringValue(lifecycleState.GetDescription()),
		IdentityCount:    types.Int32Value(lifecycleState.GetIdentityCount()),
		AccessProfileIds: accessProfileIds,
		IdentityState:    types.StringValue(lifecycleState.GetIdentityState()),
		Priority:         types.Int32Value(lifecycleState.GetPriority()),

		Created:  types.StringValue(lifecycleState.GetCreated().String()),
		Modified: types.StringValue(lifecycleState.GetModified().String()),
	}

	accessActionConfiguration, ok := lifecycleState.GetAccessActionConfigurationOk()
	if ok && accessActionConfiguration != nil {
		m.AccessActionConfiguration = &AccessActionConfigurationModel{
			RemoveAllAccessEnabled: types.BoolValue(accessActionConfiguration.GetRemoveAllAccessEnabled()),
		}
	}

	emailNotificationOption, ok := lifecycleState.GetEmailNotificationOptionOk()
	if ok && emailNotificationOption != nil {
		emailAddressList, _ := types.ListValueFrom(ctx, types.StringType, emailNotificationOption.GetEmailAddressList())

		m.EmailNotificationOption = &EmailNotificationOptionModel{
			NotifyManagers:      types.BoolValue(emailNotificationOption.GetNotifyManagers()),
			NotifyAllAdmins:     types.BoolValue(emailNotificationOption.GetNotifyAllAdmins()),
			NotifySpecificUsers: types.BoolValue(emailNotificationOption.GetNotifySpecificUsers()),
			EmailAddressList:    emailAddressList,
		}
	}

	// Handle account actions
	accountActions, ok := lifecycleState.GetAccountActionsOk()
	if ok && accountActions != nil {
		m.AccountActions = make([]AccountActionModel, 0, len(accountActions))
		for _, action := range accountActions {
			sourceIds, _ := types.ListValueFrom(ctx, types.StringType, action.GetSourceIds())
			excludeSourceIds, _ := types.ListValueFrom(ctx, types.StringType, action.GetExcludeSourceIds())

			actionModel := AccountActionModel{
				Action:           types.StringValue(action.GetAction()),
				SourceIds:        sourceIds,
				ExcludeSourceIds: excludeSourceIds,
				AllSources:       types.BoolValue(action.GetAllSources()),
			}
			m.AccountActions = append(m.AccountActions, actionModel)
		}
	}

	return m
}
