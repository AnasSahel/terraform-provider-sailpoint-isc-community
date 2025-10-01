#!/bin/bash

# Import an existing identity attribute by name
terraform import sailpoint_identity_attribute.cost_center "costCenter"

# Import another identity attribute
terraform import sailpoint_identity_attribute.department "department"

echo "Identity attributes imported successfully!"
echo "Run 'terraform plan' to see any configuration drift."