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

  owner = {
    type = "IDENTITY"
    id   = "2c91808570313110017040b06f344ec9"
    name = "john.doe"
  }

  # Optional core attributes
  authoritative    = true
  delete_threshold = 10
  features         = ["PROVISIONING", "NO_PERMISSIONS_PROVISIONING", "GROUPS_HAVE_MEMBERS"]

  # Connector-specific configuration as JSON
  connector_attributes = jsonencode({
    domain_name       = "corp.example.com"
    domain_controller = "dc1.corp.example.com"
    forest_settings = [
      {
        domain                = "corp.example.com"
        domain_controller     = "dc1.corp.example.com"
        authentication_domain = "corp.example.com"
        user_search_dn        = "DC=corp,DC=example,DC=com"
        group_search_dn       = "DC=corp,DC=example,DC=com"
        use_ssl               = "true"
      }
    ]
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
  })

  # Optional cluster reference
  cluster = {
    type = "CLUSTER"
    id   = "2c91808570313110017040b06f344ec0"
    name = "Main Cluster"
  }
}

# Example: Delimited File Source
resource "sailpoint_source" "employee_csv" {
  name        = "Employee CSV Import"
  description = "CSV file source for bulk employee data import"
  connector   = "delimited-file"

  owner = {
    type = "IDENTITY"
    id   = "2c91808570313110017040b06f344ec9"
    name = "john.doe"
  }

  # Optional attributes
  authoritative    = false
  delete_threshold = 5
  features         = ["NO_PERMISSIONS_PROVISIONING"]

  # Connector-specific configuration as JSON
  connector_attributes = jsonencode({
    file         = "employees.csv"
    delimiter    = ","
    has_header   = "true"
    column_names = "username,firstName,lastName,email,department,title,manager"
    merge_columns = {
      displayName = "$firstName $lastName"
      fullName    = "$lastName, $firstName"
    }
    identity_attribute  = "username"
    group_column_name   = "department"
    group_delimiter     = "|"
    enable_partitioning = "false"
  })
}

# Example with management workgroup
resource "sailpoint_source" "service_now" {
  name        = "ServiceNow HR System"
  description = "ServiceNow source for HR data integration"
  connector   = "web-services"

  owner = {
    type = "IDENTITY"
    id   = "2c91808570313110017040b06f344ec9"
    name = "john.doe"
  }

  # Management workgroup reference
  management_workgroup = jsonencode({
    type = "WORKGROUP"
    id   = "2c91808570313110017040b06f344ec1"
    name = "IT Operations"
  })

  # Correlation configuration
  account_correlation_config = jsonencode({
    type = "ACCOUNT_CORRELATION_CONFIG"
    id   = "2c91808570313110017040b06f344ec2"
    name = "Employee ID Correlation"
  })

  # Connector configuration
  connector_attributes = jsonencode({
    base_url     = "https://instance.service-now.com"
    username     = var.servicenow_username
    password     = var.servicenow_password
    timeout      = "60"
    batch_size   = "1000"
    table_name   = "sys_user"
    query_filter = "active=true"
  })
}

# Variables for sensitive data
variable "ad_service_password" {
  description = "Service account password for Active Directory"
  type        = string
  sensitive   = true
}

variable "servicenow_username" {
  description = "ServiceNow service account username"
  type        = string
  sensitive   = true
}

variable "servicenow_password" {
  description = "ServiceNow service account password"
  type        = string
  sensitive   = true
}
