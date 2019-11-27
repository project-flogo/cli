package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/project-flogo/cli/common"
	"github.com/project-flogo/cli/util"
)

var fileSampleEngineMain = filepath.Join("examples", "engine", "main.go")

func CreateProject(basePath, appName, appCfgPath, coreVersion string) (common.AppProject, error) {

	var err error
	var appJson string

	if appCfgPath != "" {

		if util.IsRemote(appCfgPath) {

			appJson, err = util.LoadRemoteFile(appCfgPath)
			if err != nil {
				return nil, fmt.Errorf("unable to load remote app file '%s' - %s", appCfgPath, err.Error())
			}
		} else {
			appJson, err = util.LoadLocalFile(appCfgPath)
			if err != nil {
				return nil, fmt.Errorf("unable to load app file '%s' - %s", appCfgPath, err.Error())
			}
		}
	} else {
		if len(appName) == 0 {
			return nil, fmt.Errorf("app name not specified")
		}
	}

	appName, err = getAppName(appName, appJson)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Creating Flogo App: %s\n", appName)

	appDir, err := createAppDirectory(basePath, appName)
	if err != nil {
		return nil, err
	}

	srcDir := filepath.Join(appDir, "src")
	dm := util.NewDepManager(srcDir)

	if Verbose() {
		fmt.Printf("Setting up app directory: %s\n", appDir)
	}

	err = setupAppDirectory(dm, appDir, coreVersion)
	if err != nil {
		return nil, err
	}

	if Verbose() {
		if appJson == "" {
			fmt.Println("Adding sample flogo.json")
		}
	}
	err = createAppJson(dm, appDir, appName, appJson)
	if err != nil {
		return nil, err
	}

	err = createMain(dm, appDir)
	if err != nil {
		return nil, err
	}

	project := NewAppProject(appDir)

	if Verbose() {
		fmt.Println("Importing Dependencies...")
	}

	err = importDependencies(project)
	if err != nil {
		return nil, err
	}

	if Verbose() {
		fmt.Printf("Created App: %s\n", appName)
	}

	return project, nil
}

// createAppDirectory creates the flogo app directory
func createAppDirectory(basePath, appName string) (string, error) {

	var err error

	if basePath == "." {
		basePath, err = os.Getwd()
		if err != nil {
			return "", err
		}
	}

	appPath := filepath.Join(basePath, appName)
	err = os.Mkdir(appPath, os.ModePerm)
	if err != nil {
		return "", err
	}

	return appPath, nil
}

//setupAppDirectory sets up the flogo app directory
func setupAppDirectory(dm util.DepManager, appPath, coreVersion string) error {

	err := os.Mkdir(filepath.Join(appPath, dirBin), os.ModePerm)
	if err != nil {
		return err
	}

	srcDir := filepath.Join(appPath, dirSrc)
	err = os.Mkdir(srcDir, os.ModePerm)
	if err != nil {
		return err
	}

	_, err = os.Create(filepath.Join(srcDir, fileImportsGo))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filepath.Join(srcDir, fileImportsGo), []byte("package main\n"), 0644)
	if err != nil {
		return err
	}

	err = dm.Init()
	if err != nil {
		return err
	}

	flogoCoreImport := util.NewFlogoImport(flogoCoreRepo, "", coreVersion, "")

	//todo get the actual version installed from the go.mod file
	if coreVersion == "" {
		fmt.Printf("Installing: %s@latest\n", flogoCoreImport.CanonicalImport())
	} else {
		fmt.Printf("Installing: %s\n", flogoCoreImport.CanonicalImport())
	}

	// add & fetch the core library
	err = dm.AddDependency(flogoCoreImport)
	if err != nil {
		return err
	}

	return nil
}

// createAppJson create the flogo app json
func createAppJson(dm util.DepManager, appDir, appName, appJson string) error {

	updatedJson, err := getAndUpdateAppJson(dm, appName, appJson)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath.Join(appDir, fileFlogoJson), []byte(updatedJson), 0644)
	if err != nil {
		return err
	}

	return nil
}

// importDependencies import all dependencies
func importDependencies(project common.AppProject) error {

	ai, err := util.GetAppImports(filepath.Join(project.Dir(), fileFlogoJson), project.DepManager(), true)
	if err != nil {
		return err
	}

	imports := ai.GetAllImports()

	if len(imports) == 0 {
		return nil
	}

	err = project.AddImports(true, false, imports...)
	if err != nil {
		return err
	}

	legacySupportRequired := false

	for _, details := range ai.GetAllImportDetails() {

		path, err := project.GetPath(details.Imp)
		if err != nil {
			return err
		}

		desc, err := util.GetContribDescriptor(path)

		if err != nil {
			return err
		}

		if desc != nil {

			cType := desc.GetContribType()
			if desc.IsLegacy {
				legacySupportRequired = true
				cType = "legacy " + desc.GetContribType()
				err := CreateLegacyMetadata(path, desc.GetContribType(), details.Imp.GoImportPath())
				if err != nil {
					return err
				}
			}

			fmt.Printf("Installed %s: %s\n", cType, details.Imp)
			//instStr := fmt.Sprintf("Installed %s:", cType)
			//fmt.Printf("%-20s %s\n", instStr, imp)
		}
	}

	if legacySupportRequired {
		err := InstallLegacySupport(project)
		return err
	}

	return nil
}

func createMain(dm util.DepManager, appDir string) error {

	flogoCoreImport, err := util.NewFlogoImportFromPath(flogoCoreRepo)
	if err != nil {
		return err
	}

	corePath, err := dm.GetPath(flogoCoreImport)
	if err != nil {
		return err
	}

	bytes, err := ioutil.ReadFile(filepath.Join(corePath, fileSampleEngineMain))
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath.Join(appDir, dirSrc, fileMainGo), bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func getAndUpdateAppJson(dm util.DepManager, appName, appJson string) (string, error) {

	if len(appJson) == 0 {
		appJson = emptyFlogoJson
	}

	descriptor, err := util.ParseAppDescriptor(appJson)
	if err != nil {
		return "", err
	}

	if appName != "" {
		// override the application name

		altJson := strings.Replace(appJson, `"`+descriptor.Name+`"`, `"`+appName+`"`, 1)
		altDescriptor, err := util.ParseAppDescriptor(altJson)

		//see if we can get away with simple replace so we don't reorder the existing json
		if err == nil && altDescriptor.Name == appName {
			appJson = altJson
		} else {
			//simple replace didn't work so we have to unmarshal & re-marshal the supplied json
			var appObj map[string]interface{}
			err := json.Unmarshal([]byte(appJson), &appObj)
			if err != nil {
				return "", err
			}

			appObj["name"] = appName

			updApp, err := json.MarshalIndent(appObj, "", "  ")
			if err != nil {
				return "", err
			}
			appJson = string(updApp)
		}

		descriptor.Name = appName
	} else {
		appName = descriptor.Name
	}

	return appJson, nil
}

func getAppName(appName, appJson string) (string, error) {

	if appJson != "" && appName == "" {
		descriptor, err := util.ParseAppDescriptor(appJson)
		if err != nil {
			return "", err
		}

		return descriptor.Name, nil
	}

	return appName, nil
}
func GetTempDir() (string, error) {

	tempDir, err := ioutil.TempDir("", "flogo")
	if err != nil {
		return "", err
	}
	return tempDir, nil
}

var emptyFlogoJson = `
{
	"name": "{{.AppName}}",
	"type": "flogo:app",
	"version": "0.0.1",
	"description": "My Flogo Application Description",
	"appModel": "1.1.0",
	"imports": [],
	"triggers": [],
	"resources":[]
  }
  `
