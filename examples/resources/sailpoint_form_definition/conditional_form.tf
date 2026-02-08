# Form with conditional logic to show/hide fields
resource "sailpoint_form_definition" "conditional" {
  name        = "Employee Onboarding Form"
  description = "Form with conditional fields based on employee type"

  owner = {
    type = "IDENTITY"
    id   = "2c91808a7813090a017814121e121518"
  }

  form_elements = jsonencode([
    # Main section
    {
      id          = "section-main"
      elementType = "SECTION"
      key         = "main-info"
      config = {
        label = "Employee Information"
      }
    },
    {
      id          = "employee-type"
      elementType = "SELECT"
      key         = "employeeType"
      config = {
        label   = "Employee Type"
        options = ["Full-Time", "Contractor", "Intern"]
      }
      validations = [
        { validationType = "REQUIRED" }
      ]
    },
    {
      id          = "department"
      elementType = "SELECT"
      key         = "department"
      config = {
        label   = "Department"
        options = ["Engineering", "Sales", "Marketing", "HR", "Finance"]
      }
      validations = [
        { validationType = "REQUIRED" }
      ]
    },
    # Contractor-specific section
    {
      id          = "section-contractor"
      elementType = "SECTION"
      key         = "contractor-info"
      config = {
        label = "Contractor Details"
      }
    },
    {
      id          = "contract-end-date"
      elementType = "DATE"
      key         = "contractEndDate"
      config = {
        label = "Contract End Date"
      }
      validations = [
        { validationType = "REQUIRED" }
      ]
    },
    {
      id          = "vendor-company"
      elementType = "TEXT"
      key         = "vendorCompany"
      config = {
        label       = "Vendor Company"
        placeholder = "Name of contracting company"
      }
      validations = [
        { validationType = "REQUIRED" }
      ]
    },
    # Manager approval section (for Admin access)
    {
      id          = "section-approval"
      elementType = "SECTION"
      key         = "approval-info"
      config = {
        label = "Manager Approval"
      }
    },
    {
      id          = "manager-email"
      elementType = "EMAIL"
      key         = "managerEmail"
      config = {
        label       = "Manager Email"
        placeholder = "manager@company.com"
      }
      validations = [
        { validationType = "REQUIRED" },
        { validationType = "EMAIL" }
      ]
    }
  ])

  form_conditions = jsonencode([
    # Hide contractor section for non-contractors
    {
      ruleOperator = "AND"
      rules = [
        {
          sourceType = "ELEMENT"
          source     = "employeeType"
          operator   = "NE"
          valueType  = "STRING"
          value      = "Contractor"
        }
      ]
      effects = [
        {
          effectType = "HIDE"
          config = {
            element = "section-contractor"
          }
        },
        {
          effectType = "HIDE"
          config = {
            element = "contract-end-date"
          }
        },
        {
          effectType = "HIDE"
          config = {
            element = "vendor-company"
          }
        }
      ]
    },
    # Show manager approval section only for Engineering department
    {
      ruleOperator = "AND"
      rules = [
        {
          sourceType = "ELEMENT"
          source     = "department"
          operator   = "NE"
          valueType  = "STRING"
          value      = "Engineering"
        }
      ]
      effects = [
        {
          effectType = "HIDE"
          config = {
            element = "section-approval"
          }
        },
        {
          effectType = "HIDE"
          config = {
            element = "manager-email"
          }
        }
      ]
    }
  ])
}
