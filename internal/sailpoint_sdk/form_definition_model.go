package sailpoint_sdk

type FormDefinition struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Created     string `json:"created"`
	Modified    string `json:"modified"`

	Owner FormOwner `json:"owner"`
}

type FormOwner struct {
	Type string `json:"type"`
	Id   string `json:"id"`
	Name string `json:"name"`
}
