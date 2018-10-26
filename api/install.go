package api

import (
	"fmt"

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
