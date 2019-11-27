package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/project-flogo/cli/common"
	"github.com/project-flogo/cli/util"
)

func InstallPackage(project common.AppProject, pkg string) error {

	flogoImport, err := util.ParseImport(pkg)
	if err != nil {
		return err
	}

	err = project.AddImports(false, true, flogoImport)
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

	legacySupportRequired := false
	desc, err := util.GetContribDescriptor(path)
	if desc != nil {
		cType := desc.GetContribType()
		if desc.IsLegacy {
			legacySupportRequired = true
			cType = "legacy " + desc.GetContribType()
			err := CreateLegacyMetadata(path, desc.GetContribType(), pkg)
			if err != nil {
				return err
			}
		}

		fmt.Printf("Installed %s: %s\n", cType, flogoImport)
		//instStr := fmt.Sprintf("Installed %s:", cType)
		//fmt.Printf("%-20s %s\n", instStr, imp)
	}

	if legacySupportRequired {
		err := InstallLegacySupport(project)
		if err != nil {
			return err
		}
	}

	return nil
}

func InstallReplacedPackage(project common.AppProject, replacedPath string, pkg string) error {

	err := project.DepManager().InstallReplacedPkg(pkg, replacedPath)
	if err != nil {
		return err
	}
	return InstallPackage(project, pkg+"@v0.0.0")
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
		err := InstallPackage(project, contrib)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error installing contrib '%s': %s", contrib, err.Error())
		}
	}

	return nil
}
