package api

import (
	"github.com/project-flogo/cli/common"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/project-flogo/cli/util"
)

func BuildProject(project common.AppProject) error {

	err := util.ExecCmd(exec.Command("go", "build"), project.SrcDir())
	if err != nil {
		return err
	}

	if runtime.GOOS == "windows" {
		err = os.Rename(filepath.Join(project.SrcDir(), "main.exe"), project.Executable())
	} else {
		err = os.Rename(filepath.Join(project.SrcDir(), "main"), project.Executable())
	}

	if err != nil {
		return err
	}

	return nil
}
