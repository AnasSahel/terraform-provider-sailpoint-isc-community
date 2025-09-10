
terraform {
  required_providers {
    sailpoint = {
      source = "hashicorp.com/edu/sailpoint-isc-community"
    }
  }
}

provider "sailpoint" {}

resource "sailpoint_transform" "example" {
  name = "example-transform"
  type = "upper"
  attributes_json = jsonencode({
    "input" : {
      "attributes" : {
        "value" : "active"
      },
      "type" : "static"
    }
  })
}
