# Example 1: Basic Upper Transform - Convert input to uppercase
resource "sailpoint_transform" "upper_example" {
  name = "Upper Case Transform"
  type = "upper"

  attributes = jsonencode({
    input = {
      type = "accountAttribute"
      attributes = {
        attributeName = "firstName"
        sourceName    = "My Source"
      }
    }
  })
}

# Example 2: Lower Case Transform - Convert input to lowercase
resource "sailpoint_transform" "lower_example" {
  name = "Lower Case Transform"
  type = "lower"

  attributes = jsonencode({
    input = {
      type = "accountAttribute"
      attributes = {
        attributeName = "lastName"
        sourceName    = "My Source"
      }
    }
  })
}

# Example 3: Substring Transform - Extract part of a string
resource "sailpoint_transform" "substring_example" {
  name = "Email Domain Extract"
  type = "substring"

  attributes = jsonencode({
    input = {
      type = "accountAttribute"
      attributes = {
        attributeName = "email"
        sourceName    = "My Source"
      }
    }
    begin = 0
    end = {
      type = "indexOf"
      attributes = {
        input = {
          type = "accountAttribute"
          attributes = {
            attributeName = "email"
            sourceName    = "My Source"
          }
        }
        substring = "@"
      }
    }
  })
}

# Example 4: Concatenation Transform - Join multiple values
resource "sailpoint_transform" "concatenation_example" {
  name = "Full Name Builder"
  type = "concatenation"

  attributes = jsonencode({
    values = [
      {
        type = "accountAttribute"
        attributes = {
          attributeName = "firstName"
          sourceName    = "My Source"
        }
      },
      {
        type = "static"
        attributes = {
          value = " "
        }
      },
      {
        type = "accountAttribute"
        attributes = {
          attributeName = "lastName"
          sourceName    = "My Source"
        }
      }
    ]
  })
}

# Example 5: Replace Transform - Replace text patterns
resource "sailpoint_transform" "replace_example" {
  name = "Phone Number Formatter"
  type = "replace"

  attributes = jsonencode({
    input = {
      type = "accountAttribute"
      attributes = {
        attributeName = "phoneNumber"
        sourceName    = "My Source"
      }
    }
    regex       = "^(\\d{3})(\\d{3})(\\d{4})$"
    replacement = "($1) $2-$3"
  })
}

# Example 6: Conditional Transform - If-then-else logic
resource "sailpoint_transform" "conditional_example" {
  name = "Department Code Mapper"
  type = "conditional"

  attributes = jsonencode({
    expression        = "$department == 'Engineering'"
    positiveCondition = "ENG"
    negativeCondition = {
      type = "conditional"
      attributes = {
        expression        = "$department == 'Marketing'"
        positiveCondition = "MKT"
        negativeCondition = "OTHER"
      }
    }
  })
}

# Example 7: Date Format Transform - Format dates
resource "sailpoint_transform" "date_format_example" {
  name = "Birth Date Formatter"
  type = "dateFormat"

  attributes = jsonencode({
    input = {
      type = "accountAttribute"
      attributes = {
        attributeName = "birthDate"
        sourceName    = "My Source"
      }
    }
    inputFormat  = "yyyy-MM-dd"
    outputFormat = "MM/dd/yyyy"
  })
}

# Example 8: Static Transform - Return a fixed value
resource "sailpoint_transform" "static_example" {
  name = "Default Department"
  type = "static"

  attributes = jsonencode({
    value = "General"
  })
}

# Example 9: Using variables for flexible configuration
variable "source_name" {
  description = "Name of the source system"
  type        = string
  default     = "Active Directory"
}

variable "environment" {
  description = "Environment (dev, staging, prod)"
  type        = string
  default     = "dev"
}

resource "sailpoint_transform" "variable_example" {
  name = "Email Builder - ${upper(var.environment)}"
  type = "concatenation"

  attributes = jsonencode({
    values = [
      {
        type = "accountAttribute"
        attributes = {
          attributeName = "firstName"
          sourceName    = var.source_name
        }
      },
      {
        type = "static"
        attributes = {
          value = "."
        }
      },
      {
        type = "accountAttribute"
        attributes = {
          attributeName = "lastName"
          sourceName    = var.source_name
        }
      },
      {
        type = "static"
        attributes = {
          value = "@${var.environment}.company.com"
        }
      }
    ]
  })
}

# Example 10: Multiple transforms with for_each
variable "department_mappings" {
  description = "Map of departments to codes"
  type = map(object({
    code        = string
    description = string
  }))
  default = {
    "engineering" = {
      code        = "ENG"
      description = "Engineering Department"
    }
    "marketing" = {
      code        = "MKT"
      description = "Marketing Department"
    }
    "sales" = {
      code        = "SAL"
      description = "Sales Department"
    }
  }
}

resource "sailpoint_transform" "department_mappers" {
  for_each = var.department_mappings

  name = "Department Mapper - ${title(each.key)}"
  type = "conditional"

  attributes = jsonencode({
    expression        = "$department == '${title(each.key)}'"
    positiveCondition = each.value.code
    negativeCondition = "OTHER"
  })
}

# Example 11: Complex nested transform
resource "sailpoint_transform" "complex_example" {
  name = "Complex User ID Generator"
  type = "concatenation"

  attributes = jsonencode({
    values = [
      # First initial
      {
        type = "substring"
        attributes = {
          input = {
            type = "accountAttribute"
            attributes = {
              attributeName = "firstName"
              sourceName    = var.source_name
            }
          }
          begin = 0
          end   = 1
        }
      },
      # Last name (up to 7 chars, lowercase)
      {
        type = "lower"
        attributes = {
          input = {
            type = "substring"
            attributes = {
              input = {
                type = "accountAttribute"
                attributes = {
                  attributeName = "lastName"
                  sourceName    = var.source_name
                }
              }
              begin = 0
              end   = 7
            }
          }
        }
      },
      # Employee ID (last 3 digits)
      {
        type = "substring"
        attributes = {
          input = {
            type = "accountAttribute"
            attributes = {
              attributeName = "employeeId"
              sourceName    = var.source_name
            }
          }
          begin = -3
        }
      }
    ]
  })
}

# Example 12: Transform outputs for reference
output "transform_ids" {
  description = "IDs of created transforms"
  value = {
    upper         = sailpoint_transform.upper_example.id
    lower         = sailpoint_transform.lower_example.id
    substring     = sailpoint_transform.substring_example.id
    concatenation = sailpoint_transform.concatenation_example.id
    replace       = sailpoint_transform.replace_example.id
    conditional   = sailpoint_transform.conditional_example.id
    date_format   = sailpoint_transform.date_format_example.id
    static        = sailpoint_transform.static_example.id
    variable      = sailpoint_transform.variable_example.id
    complex       = sailpoint_transform.complex_example.id
  }
}

output "transform_info" {
  description = "Detailed information about transforms"
  value = {
    total_count = length([
      sailpoint_transform.upper_example,
      sailpoint_transform.lower_example,
      sailpoint_transform.substring_example,
      sailpoint_transform.concatenation_example,
      sailpoint_transform.replace_example,
      sailpoint_transform.conditional_example,
      sailpoint_transform.date_format_example,
      sailpoint_transform.static_example,
      sailpoint_transform.variable_example,
      sailpoint_transform.complex_example
    ]) + length(sailpoint_transform.department_mappers)

    basic_examples = {
      upper_transform = {
        id   = sailpoint_transform.upper_example.id
        name = sailpoint_transform.upper_example.name
        type = sailpoint_transform.upper_example.type
      }
      lower_transform = {
        id   = sailpoint_transform.lower_example.id
        name = sailpoint_transform.lower_example.name
        type = sailpoint_transform.lower_example.type
      }
    }
  }
}

# Example 13: Import existing transform
# To import: terraform import sailpoint_transform.imported_transform "existing-transform-id"
resource "sailpoint_transform" "imported_transform" {
  name = "Imported Transform"
  type = "upper"

  attributes = jsonencode({
    input = {
      type = "accountAttribute"
      attributes = {
        attributeName = "sAMAccountName"
        sourceName    = "Active Directory"
      }
    }
  })

  lifecycle {
    # Prevent accidental deletion of important transforms
    prevent_destroy = true
  }
}
