package client

// JSON Patch operation constants as defined in RFC 6902.
const (
	JSONPatchOpReplace = "replace"
	JSONPatchOpAdd     = "add"
	JSONPatchOpRemove  = "remove"
)

// JSONPatchOperation represents a JSON Patch operation as defined in RFC 6902.
// See: https://datatracker.ietf.org/doc/html/rfc6902
type JSONPatchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

// NewReplaceJSONPatchOperation creates a "replace" JSON Patch operation.
func NewReplaceJSONPatchOperation(path string, value interface{}) JSONPatchOperation {
	return JSONPatchOperation{
		Op:    JSONPatchOpReplace,
		Path:  path,
		Value: value,
	}
}

// NewAddJSONPatchOperation creates an "add" JSON Patch operation.
func NewAddJSONPatchOperation(path string, value interface{}) JSONPatchOperation {
	return JSONPatchOperation{
		Op:    JSONPatchOpAdd,
		Path:  path,
		Value: value,
	}
}

// NewRemoveJSONPatchOperation creates a "remove" JSON Patch operation.
func NewRemoveJSONPatchOperation(path string) JSONPatchOperation {
	return JSONPatchOperation{
		Op:   JSONPatchOpRemove,
		Path: path,
	}
}
