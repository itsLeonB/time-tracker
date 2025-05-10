package model

type ExternalTask struct {
	Provider string
	Number   string
	Name     string
	Project  string
}

type YoutrackTask struct {
	IdReadable   string                `json:"idReadable"`
	Summary      string                `json:"summary"`
	CustomFields []YoutrackCustomField `json:"customFields"`
}

type YoutrackCustomField struct {
	Value *YoutrackGenericField `json:"value"`
	Name  string                `json:"name"`
}

type YoutrackGenericField struct {
	Name string `json:"name"`
}

type YoutrackError struct {
	Error                 string          `json:"error"`
	ErrorDescription      string          `json:"error_description"`
	ErrorDeveloperMessage string          `json:"error_developer_message"`
	ErrorField            string          `json:"error_field"`
	ErrorChildren         []YoutrackError `json:"error_children"`
}

func (yt *YoutrackTask) GetEpicName() string {
	for _, field := range yt.CustomFields {
		if field.Name == "Epic" {
			if field.Value == nil {
				return ""
			}

			return field.Value.Name
		}
	}

	return ""
}
