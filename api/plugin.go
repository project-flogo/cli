package api

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

func InstallPluginHelper(args string) error {

	//Get the current GOPATH.
	path := os.Getenv("GOPATH")
	//Change the Dir.
	err := os.Chdir(Concat(path, "/src/github.com/project-flogo/cli/cmd/flogo"))

	die(err)

	currdir, _ := os.Getwd()

	AddModToImportPlugin(args, currdir) //Edit the imports.go file

	//Download all the modules. This is just to ensure all packages are downloaded before go build.
	cliCmd, err := exec.Command("go", "mod", "download").CombinedOutput()
	if err != nil {
		RemoveModFromImportPlugin(args, currdir)

		fmt.Println(string(cliCmd))

		log.Fatal(err)
	}
	die(err)

	return err
}

func BuildModule(args string, flag bool) {

	//Get the current GOPATH.
	path := os.Getenv("GOPATH")
	//Change the Dir.
	err := os.Chdir(Concat(path, "/src/github.com/project-flogo/cli/cmd/flogo"))

	die(err)

	currdir, _ := os.Getwd()

	//Build the modules.
	cliCmd, err := exec.Command("go", "build").CombinedOutput()

	if err != nil {
		RemoveModFromImportPlugin(args, currdir)

		fmt.Println(string(cliCmd))

		log.Fatal(err)
	}
	cliCmd, err = exec.Command("cp", Concat(os.Getenv("GOPATH"), "/src/github.com/project-flogo/cli/cmd/flogo/flogo"), Concat(os.Getenv("GOPATH"), "/bin")).CombinedOutput()
	if err != nil {
		RemoveModFromImportPlugin(args, currdir)

		fmt.Println(string(cliCmd))

		log.Fatal(err)
	}
	if flag {
		//Done.
		fmt.Println("Module Successfully Installed")
	}

}

func AddModToImportPlugin(pkg string, fpath string) {

	byteArray, err := ioutil.ReadFile(Concat(fpath, "/imports.go"))

	die(err)

	text := string(byteArray)

	index := strings.Index(text, ")")

	err = ioutil.WriteFile(Concat(fpath, "/imports.go"), []byte(Concat(text[:index-1], "\n _ \"", pkg, "\" \n", ")")), 0)

	die(err)
}

func RemoveModFromImportPlugin(pkg string, fpath string) {

	byteArray, err := ioutil.ReadFile(Concat(fpath, "/imports.go"))

	die(err)

	text := string(byteArray)

	newText := strings.Replace(text, Concat(" _ \"", pkg, "\" "), "", -1)
	err = ioutil.WriteFile(Concat(fpath, "/imports.go"), []byte(newText), 0)

	die(err)
}
