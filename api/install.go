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

func InstallPackage(project common.AppProject, pkg string) error {

	err := project.AddImports(false, pkg)
	if err != nil {
		return err
	}

	path, err := project.GetPath(pkg)
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
func InstallPalette(project common.AppProject, path string) error {

	file, err := ioutil.ReadFile(path)

	if err != nil {
		return err
	}

	var paletteDescriptor []*util.FlogoPaletteDescriptor
	err = json.Unmarshal(file, &paletteDescriptor)

	if err != nil {
		return err
	}

	for _, palette := range paletteDescriptor {
		InstallPackage(project, fmt.Sprintf("%v", palette.Reference))
	}

	return nil

}
func ListPackages(project common.AppProject, format bool, all bool) error {
	var contribs []string

	if all {
		contribs, _ = util.GetAllImports(filepath.Join(project.SrcDir(), fileImportsGo)) // Get Imports from imports.go

	} else {
		contribs, _ = util.GetImports(filepath.Join(project.Dir(), fileFlogoJson)) // Get Imports from flogo.json

	}

	var result []interface{}

	for _, contrib := range contribs {
		path, err := project.GetPath(contrib)

		if err != nil {
			return err
		}

		desc, err := util.GetContribDescriptor(path)
		if err != nil || desc == nil {
			return err
		}
		data := struct {
			Name        string `json:"name"`
			Type        string `json:"type"`
			Description string `json:"descriptiom"`
			Homepage    string `json:"homepage"`
			Ref         string `json:"ref"`
			Path        string `json:"path"`
		}{
			desc.Name,
			desc.Type,
			desc.Description,
			desc.Homepage,
			desc.Ref,
			getDescriptorFile(path),
		}

		result = append(result, data)
	}
	if format {
		resp, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return err
		}

		fmt.Fprintf(os.Stderr, "%v \n", string(resp))
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
