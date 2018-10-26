package api

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/project-flogo/cli/common"
	"github.com/project-flogo/cli/util"
)

func BuildProject(project common.AppProject) error {

	err := util.ExecCmd(exec.Command("go", "build"), project.SrcDir())
	if err != nil {
		return err
	}

	exe := filepath.Join(project.SrcDir(), "main")

	if runtime.GOOS == "windows" {
		exe = filepath.Join(project.SrcDir(), "main.exe")
	}

	if _, err := os.Stat(exe); err == nil {
		err = os.Rename(exe, project.Executable())
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("failed to build application, run with --verbose to see details")
	}

	return nil
}
