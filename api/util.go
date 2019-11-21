package api

import (
	"encoding/json"
	"fmt"
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

func backupMain(project common.AppProject) error {
	mainGo := filepath.Join(project.SrcDir(), fileMainGo)
	mainGoBak := filepath.Join(project.SrcDir(), fileMainGo+".bak")

	if _, err := os.Stat(mainGo); err == nil {
		//main found, check for backup main in case we have to remove it
		if _, err := os.Stat(mainGoBak); err == nil {

			//remove old main backup
			if Verbose() {
				fmt.Printf("Removing old main backup: %s\n", mainGoBak)
			}
			err = os.Rename(mainGoBak, mainGo)
			if err != nil {
				return err
			}
		}
		if Verbose() {
			fmt.Println("Backing up main.go")
		}
		err = os.Rename(mainGo, mainGoBak)
		if err != nil {
			return err
		}
	}

	return nil
}

func restoreMain(project common.AppProject) error {

	mainGo := filepath.Join(project.SrcDir(), fileMainGo)
	mainGoBak := filepath.Join(project.SrcDir(), fileMainGo+".bak")

	if _, err := os.Stat(mainGo); err != nil {
		//main not found, check for backup main
		if _, err := os.Stat(mainGoBak); err == nil {
			if Verbose() {
				fmt.Printf("Restoring main from: %s\n", mainGoBak)
			}
			err = os.Rename(mainGoBak, mainGo)
			if err != nil {
				return err
			}
		} else if _, err := os.Stat(mainGo); err != nil {
			return fmt.Errorf("project corrupt, main missing")
		}
	}

	return nil
}
