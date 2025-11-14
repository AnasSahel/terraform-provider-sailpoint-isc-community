# SailPoint Form Definition Resource Examples
#
# IMPORTANT NOTES:
# - Forms are composed of sections and fields for data collection
# - The 'name' field is required
# - Form elements, inputs, and conditions are specified as JSON strings
# - The 'owner' field references an identity who owns this form
#
# For more information, see:
# https://developer.sailpoint.com/docs/api/v2025/custom-forms/

# Example 1: Basic Form with Single Section and Text Field
resource "sailpoint_form_definition" "basic_form" {
  name        = "Employee Information Form"
  description = "Basic form to collect employee information"

  owner = {
    type = "IDENTITY"
    id   = "2c9180867624cbd7017642d8c8c81f67"
    name = "John Doe"
  }

  form_elements = jsonencode([
    {
      id          = "section1"
      elementType = "SECTION"
      config = {
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
      }
    }
  ])
}

# Example 2: Form with Multiple Sections and Field Types
resource "sailpoint_form_definition" "advanced_form" {
  name        = "Access Request Form"
  description = "Comprehensive form for requesting system access"

  owner = {
    type = "IDENTITY"
    id   = "2c9180867624cbd7017642d8c8c81f67"
    name = "Admin User"
  }

  form_elements = jsonencode([
    {
      id          = "section1"
      elementType = "SECTION"
      config = {
        label = "Requestor Information"
        formElements = [
          {
            id          = "email"
            elementType = "EMAIL"
            key         = "requesterEmail"
            config = {
              label = "Email Address"
            }
          },
          {
            id          = "phone"
            elementType = "PHONE"
            key         = "requesterPhone"
            config = {
              label = "Phone Number"
            }
          }
        ]
      }
    },
    {
      id          = "section2"
      elementType = "SECTION"
      config = {
        label = "Access Details"
        formElements = [
          {
            id          = "accessType"
            elementType = "SELECT"
            key         = "accessType"
            config = {
              label = "Access Type"
              options = [
                { label = "Read Only", value = "read" },
                { label = "Read/Write", value = "write" },
                { label = "Admin", value = "admin" }
              ]
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
      }
    }
  ])

  form_input = jsonencode([
    {
      id    = "requesterIdentity"
      type  = "IDENTITY"
      label = "Requester Identity"
    }
  ])
}

# Example 3: Form with Conditional Logic
resource "sailpoint_form_definition" "conditional_form" {
  name        = "Account Request with Conditions"
  description = "Form with conditional fields based on user selections"

  owner = {
    type = "IDENTITY"
    id   = "2c9180867624cbd7017642d8c8c81f67"
    name = "Form Administrator"
  }

  form_elements = jsonencode([
    {
      id          = "section1"
      elementType = "SECTION"
      config = {
        label = "Account Details"
        formElements = [
          {
            id          = "accountType"
            elementType = "SELECT"
            key         = "accountType"
            config = {
              label = "Account Type"
              options = [
                { label = "Standard", value = "standard" },
                { label = "Privileged", value = "privileged" }
              ]
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
      }
    }
  ])

  form_conditions = jsonencode([
    {
      ruleOperator = "AND"
      rules = [
        {
          sourceType = "ELEMENT"
          source     = "accountType"
          operator   = "EQ"
          valueType  = "STRING"
          value      = "privileged"
        }
      ]
      effects = [
        {
          effectType = "SHOW"
          config = {
            element = "approverField"
          }
        }
      ]
    }
  ])
}

# Example 4: Simple Onboarding Form
resource "sailpoint_form_definition" "onboarding_form" {
  name        = "New Hire Onboarding"
  description = "Collect information for new employee onboarding"

  form_elements = jsonencode([
    {
      id          = "section1"
      elementType = "SECTION"
      config = {
        label = "New Hire Information"
        formElements = [
          {
            id          = "hireDate"
            elementType = "DATE"
            key         = "hireDate"
            config = {
              label = "Start Date"
            }
          },
          {
            id          = "department"
            elementType = "TEXT"
            key         = "department"
            config = {
              label = "Department"
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
      }
    }
  ])
}
