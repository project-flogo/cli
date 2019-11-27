package api

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/project-flogo/cli/common"
	"github.com/project-flogo/cli/util"
	"github.com/project-flogo/core/app"
)

func ListProjectImports(project common.AppProject) error {

	appImports, err := util.GetAppImports(filepath.Join(project.Dir(), fileFlogoJson), project.DepManager(), false)
	if err != nil {
		return err
	}

	for _, imp := range appImports.GetAllImports() {

		fmt.Fprintf(os.Stdout, "  %s\n", imp)
	}

	return nil
}

func SyncProjectImports(project common.AppProject) error {

	appImports, err := util.GetAppImports(filepath.Join(project.Dir(), fileFlogoJson), project.DepManager(), false)
	if err != nil {
		return err
	}
	appImportsMap := make(map[string]util.Import)
	for _, imp := range appImports.GetAllImports() {
		appImportsMap[imp.GoImportPath()] = imp
	}

	goImports, err := project.GetGoImports(false)
	if err != nil {
		return err
	}
	goImportsMap := make(map[string]util.Import)
	for _, imp := range goImports {
		goImportsMap[imp.GoImportPath()] = imp
	}

	var toAdd []util.Import
	for goPath, imp := range appImportsMap {
		if _, ok := goImportsMap[goPath]; !ok {
			toAdd = append(toAdd, imp)
			if Verbose() {
				fmt.Println("Adding missing Go import: ", goPath)
			}
		}
	}

	if util.FileExists(filepath.Join(project.Dir(), fileEngineJson)) {
		engineImports, err := util.GetEngineImports(filepath.Join(project.Dir(), fileEngineJson), project.DepManager())
		if err != nil {
			return err
		}

		engImportsMap := make(map[string]util.Import)
		for _, imp := range engineImports.GetAllImports() {
			engImportsMap[imp.GoImportPath()] = imp
		}

		for goPath, imp := range engImportsMap {
			if _, ok := goImportsMap[goPath]; !ok {
				toAdd = append(toAdd, imp)
				if Verbose() {
					fmt.Println("Adding missing Go import: ", goPath)
				}
			}
		}
	}

	var toRemove []string
	for goPath := range goImportsMap {
		if _, ok := appImportsMap[goPath]; !ok {
			toRemove = append(toRemove, goPath)
			if Verbose() {
				fmt.Println("Removing extraneous Go import: ", goPath)
			}
		}
	}

	err = project.RemoveImports(toRemove...)
	if err != nil {
		return err
	}

	err = project.AddImports(false, false, toAdd...)
	if err != nil {
		return err
	}

	return nil
}

func ResolveProjectImports(project common.AppProject) error {
	if Verbose() {
		fmt.Fprintln(os.Stdout, "Synchronizing project imports")
	}
	err := SyncProjectImports(project)
	if err != nil {
		return err
	}

	if Verbose() {
		fmt.Fprintln(os.Stdout, "Reading flogo.json")
	}
	appDescriptor, err := readAppDescriptor(project)
	if err != nil {
		return err
	}

	if Verbose() {
		fmt.Fprintln(os.Stdout, "Updating flogo.json import versions")
	}
	err = updateDescriptorImportVersions(project, appDescriptor)
	if err != nil {
		return err
	}

	if Verbose() {
		fmt.Fprintln(os.Stdout, "Saving updated flogo.json")
	}
	err = writeAppDescriptor(project, appDescriptor)
	if err != nil {
		return err
	}

	return nil
}

func updateDescriptorImportVersions(project common.AppProject, appDescriptor *app.Config) error {

	goModImports, err := project.DepManager().GetAllImports()
	if err != nil {
		return err
	}

	appImports, err := util.ParseImports(appDescriptor.Imports)
	if err != nil {
		return err
	}

	var result []string

	for _, appImport := range appImports {

		if goModImport, ok := goModImports[appImport.ModulePath()]; ok {
			updatedImp := util.NewFlogoImportWithVersion(appImport, goModImport.Version())
			result = append(result, updatedImp.CanonicalImport())
		} else {
			//not found, look for import of parent package
			for pkg, goModImport := range goModImports {
				if strings.Contains(appImport.ModulePath(), pkg) {
					updatedImp := util.NewFlogoImportWithVersion(appImport, goModImport.Version())
					result = append(result, updatedImp.CanonicalImport())
				}
			}
		}
	}

	appDescriptor.Imports = result

	return nil
}


