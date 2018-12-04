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
func InstallLocalPackage(project common.AppProject, pkgs []string) error {

	project.DepManager().InstallLocalPkg(pkgs[0], pkgs[1])

	return InstallPackage(project, pkgs[0])
}
func ListPackages(project common.AppProject, format bool) error {

	contribs, _ := util.GetImports(filepath.Join(project.Dir(), fileFlogoJson))
	var result []interface{}
	for _, contrib := range contribs {
		path, err := project.GetPath(contrib)

		if err != nil {
			return err
		}

		desc, err := util.GetContribDescriptor(path)
		data := struct {
			Name        string `json:"name"`
			Type        string `json:"type"`
			Description string `json:"descriptiom"`
			Ref         string `json:"ref"`
			Path        string `json:"path"`
		}{
			desc.Name,
			desc.Type,
			desc.Description,
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
