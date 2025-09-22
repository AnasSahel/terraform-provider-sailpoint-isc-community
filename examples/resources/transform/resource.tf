# SailPoint Transform Resource Examples
# 
# IMPORTANT NOTES (Enhanced in v0.2.0):
# - The 'name' and 'type' fields are IMMUTABLE after creation (RequiresReplace)
# - The 'type' field is validated against 31 supported transform types
# - The 'attributes' field must contain valid JSON
# - Use 'terraform plan' to see if changes will force resource recreation
#
# For a complete list of supported transform types, see:
# https://documentation.sailpoint.com/saas/help/transforms/

# Example 1: Basic Upper Transform - Convert input to uppercase
resource "sailpoint_transform" "upper_example" {
  name = "Upper Case Transform" # IMMUTABLE - changing this will recreate the resource
  type = "upper"                # IMMUTABLE - validated against supported types

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

# ============================================================================
# VALIDATION EXAMPLES (Enhanced in v0.2.0)
# ============================================================================

# Example 14: Demonstrates validation features
resource "sailpoint_transform" "validation_example" {
  # NOTE: Both 'name' and 'type' are IMMUTABLE (RequiresReplace)
  # Changing either will force Terraform to destroy and recreate the resource
  name = "Validation Demo Transform"

  # The 'type' field is validated against supported transform types:
  # accountAttribute, base64Decode, base64Encode, concatenation, conditional,
  # dateCompare, dateFormat, dateMath, decompose, displayName, e164phone,
  # firstValid, getReference, getReferenceIdentityAttribute, identityAttribute,
  # indexOf, iso3166, lastIndexOf, leftPad, lookup, lower, normalizeNames,
  # randomAlphaNumeric, randomNumeric, replace, replaceAll, rightPad, rule,
  # split, static, substring, trim, upper, uuid
  type = "concatenation"

  # The 'attributes' field must contain valid JSON
  # Invalid JSON will cause a validation error before API calls
  attributes = jsonencode({
    values = [
      {
        type = "accountAttribute"
        attributes = {
          attributeName = "firstName"
          sourceName    = "HR System"
        }
      },
      " ", # Static space
      {
        type = "accountAttribute"
        attributes = {
          attributeName = "lastName"
          sourceName    = "HR System"
        }
      }
    ]
  })
}

# Example 15: All supported transform types (for reference)
locals {
  supported_transform_types = [
    "accountAttribute",
    "base64Decode",
    "base64Encode",
    "concatenation",
    "conditional",
    "dateCompare",
    "dateFormat",
    "dateMath",
    "decompose",
    "displayName",
    "e164phone",
    "firstValid",
    "getReference",
    "getReferenceIdentityAttribute",
    "identityAttribute",
    "indexOf",
    "iso3166",
    "lastIndexOf",
    "leftPad",
    "lookup",
    "lower",
    "normalizeNames",
    "randomAlphaNumeric",
    "randomNumeric",
    "replace",
    "replaceAll",
    "rightPad",
    "rule",
    "split",
    "static",
    "substring",
    "trim",
    "upper",
    "uuid"
  ]
}

# VALIDATION ERRORS YOU MIGHT ENCOUNTER:
# 
# 1. Invalid transform type:
# Error: Invalid Attribute Value Match
# │ "invalidType" is not a valid transform type
# 
# 2. Invalid JSON in attributes:
# Error: Invalid Attribute Value Match  
# │ must be valid JSON object
#
# 3. Attempting to change immutable fields:
# Plan: 1 to add, 0 to change, 1 to destroy.
# │ # sailpoint_transform.example must be replaced
# │ # (because name/type cannot be updated in-place)}
