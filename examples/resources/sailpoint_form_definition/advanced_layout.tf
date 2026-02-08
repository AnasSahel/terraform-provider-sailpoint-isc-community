# Advanced form with multi-column layout and various element types
resource "sailpoint_form_definition" "advanced" {
  name        = "Advanced Layout Form"
  description = "Form demonstrating various element types and layouts"

  owner = {
    type = "IDENTITY"
    id   = "2c91808a7813090a017814121e121518"
  }

  form_elements = jsonencode([
    # Header section with description
    {
      id          = "section-header"
      elementType = "SECTION"
      key         = "header"
      config = {
        label = "User Registration"
      }
    },
    {
      id          = "description-text"
      elementType = "DESCRIPTION"
      key         = "intro"
      config = {
        text = "Please fill out all required fields to complete your registration."
      }
    },
    # Personal Information Section
    {
      id          = "section-personal"
      elementType = "SECTION"
      key         = "personal-info"
      config = {
        label = "Personal Information"
      }
    },
    # Column set for side-by-side fields
    {
      id          = "columns-name"
      elementType = "COLUMN_SET"
      key         = "name-columns"
      config = {
        columns = 2
      }
    },
    {
      id          = "first-name"
      elementType = "TEXT"
      key         = "firstName"
      config = {
        label       = "First Name"
        placeholder = "John"
        columnIndex = 0
      }
      validations = [
        { validationType = "REQUIRED" }
      ]
    },
    {
      id          = "last-name"
      elementType = "TEXT"
      key         = "lastName"
      config = {
        label       = "Last Name"
        placeholder = "Doe"
        columnIndex = 1
      }
      validations = [
        { validationType = "REQUIRED" }
      ]
    },
    # Contact Information
    {
      id          = "section-contact"
      elementType = "SECTION"
      key         = "contact-info"
      config = {
        label = "Contact Information"
      }
    },
    {
      id          = "email"
      elementType = "EMAIL"
      key         = "email"
      config = {
        label       = "Email Address"
        placeholder = "user@example.com"
      }
      validations = [
        { validationType = "REQUIRED" },
        { validationType = "EMAIL" }
      ]
    },
    {
      id          = "phone"
      elementType = "PHONE"
      key         = "phoneNumber"
      config = {
        label       = "Phone Number"
        placeholder = "+1 (555) 123-4567"
      }
      validations = [
        { validationType = "PHONE" }
      ]
    },
    # Preferences Section
    {
      id          = "section-preferences"
      elementType = "SECTION"
      key         = "preferences"
      config = {
        label = "Preferences"
      }
    },
    {
      id          = "notifications"
      elementType = "TOGGLE"
      key         = "enableNotifications"
      config = {
        label       = "Enable Email Notifications"
        description = "Receive updates about your account"
        default     = true
      }
    },
    {
      id          = "timezone"
      elementType = "SELECT"
      key         = "timezone"
      config = {
        label   = "Timezone"
        options = ["UTC", "America/New_York", "America/Los_Angeles", "Europe/London", "Asia/Tokyo"]
      }
    },
    # Hidden field for tracking
    {
      id          = "tracking-id"
      elementType = "HIDDEN"
      key         = "trackingId"
      config = {
        defaultValue = "web-registration-v2"
      }
    },
    # Additional Notes
    {
      id          = "section-notes"
      elementType = "SECTION"
      key         = "notes"
      config = {
        label = "Additional Information"
      }
    },
    {
      id          = "notes"
      elementType = "TEXTAREA"
      key         = "additionalNotes"
      config = {
        label       = "Notes"
        placeholder = "Any additional information you'd like to provide"
        maxLength   = 500
      }
    }
  ])

  form_conditions = jsonencode([
    # If notifications are disabled, hide timezone selection
    {
      ruleOperator = "AND"
      rules = [
        {
          sourceType = "ELEMENT"
          source     = "enableNotifications"
          operator   = "EQ"
          valueType  = "BOOLEAN"
          value      = "false"
        }
      ]
      effects = [
        {
          effectType = "HIDE"
          config = {
            element = "timezone"
          }
        }
      ]
    }
  ])
}
