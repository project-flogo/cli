package api

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/project-flogo/cli/common"
)

func UpdatePkg(project common.AppProject, pkg string) error {

	if Verbose() {
		fmt.Println("Updating ", pkg)
	}

	os.Chdir(project.SrcDir())
	output, err := exec.Command("go", "get", "-u", pkg).CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		return err
	}
	return nil
}
