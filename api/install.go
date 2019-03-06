package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/project-flogo/cli/common"
	"github.com/project-flogo/cli/util"
)

func InstallPackage(project common.AppProject, pkg string) error {

	flogoImport, err := util.ParseImport(pkg)
	if err != nil {
		return err
	}

	err = project.AddImports(false, flogoImport)
	if err != nil {
		return err
	}

	path, err := project.GetPath(flogoImport)
	if Verbose() {
		fmt.Println("Installed path", path)
	}
	if err != nil {
		return err
	}

	desc, err := util.GetContribDescriptor(path)
	if desc != nil {
		fmt.Printf("Installed %s: %s\n", desc.GetContribType(), pkg)
	}

	legacySupportRequired, err := IsLegacySupportRequired(desc, path, pkg, true)
	if err != nil {
		return err
	}

	if legacySupportRequired {
		InstallLegacySupport(project)
	}

	return nil
}
func InstallLocalPackage(project common.AppProject, localPath string, pkg string) error {

	project.DepManager().InstallLocalPkg(pkg, localPath)

	return InstallPackage(project, pkg)
}
func InstallContribBundle(project common.AppProject, path string) error {

	file, err := ioutil.ReadFile(path)

	if err != nil {
		return err
	}

	var contribBundleDescriptor util.FlogoContribBundleDescriptor
	err = json.Unmarshal(file, &contribBundleDescriptor)

	if err != nil {
		return err
	}

	for _, contrib := range contribBundleDescriptor.Contribs {
		InstallPackage(project, contrib)
	}

	return nil

}

func ListPackages(project common.AppProject, format bool, all bool) error {

	err := util.ExecCmd(exec.Command("go", "mod", "tidy"), project.SrcDir())
	if err != nil {
		fmt.Println("Error in tidying up modules")
		return err
	}

	var contribs util.Imports

	if all {
		imports, _ := util.GetAllImports(filepath.Join(project.SrcDir(), fileImportsGo)) // Get Imports from imports.go
		for _, i := range imports {
			flogoImport, _ := util.ParseImport(i)
			contribs = append(contribs, flogoImport)
		}
	} else {
		contribs, _ = util.GetImports(filepath.Join(project.Dir(), fileFlogoJson)) // Get Imports from flogo.json
	}

	var result []interface{}

	for _, contrib := range contribs {
		path, err := project.GetPath(contrib)
		if Verbose() {
			fmt.Println("Path of contrib", path, "for contrib", contrib)
		}

		if err != nil {
			return err
		}

		desc, err := util.GetContribDescriptor(path)

		if Verbose() {
			fmt.Println("Path of contrib descriptor", desc)
		}
		if err != nil {
			return err
		}
		if desc == nil {
			continue
		}
		data := struct {
			Name        string `json:"name"`
			Type        string `json:"type"`
			Description string `json:"description"`
			Homepage    string `json:"homepage"`
			Ref         string `json:"ref"`
			Path        string `json:"path"`
		}{
			desc.Name,
			desc.Type,
			desc.Description,
			desc.Homepage,
			contrib.ModulePath(),
			getDescriptorFile(path),
		}

		result = append(result, data)
	}
	if format {
		resp, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return err
		}

		fmt.Fprintf(os.Stdout, "%v \n", string(resp))
	}

	return nil
}
func getDescriptorFile(path string) string {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return ""
	}

	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".json") {
			return filepath.Join(path, f.Name())
		}
	}
	return ""
}
