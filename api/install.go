package api

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/project-flogo/cli/common"
	"github.com/project-flogo/cli/util"
)

func InstallPackage(project common.AppProject, pkgs []string, local bool) error {

	if local {
		project.DepManager().InstallLocalPkg(pkgs[0], pkgs[1])
	}
	pkg := pkgs[0]

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

func ListPackages(project common.AppProject, format bool) error {

	contribs, _ := util.GetImports(filepath.Join(project.Dir(), fileFlogoJson))
	var result []util.FlogoContribDescriptor
	for _, contrib := range contribs {
		path, err := project.GetPath(contrib)

		if err != nil {
			return err
		}

		desc, err := util.GetContribDescriptor(path)

		result = append(result, *desc)
	}
	if format {
		resp, err := json.MarshalIndent(result, "", "   ")
		if err != nil {
			return err
		}

		fmt.Fprintf(os.Stderr, "%v \n", string(resp))
	}

	return nil
}
