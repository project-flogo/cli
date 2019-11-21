package api

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/project-flogo/cli/common"
	"github.com/project-flogo/cli/util"
)

type AppBuilder struct {

}

func (*AppBuilder) Build(project common.AppProject) error {

	err := restoreMain(project)
	if err != nil {
		return err
	}

	err = simpleGoBuild(project)
	if err != nil {
		return err
	}

	return nil
}


func simpleGoBuild(project common.AppProject) error {
	if _, err := os.Stat(project.BinDir()); err != nil {
		if Verbose() {
			fmt.Println("Creating 'bin' directory")
		}
		err = os.MkdirAll(project.BinDir(), os.ModePerm)
		if err != nil {
			return err
		}
	}

	if Verbose() {
		fmt.Println("Performing 'go build'...")
	}

	err := util.ExecCmd(exec.Command("go", "build", "-o", project.Executable()), project.SrcDir())
	if err != nil {
		fmt.Println("Error in building", project.SrcDir())
		return err
	}

	return nil
}