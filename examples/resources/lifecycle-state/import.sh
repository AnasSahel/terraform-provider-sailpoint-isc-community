#!/bin/bash

# Import existing lifecycle state
# Usage: ./import.sh <identity_profile_id> <lifecycle_state_id>

if [ $# -ne 2 ]; then
    echo "Usage: $0 <identity_profile_id> <lifecycle_state_id>"
    echo "Example: $0 55ecd185917d4b2e9f6d42d23656fdcb 6e629eac97ed4e08a35bcb66a301806e"
    exit 1
fi

IDENTITY_PROFILE_ID=$1
LIFECYCLE_STATE_ID=$2
IMPORT_ID="${IDENTITY_PROFILE_ID}:${LIFECYCLE_STATE_ID}"

echo "Importing lifecycle state with ID: ${IMPORT_ID}"
terraform import sailpoint_lifecycle_state.example "${IMPORT_ID}"