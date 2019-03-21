package api

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/project-flogo/cli/common"
	"github.com/project-flogo/core/app"
)

func readAppDescriptor(project common.AppProject) (*app.Config, error) {

	appDescriptorFile, err := os.Open(filepath.Join(project.Dir(), fileFlogoJson))
	if err != nil {
		return nil, err
	}
	defer appDescriptorFile.Close()

	appDescriptorData, err := ioutil.ReadAll(appDescriptorFile)
	if err != nil {
		return nil, err
	}

	var appDescriptor app.Config
	err = json.Unmarshal([]byte(appDescriptorData), &appDescriptor)
	if err != nil {
		return nil, err
	}

	return &appDescriptor, nil
}

func writeAppDescriptor(project common.AppProject, appDescriptor *app.Config)  error {

	appDescriptorUpdated, err := json.MarshalIndent(appDescriptor, "", "  ")
	if err != nil {
		return err
	}

	appDescriptorUpdatedJson := string(appDescriptorUpdated)

	err = ioutil.WriteFile(filepath.Join(project.Dir(), fileFlogoJson), []byte(appDescriptorUpdatedJson), 0644)
	if err != nil {
		return err
	}

	return nil
}