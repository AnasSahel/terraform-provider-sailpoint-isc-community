terraform {
  required_providers {
    sailpoint = {
      source = "anasSahel/sailpoint-isc-community"
    }
  }
}

provider "sailpoint" {
  # Configuration options
}

# Example: Active Directory Source
resource "sailpoint_source" "active_directory" {
  name        = "Corporate Active Directory"
  description = "Main Active Directory source for employee identities"
  connector   = "active-directory"

  owner = jsonencode({
    type = "IDENTITY"
    id   = "2c91808570313110017040b06f344ec9"
    name = "john.doe"
  })

  configuration = {
    domain_name       = "corp.example.com"
    domain_controller = "dc1.corp.example.com"
    forest_settings = jsonencode([
      {
        domain                = "corp.example.com"
        domain_controller     = "dc1.corp.example.com"
        authentication_domain = "corp.example.com"
        user_search_dn        = "DC=corp,DC=example,DC=com"
        group_search_dn       = "DC=corp,DC=example,DC=com"
        use_ssl               = "true"
      }
    ])
    group_search_filter            = "(objectClass=group)"
    user_search_filter             = "(&(objectClass=user)(objectCategory=person))"
    authorization_type             = "simple"
    account_username               = "svc-sailpoint@corp.example.com"
    account_password               = var.ad_service_password
    search_dn                      = "DC=corp,DC=example,DC=com"
    use_tls                        = "true"
    partition_mode                 = "auto"
    group_membership_search_dn     = "DC=corp,DC=example,DC=com"
    group_membership_search_filter = "(member={0})"
  }
}

# Example: Delimited File Source
resource "sailpoint_source" "employee_csv" {
  name        = "Employee CSV Import"
  description = "CSV file source for bulk employee data import"
  connector   = "delimited-file"

  owner = jsonencode({
    type = "IDENTITY"
    id   = "2c91808570313110017040b06f344ec9"
    name = "john.doe"
  })

  configuration = {
    file         = "employees.csv"
    delimiter    = ","
    has_header   = "true"
    column_names = "username,firstName,lastName,email,department,title,manager"
    merge_columns = jsonencode({
      displayName = "$firstName $lastName"
      fullName    = "$lastName, $firstName"
    })
    identity_attribute  = "username"
    group_column_name   = "department"
    group_delimiter     = "|"
    enable_partitioning = "false"
  }
}

# Variables for sensitive data
variable "ad_service_password" {
  description = "Service account password for Active Directory"
  type        = string
  sensitive   = true
}
