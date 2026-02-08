# Access request form with multiple fields and validation
resource "sailpoint_form_definition" "access_request" {
  name        = "Application Access Request"
  description = "Form for requesting access to applications"

  owner = {
    type = "IDENTITY"
    id   = "2c91808a7813090a017814121e121518"
  }

  form_input = jsonencode([
    {
      id          = "requestedFor"
      type        = "STRING"
      label       = "Requested For"
      description = "The identity requesting access"
    },
    {
      id          = "applicationName"
      type        = "STRING"
      label       = "Application"
      description = "The application being requested"
    }
  ])

  form_elements = jsonencode([
    {
      id          = "section-details"
      elementType = "SECTION"
      key         = "request-details"
      config = {
        label       = "Request Details"
        description = "Please provide details about your access request"
      }
    },
    {
      id          = "access-level"
      elementType = "SELECT"
      key         = "accessLevel"
      config = {
        label   = "Access Level"
        options = ["Read Only", "Read/Write", "Admin"]
      }
      validations = [
        { validationType = "REQUIRED" }
      ]
    },
    {
      id          = "justification"
      elementType = "TEXTAREA"
      key         = "businessJustification"
      config = {
        label       = "Business Justification"
        placeholder = "Explain why you need this access"
        maxLength   = 1000
      }
      validations = [
        { validationType = "REQUIRED" },
        { validationType = "MIN_LENGTH", minLength = 20 }
      ]
    },
    {
      id          = "duration"
      elementType = "SELECT"
      key         = "accessDuration"
      config = {
        label   = "Access Duration"
        options = ["30 Days", "90 Days", "1 Year", "Permanent"]
      }
      validations = [
        { validationType = "REQUIRED" }
      ]
    },
    {
      id          = "start-date"
      elementType = "DATE"
      key         = "startDate"
      config = {
        label = "Requested Start Date"
      }
    }
  ])

  form_conditions = jsonencode([
    {
      ruleOperator = "AND"
      rules = [
        {
          sourceType = "ELEMENT"
          source     = "accessLevel"
          operator   = "EQ"
          valueType  = "STRING"
          value      = "Admin"
        }
      ]
      effects = [
        {
          effectType = "REQUIRE"
          config = {
            element = "justification"
          }
        }
      ]
    }
  ])
}
