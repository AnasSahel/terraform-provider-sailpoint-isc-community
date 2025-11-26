#!/bin/bash

# Import an existing SailPoint Access Profile into Terraform state
# The import ID is the access profile's UUID from SailPoint ISC

# Basic import command
terraform import sailpoint_access_profile.example "00000000000000000000000000000001"

# Import with specific resource name
terraform import sailpoint_access_profile.my_profile "00000000000000000000000000000002"

# Steps to find the access profile ID:
# 1. Via UI: Navigate to Access Profiles, select the profile, and copy the ID from the URL
# 2. Via API: Use the List Access Profiles endpoint
#    curl -H "Authorization: Bearer $TOKEN" \
#         "https://{tenant}.api.identitynow.com/v2025/access-profiles?filters=name eq \"Profile Name\""
# 3. Via ISC CLI: Use the SailPoint CLI to list access profiles
#    sail access-profile list --name "Profile Name"
