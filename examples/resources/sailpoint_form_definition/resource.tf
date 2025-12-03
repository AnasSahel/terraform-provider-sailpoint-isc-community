# SailPoint Form Definition Resource Examples
#
# IMPORTANT NOTES:
# - REQUIRED FIELDS: name, owner (with type and id), form_elements
# - Forms are composed of sections and fields for data collection
# - Form elements are now structured objects (not JSON strings)
# - Each form element MUST have an 'id' and 'element_type'
# - The 'config' field within elements uses jsonencode() for complex configurations
# - The 'owner' field references an identity who owns this form
#
# For more information, see:
# https://developer.sailpoint.com/docs/api/v2025/custom-forms/

# Example 1: Basic Form with Single Section
resource "sailpoint_form_definition" "basic_form" {
  name        = "Employee Information Form"
  description = "Basic form to collect employee information"

  owner = {
    type = "IDENTITY"
    id   = "00000000000000000000000000000001"
    name = "John Doe"
  }

  form_elements = [
    {
      id           = "section1"
      element_type = "SECTION"
      config = jsonencode({
        label = "Personal Information"
        formElements = [
          {
            id          = "firstName"
            elementType = "TEXT"
            key         = "firstName"
            config = {
              label = "First Name"
            }
          },
          {
            id          = "lastName"
            elementType = "TEXT"
            key         = "lastName"
            config = {
              label = "Last Name"
            }
          }
        ]
      })
    }
  ]
}

# Example 2: Form with Multiple Sections and Validations
resource "sailpoint_form_definition" "advanced_form" {
  name        = "Access Request Form"
  description = "Comprehensive form for requesting system access"

  owner = {
    type = "IDENTITY"
    id   = "00000000000000000000000000000001"
    name = "Admin User"
  }

  form_elements = [
    {
      id           = "section1"
      element_type = "SECTION"
      config = jsonencode({
        label = "Requestor Information"
        formElements = [
          {
            id          = "email"
            elementType = "TEXT"
            key         = "requesterEmail"
            config = {
              label = "Email Address"
            }
          },
          {
            id          = "phone"
            elementType = "TEXT"
            key         = "requesterPhone"
            config = {
              label = "Phone Number"
            }
          }
        ]
      })
    },
    {
      id           = "section2"
      element_type = "SECTION"
      validations  = []
      config = jsonencode({
        label = "Access Details"
        formElements = [
          {
            id          = "accessType"
            elementType = "SELECT"
            key         = "accessType"
            config = {
              label = "Access Type"
              dataSource = {
                dataSourceType = "STATIC"
                config = {
                  options = [
                    { label = "Read Only", value = "read" },
                    { label = "Read/Write", value = "write" },
                    { label = "Admin", value = "admin" }
                  ]
                }
              }
            }
          },
          {
            id          = "startDate"
            elementType = "DATE"
            key         = "startDate"
            config = {
              label = "Start Date"
            }
          },
          {
            id          = "justification"
            elementType = "TEXTAREA"
            key         = "justification"
            config = {
              label       = "Business Justification"
              placeholder = "Please provide a business justification for this access request"
            }
          }
        ]
      })
    }
  ]

  form_input = [
    {
      id    = "requesterIdentity"
      type  = "IDENTITY"
      label = "Requester Identity"
    }
  ]
}

# Example 3: Form with Conditional Logic
resource "sailpoint_form_definition" "conditional_form" {
  name        = "Account Request with Conditions"
  description = "Form with conditional fields based on user selections"

  owner = {
    type = "IDENTITY"
    id   = "00000000000000000000000000000001"
    name = "Form Administrator"
  }

  form_elements = [
    {
      id           = "section1"
      element_type = "SECTION"
      config = jsonencode({
        label = "Account Details"
        formElements = [
          {
            id          = "accountType"
            elementType = "SELECT"
            key         = "accountType"
            config = {
              label = "Account Type"
              dataSource = {
                dataSourceType = "STATIC"
                config = {
                  options = [
                    { label = "Standard", value = "standard" },
                    { label = "Privileged", value = "privileged" }
                  ]
                }
              }
            }
          },
          {
            id          = "approverField"
            elementType = "TEXT"
            key         = "approver"
            config = {
              label = "Manager Approval Required"
            }
          }
        ]
      })
    }
  ]

  form_conditions = [
    {
      rule_operator = "AND"
      rules = [
        {
          source_type = "ELEMENT"
          source      = "accountType"
          operator    = "EQ"
          value_type  = "STRING"
          value       = "privileged"
        }
      ]
      effects = [
        {
          effect_type = "SHOW"
          config = {
            element = "approverField"
          }
        }
      ]
    }
  ]
}

# Example 4: Form with Required Validations
resource "sailpoint_form_definition" "onboarding_form" {
  name        = "New Hire Onboarding"
  description = "Collect information for new employee onboarding"

  owner = {
    type = "IDENTITY"
    id   = "00000000000000000000000000000001"
  }

  form_elements = [
    {
      id           = "section1"
      element_type = "SECTION"
      validations = [
        {
          validation_type = "REQUIRED"
        }
      ]
      config = jsonencode({
        label = "New Hire Information"
        formElements = [
          {
            id          = "hireDate"
            elementType = "DATE"
            key         = "hireDate"
            config = {
              label    = "Start Date"
              required = true
            }
          },
          {
            id          = "department"
            elementType = "TEXT"
            key         = "department"
            config = {
              label    = "Department"
              required = true
            }
          },
          {
            id          = "equipmentNeeded"
            elementType = "TOGGLE"
            key         = "equipmentNeeded"
            config = {
              label = "Equipment Required"
            }
          }
        ]
      })
    }
  ]
}
