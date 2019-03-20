package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/project-flogo/cli/common"
	"github.com/project-flogo/cli/util"
	"github.com/project-flogo/core/app"
)

func SyncPkg(project common.AppProject) error {

	pkgs, err := util.GetAppImports(filepath.Join(project.Dir(), fileFlogoJson), project.DepManager(), true)

	if err != nil {
		return err
	}
	for _, pkg := range pkgs.GetAllImports() {
		project.RemoveImports(pkg.GoImportPath())
	}

	err = updateGoMod(project, pkgs.GetAllImports())

	if err != nil {
		return err
	}

	return nil
}
func updateGoMod(project common.AppProject, pkgs util.Imports) error {

	var err error
	for _, pkg := range pkgs {
		if Verbose() {
			fmt.Println("Adding dependency for ", pkg)
		}
		err = project.DepManager().AddDependency(pkg)
		if err != nil {
			return err
		}
	}
	return nil
}

func ResolvePkg(project common.AppProject) error {
	/*err := SyncPkg(project)
	if err != nil {
		return err
	}*/
	imports, err := project.DepManager().GetAllImports()
	if err != nil {
		return err
	}

	addImportToJSON(project, imports)

	return nil
}

func addImportToJSON(project common.AppProject, imports map[string]util.Import) error {
	appDescriptorFile := filepath.Join(project.Dir(), fileFlogoJson)
	appDescriptorJsonFile, err := os.Open(appDescriptorFile)
	if err != nil {
		return err
	}
	defer appDescriptorJsonFile.Close()

	appDescriptorData, err := ioutil.ReadAll(appDescriptorJsonFile)
	if err != nil {
		return err
	}

	var appDescriptor app.Config
	json.Unmarshal([]byte(appDescriptorData), &appDescriptor)

	parsedImports, err := util.ParseImports(appDescriptor.Imports)

	if err != nil {
		return err
	}

	var result []string

	for _, Import := range parsedImports {

		if val, ok := imports[Import.ModulePath()]; ok {
			result = append(result, val.CanonicalImport())
		}

	}

	appDescriptor.Imports = result

	appDescriptorUpdated, err := json.MarshalIndent(appDescriptor, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(appDescriptorFile, []byte(appDescriptorUpdated), 0644)
	if err != nil {
		return err
	}

	return nil
}
