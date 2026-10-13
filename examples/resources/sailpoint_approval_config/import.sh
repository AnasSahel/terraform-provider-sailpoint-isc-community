#!/bin/bash
# Import format: <scope>/<object_id>. Scope is ROLE, ACCESS_PROFILE, ENTITLEMENT or TENANT.
# Importing an object that only inherits the tenant config fails with "Cannot import non-existent remote object".
terraform import sailpoint_approval_config.finance_role "ROLE/REPLACE_WITH_ROLE_ID"
