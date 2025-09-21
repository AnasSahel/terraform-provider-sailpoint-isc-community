#!/bin/bash

# Import an existing SailPoint source into Terraform state
# Usage: ./import.sh <source_id> <resource_name>

if [ $# -ne 2 ]; then
    echo "Usage: $0 <source_id> <resource_name>"
    echo "Example: $0 2c91808570313110017040b06f344ec9 sailpoint_source.imported_ad"
    exit 1
fi

SOURCE_ID="$1"
RESOURCE_NAME="$2"

echo "Importing SailPoint source with ID: $SOURCE_ID"
echo "Into Terraform resource: $RESOURCE_NAME"

# Import the source
terraform import "$RESOURCE_NAME" "$SOURCE_ID"

if [ $? -eq 0 ]; then
    echo "Successfully imported source!"
    echo ""
    echo "Next steps:"
    echo "1. Run 'terraform plan' to see the configuration diff"
    echo "2. Update your .tf file to match the imported state"
    echo "3. Run 'terraform plan' again to verify no changes needed"
    echo ""
    echo "Example configuration for the imported source:"
    echo "resource \"sailpoint_source\" \"imported_ad\" {"
    echo "  # Add the configuration based on terraform plan output"
    echo "  name        = \"...\""
    echo "  description = \"...\""
    echo "  connector   = \"...\""
    echo "  owner       = jsonencode({...})"
    echo "  configuration = {...}"
    echo "}"
else
    echo "Import failed! Please check:"
    echo "1. Source ID is correct and exists in SailPoint"
    echo "2. SailPoint credentials are properly configured"
    echo "3. Network connectivity to SailPoint tenant"
fi