package client

const (
	JSONPatchOpReplace = "replace"
	JSONPatchOpAdd     = "add"
	JSONPatchOpRemove  = "remove"
)

type JSONPatchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

func NewReplaceJSONPatchOperation(path string, value interface{}) JSONPatchOperation {
	return JSONPatchOperation{
		Op:    JSONPatchOpReplace,
		Path:  path,
		Value: value,
	}
}

func NewAddJSONPatchOperation(path string, value interface{}) JSONPatchOperation {
	return JSONPatchOperation{
		Op:    JSONPatchOpAdd,
		Path:  path,
		Value: value,
	}
}

func NewRemoveJSONPatchOperation(path string) JSONPatchOperation {
	return JSONPatchOperation{
		Op:   JSONPatchOpRemove,
		Path: path,
	}
}

type ObjectRef struct {
	Type string `json:"type"`
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}
