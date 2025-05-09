package model

type ExternalTask struct {
	Provider string
	Number   string
	Name     string
	Project  string
}

type YoutrackTask struct {
	Id           string                `json:"id"`
	IdReadable   string                `json:"idReadable"`
	Summary      string                `json:"summary"`
	CustomFields []YoutrackCustomField `json:"customFields"`
	Type         string                `json:"$type"`
}

type YoutrackCustomField struct {
	Value *YoutrackGenericField `json:"value"`
	Name  string                `json:"name"`
	Type  string                `json:"$type"`
}

type YoutrackGenericField struct {
	Name string `json:"name"`
	Type string `json:"$type"`
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
