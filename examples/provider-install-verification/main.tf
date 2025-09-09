
terraform {
  required_providers {
    sailpoint = {
      source = "hashicorp.com/edu/sailpoint-isc-community"
    }
  }
}

provider "sailpoint" {}

data "sailpoint_transforms" "example" {}
