package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	fileDescriptorJson string = "descriptor.json"
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
	oldDescriptor := false

	for _, descriptorName := range contribDescriptors {
		dPath := filepath.Join(contribPath, descriptorName)
		if _, err := os.Stat(dPath); err == nil {
			if descriptorName != "descriptor.json" {
				oldDescriptor = true
			}
			descriptorPath = dPath
		}
	}

	if descriptorPath == "" {
		//descriptor not found
		return nil, nil
	}

	desc, err := ReadContribDescriptor(descriptorPath)
	if err != nil {
		return nil, err
	}

	if desc.Ref != "" && oldDescriptor {
		desc.IsLegacy = true
	}

	return desc, nil
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

func GetContribType(depManager DepManager, ref string) (string, error) {

	refAsFlogoImport, err := NewFlogoImportFromPath(ref)
	if err != nil {
		return "", err
	}

	impPath, err := depManager.GetPath(refAsFlogoImport)//(refAsFlogoImport)
	if err != nil {
		return "", err
	}
	var descriptorPath string

	if _, err := os.Stat(filepath.Join(impPath, fileDescriptorJson)); err == nil {
		descriptorPath = filepath.Join(impPath, fileDescriptorJson)

	} else if _, err := os.Stat(filepath.Join(impPath, "activity.json")); err == nil {
		descriptorPath = filepath.Join(impPath, "activity.json")
	} else if _, err := os.Stat(filepath.Join(impPath, "trigger.json")); err == nil {
		descriptorPath = filepath.Join(impPath, "trigger.json")
	} else if _, err := os.Stat(filepath.Join(impPath, "action.json")); err == nil {
		descriptorPath = filepath.Join(impPath, "action.json")
	}

	if _, err := os.Stat(descriptorPath); descriptorPath != "" && err == nil {

		desc, err := ReadContribDescriptor(descriptorPath)
		if err != nil {
			return "", err
		}

		return desc.Type, nil
	}

	return "", nil
}
