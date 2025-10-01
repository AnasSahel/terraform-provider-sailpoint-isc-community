#!/bin/bash

# SailPoint ISC Identity Attribute Import Script
# This script demonstrates how to import existing identity attributes into Terraform management

echo "🚀 Starting SailPoint Identity Attribute Import Process..."
echo ""

# Function to import an attribute and handle errors
import_attribute() {
    local resource_name=$1
    local attribute_name=$2
    
    echo "📥 Importing identity attribute: $attribute_name"
    echo "   Resource name: sailpoint_identity_attribute.$resource_name"
    
    if terraform import "sailpoint_identity_attribute.$resource_name" "$attribute_name"; then
        echo "✅ Successfully imported $attribute_name"
    else
        echo "❌ Failed to import $attribute_name"
        echo "   - Check that the attribute exists in SailPoint"
        echo "   - Verify the attribute name is correct (case-sensitive)"
        echo "   - Ensure it's not a system attribute"
    fi
    echo ""
}

# Import common business attributes
echo "📋 Importing common business identity attributes..."
import_attribute "cost_center" "costCenter"
import_attribute "department" "department"
import_attribute "employee_id" "employeeId"
import_attribute "manager" "manager"
import_attribute "location" "location"

# Import standard SailPoint attributes (if they exist as custom)
echo "📋 Importing standard identity attributes..."
import_attribute "job_title" "jobTitle"
import_attribute "division" "division"

# Import security-related attributes
echo "🔒 Importing security-related identity attributes..."
import_attribute "security_clearance" "securityClearance"
import_attribute "access_level" "accessLevel"

echo "🎉 Import process completed!"
echo ""
echo "📊 Next steps:"
echo "1. Run 'terraform plan' to see any configuration drift"
echo "2. Update your .tf files to match the imported state"
echo "3. Run 'terraform plan' again to verify no changes are needed"
echo "4. Consider adding lifecycle rules for sensitive attributes"
echo ""
echo "💡 Tips:"
echo "- Use 'terraform show' to see the current state of imported resources"
echo "- Use 'terraform state list' to see all managed resources"
echo "- Use 'terraform state show <resource>' to see detailed resource state"