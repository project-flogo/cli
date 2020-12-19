package util

import (
	"encoding/json"
)

// ParseAppDescriptor parse the application descriptor
func ParseAppDescriptor(appJson string) (*FlogoAppDescriptor, error) {
	descriptor := &FlogoAppDescriptor{}

	err := json.Unmarshal([]byte(appJson), descriptor)

	if err != nil {
		return nil, err
	}

	return descriptor, nil
}

// FlogoAppDescriptor is the descriptor for a Flogo application
type FlogoAppDescriptor struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Version     string   `json:"version"`
	Description string   `json:"description"`
	AppModel    string   `json:"appModel,omitempty"`
	Imports     []string `json:"imports"`

	Triggers []*FlogoTriggerConfig `json:"triggers"`
	Resources []*FlogoResourceConfig `json:"resources"`
}

type FlogoTriggerConfig struct {
	Id   string `json:"id"`
	Ref  string `json:"ref"`
	Type string `json:"type"`
}

type FlogoResourceConfig struct {
	Id   string `json:"id"`
	Data FlogoResourceDataConfig `json:"data"`
}

type FlogoResourceDataConfig struct {
	Tasks []FlogoResourceDataTaskConfig `json:"tasks"`
}

type FlogoResourceDataTaskConfig struct {
	Id   string `json:"id"`
	Activity FlogoResourceDataTaskActivityConfig `json:"activity"`
}

type FlogoResourceDataTaskActivityConfig struct {
	Ref  string `json:"ref"`
}