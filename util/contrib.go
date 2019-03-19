package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var contribDescriptors = []string{"descriptor.json", "activity.json", "trigger.json", "action.json"}

// FlogoAppDescriptor is the descriptor for a Flogo application
type FlogoContribDescriptor struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Homepage    string `json:"homepage"`
	Shim        string `json:"shim"`
	Ref         string `json:"ref"` //legacy

	IsLegacy bool `json:"-"`
}

type FlogoContribBundleDescriptor struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Contribs    []string `json:"contributions"`
}

func (d *FlogoContribDescriptor) GetContribType() string {
	parts := strings.Split(d.Type, ":")
	if len(parts) > 1 {
		return parts[1]
	}
	return ""
}


func GetContribDescriptorFromImport(depManager DepManager, contribImport Import) (*FlogoContribDescriptor, error) {

	contribPath, err := depManager.GetPath(contribImport)
	if err != nil {
		return nil, err
	}

	return GetContribDescriptor(contribPath)
}

func GetContribDescriptor(contribPath string) (*FlogoContribDescriptor, error) {

	var descriptorPath string

	for _, descriptorName := range contribDescriptors {
		dPath := filepath.Join(contribPath, descriptorName)
		if _, err := os.Stat(dPath); err == nil {
			descriptorPath = dPath
		}
	}

	if descriptorPath == "" {
		//descriptor not found
		return nil, nil
	}

	if _, err := os.Stat(descriptorPath); descriptorPath != "" && err == nil {

		desc, err := ReadContribDescriptor(descriptorPath)
		if err != nil {
			return nil, err
		}

		return desc, nil
	}

	return nil, nil
}

func ReadContribDescriptor(descriptorFile string) (*FlogoContribDescriptor, error) {

	descriptorJson, err := os.Open(descriptorFile)
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(descriptorJson)
	if err != nil {
		return nil, err
	}

	descriptor := &FlogoContribDescriptor{}

	err = json.Unmarshal(bytes, descriptor)
	if err != nil {
		return nil, fmt.Errorf("failed to parse descriptor '%s': %s", descriptorFile, err.Error())
	}

	return descriptor, nil
}
