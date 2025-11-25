# Example: Basic launcher for an interactive process
resource "sailpoint_launcher" "group_creation" {
  name        = "Group Creation Launcher"
  description = "Interactive launcher to create Active Directory groups"
  type        = "INTERACTIVE_PROCESS"
  disabled    = false

  reference = {
    type = "WORKFLOW"
    id   = "2c91808a7b5c3e1d017b5c4a8f6d0001"
  }

  config = jsonencode({
    label       = "Create AD Group"
    description = "Create a new Active Directory group with specified members"
    inputs = [
      {
        id       = "groupName"
        label    = "Group Name"
        type     = "text"
        required = true
      },
      {
        id       = "groupDescription"
        label    = "Description"
        type     = "text"
        required = false
      },
      {
        id       = "members"
        label    = "Initial Members"
        type     = "multiselect"
        required = false
      }
    ]
    submitLabel = "Create Group"
    cancelLabel = "Cancel"
  })
}

# Example: Launcher for user onboarding workflow
resource "sailpoint_launcher" "user_onboarding" {
  name        = "User Onboarding"
  description = "Interactive process to onboard new employees"
  type        = "INTERACTIVE_PROCESS"
  disabled    = false

  reference = {
    type = "WORKFLOW"
    id   = "2c91808568c529c60168cca6f90c1314"
  }

  config = jsonencode({
    label       = "Onboard New Employee"
    description = "Start the onboarding process for a new employee"
    inputs = [
      {
        id       = "firstName"
        label    = "First Name"
        type     = "text"
        required = true
      },
      {
        id       = "lastName"
        label    = "Last Name"
        type     = "text"
        required = true
      },
      {
        id       = "email"
        label    = "Email Address"
        type     = "email"
        required = true
      },
      {
        id       = "department"
        label    = "Department"
        type     = "select"
        required = true
        options  = ["Engineering", "Sales", "Marketing", "Finance"]
      },
      {
        id       = "manager"
        label    = "Manager"
        type     = "identity"
        required = true
      },
      {
        id       = "startDate"
        label    = "Start Date"
        type     = "date"
        required = true
      }
    ]
    submitLabel = "Start Onboarding"
    cancelLabel = "Cancel"
  })
}

# Example: Disabled launcher (can be enabled later)
resource "sailpoint_launcher" "access_request" {
  name        = "Request Access"
  description = "Interactive launcher for requesting application access"
  type        = "INTERACTIVE_PROCESS"
  disabled    = true

  reference = {
    type = "WORKFLOW"
    id   = "2c91808568c529c60168cca6f90c1315"
  }

  config = jsonencode({
    label       = "Request Application Access"
    description = "Submit a request for access to an application"
    inputs = [
      {
        id       = "application"
        label    = "Application"
        type     = "select"
        required = true
        options  = ["Salesforce", "Workday", "ServiceNow"]
      },
      {
        id       = "accessLevel"
        label    = "Access Level"
        type     = "select"
        required = true
        options  = ["Read Only", "Standard User", "Power User", "Admin"]
      },
      {
        id       = "justification"
        label    = "Business Justification"
        type     = "textarea"
        required = true
      }
    ]
    submitLabel = "Submit Request"
    cancelLabel = "Cancel"
  })
}
