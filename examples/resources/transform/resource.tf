resource "sailpoint_transform" "example" {
  name = "example-transform"
  type = "upper"
  attributes = jsonencode({
    "input" : {
      "attributes" : {
        "value" : "about to be deleted"
      },
      "type" : "static"
    }
  })
}
