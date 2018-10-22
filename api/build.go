package api

import (
	"fmt"
	"os"
	"os/exec"
)

func BuildProject() {
	if CheckCurrDir() {
		path, _ := os.Getwd()
		//Move to the src/ dir in the App.
		os.Chdir(Concat(path, "/src"))
		cliCmd, err := exec.Command("go", "build").CombinedOutput()
		if err != nil {
			fmt.Println(string(cliCmd))
		}
		die(err)
		_, err = exec.Command("cp", "main", "../bin/").Output()
		die(err)
		//Reset the Dir.
		os.Chdir(path)
	}
}

func CheckCurrDir() bool {
	currDir, _ := os.Getwd()

	_, err := os.Stat(Concat(currDir, "/src/main.go"))
	_, err1 := os.Stat(Concat(currDir, "/src/imports.go"))
	_, err2 := os.Stat(Concat(currDir, "/src/go.mod"))

	return !(os.IsNotExist(err) && os.IsNotExist(err1) && os.IsNotExist(err2))
}
