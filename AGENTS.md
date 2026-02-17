## Do
- use resty v3. Make sure your code is v3 compatible
- use `jsontypes.Normalized` for JSON strings
- always ask me to compile or lint the code and wait for my answer
- Always redact sensitive and personnal information: id, uid, name, full name, first name, last name, email etc

## Don't
- do not harcode id
- do not try to compile or lint code

## Safety and permissions
Allowed without prompt:
- read files, list files
- find a string in files within the project folder

## Project structure
- Generate REST client in `internal/client/<endpoint_name>.go`, where `endpoint_name` is the endpoint name. (For example: Form Definition becomes form_definition)
- Generate Terraform models related to an endpoint in `internal/services/<endpoint_name>/<endpoint_name>_model.go`
- Generate Terraform Resource related to an endpoint in `internal/services/<endpoint_name>/<endpoint_name>_resource.go`
- Generate Terraform Data Source related to an endpoint in `internal/services/<endpoint_name>/<endpoint_name>_data_source.go`

## Resources
You can find Terraform resources in:
- https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-provider

You can find SailPoint API documentation in:
- https://developer.sailpoint.com/docs/api/v2025