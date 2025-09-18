#!/bin/bash

# Import script for SailPoint managed cluster resource
# This script demonstrates how to import an existing managed cluster into Terraform

# Set your cluster ID - replace with actual cluster ID from SailPoint ISC
CLUSTER_ID="2c918085-74f3-4b96-8c31-3c3a7cb8f5e2"

# Set the Terraform resource name you want to use
RESOURCE_NAME="sailpoint_managed_cluster.imported_cluster"

echo "Importing SailPoint managed cluster..."
echo "Cluster ID: ${CLUSTER_ID}"
echo "Resource Name: ${RESOURCE_NAME}"

# Run the import command
terraform import "${RESOURCE_NAME}" "${CLUSTER_ID}"

if [ $? -eq 0 ]; then
    echo "✅ Successfully imported managed cluster!"
    echo ""
    echo "Next steps:"
    echo "1. Run 'terraform show' to see the imported state"
    echo "2. Create a corresponding resource block in your .tf file"
    echo "3. Run 'terraform plan' to see any differences"
    echo "4. Update your resource configuration to match the imported state"
    echo "5. Run 'terraform plan' again to ensure no changes are detected"
else
    echo "❌ Failed to import managed cluster"
    echo "Please check:"
    echo "- The cluster ID is correct and exists in SailPoint"
    echo "- Your SailPoint credentials are properly configured"
    echo "- The Terraform resource name doesn't already exist in state"
fi