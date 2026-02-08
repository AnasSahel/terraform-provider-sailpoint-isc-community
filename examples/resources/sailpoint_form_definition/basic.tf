# Basic form definition with a simple text field
resource "sailpoint_form_definition" "basic" {
  name        = "Basic Form"
  description = "A simple form with one text field"

  owner = {
    type = "IDENTITY"
    id   = "2c91808a7813090a017814121e121518"
  }

  form_elements = jsonencode([
    {
      id          = "section-1"
      elementType = "SECTION"
      key         = "main-section"
      config = {
        label = "Information"
      }
    },
    {
      id          = "name-field"
      elementType = "TEXT"
      key         = "fullName"
      config = {
        label       = "Full Name"
        placeholder = "Enter your full name"
      }
      validations = [
        { validationType = "REQUIRED" }
      ]
    }
  ])
}
