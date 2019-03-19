package api

import (
	"os/exec"
	"path/filepath"

	"github.com/project-flogo/cli/common"
	"github.com/project-flogo/cli/util"
)

func SyncPkg(project common.AppProject) error {

	pkgs, err := util.GetImportsFromJSON(filepath.Join(project.Dir(), "flogo.json"))

	if err != nil {
		return err
	}

	err = updateGoMod(project, pkgs)

	if err != nil {
		return err
	}

	return nil
}
func updateGoMod(project common.AppProject, pkgs util.Imports) error {

	err := clearGoMod(project.SrcDir())

	if err != nil {
		return err
	}

	for _, pkg := range pkgs {
		err = project.DepManager().AddDependency(pkg)
		if err != nil {
			return err
		}
	}
	return nil
}

func clearGoMod(src string) error {
	err := util.ExecCmd(exec.Command("rm", "go.mod"), src)
	if err != nil {
		return err
	}
	err = util.ExecCmd(exec.Command("rm", "go.sum"), src)
	if err != nil {
		return err
	}
	err = util.ExecCmd(exec.Command("go", "mod", "init", "main"), src)
	if err != nil {
		return err
	}
	return nil

}
