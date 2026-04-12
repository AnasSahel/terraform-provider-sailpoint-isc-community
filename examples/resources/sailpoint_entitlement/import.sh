#!/bin/bash
# Import an existing entitlement by its ID.
# Note: terraform destroy is a no-op for this resource — entitlements cannot be deleted via the API.
terraform import sailpoint_entitlement.admin_group "REPLACE_WITH_ENTITLEMENT_ID"
