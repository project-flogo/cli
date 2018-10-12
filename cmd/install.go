package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the version module",
	Long:  `All software has versions. This is Hugo's`,
	Run: func(cmd *cobra.Command, args []string) {
		if checkCurrDir() {
			currDir, _ := os.Getwd()
			InstallPackage(os.Args[2], currDir)
		} else {
			fmt.Println("Error in detecting app")
		}

	},
}

func init() {
	RootCmd.AddCommand(installCmd)
}
func InstallPackage(args string, path string) {
	os.Chdir(path)
	legacySupport = true
	if InstallPackageHelper(args) == nil {

		BuildAppModule(args)
		path, _ := os.Getwd()

		if NeedsLegacySupport(args) && legacySupport {
			fmt.Println("Needs Support")
			updateLegacyBridge(path)

			InstallPackageHelper("github.com/project-flogo/legacybridge")

			BuildAppModule(args)

			legacySupport = false
		} else {
			BuildAppModule(args)
		}

	}

}

func updateLegacyBridge(path string) {

	os.Chdir(Concat(path, "/src"))
	_, err := exec.Command("go", "get", "github.com/project-flogo/legacybridge").CombinedOutput()
	die(err)
	os.Chdir(path)
}

func InstallPackageHelper(args string) error {

	currdir, _ := os.Getwd()

	AddModToImport(args, currdir) //Edit the imports.go file
	os.Chdir(Concat(currdir, "/src"))
	//Download all the modules. This is just to ensure all packages are downloaded before go build.
	cliCmd, err := exec.Command("go", "mod", "download").CombinedOutput()
	if err != nil {
		RemoveModFromImport(args, currdir)

		fmt.Println(string(cliCmd))

		log.Fatal(err)
	}
	die(err)
	cliCmd, err = exec.Command("go", "mod", "tidy").CombinedOutput()
	if err != nil {
		RemoveModFromImport(args, currdir)

		fmt.Println(string(cliCmd))

		log.Fatal(err)
	}
	die(err)
	os.Chdir(currdir)
	return err
}

func BuildAppModule(args string) {

	currdir, _ := os.Getwd()
	os.Chdir(Concat(currdir, "/src"))

	//Build the modules.
	cliCmd, err := exec.Command("go", "build").CombinedOutput()

	if err != nil {

		fmt.Println(string(cliCmd))
		os.Chdir(currdir)
		RemoveModFromImport(args, currdir)

		log.Fatal(err)
	}

	os.Chdir(currdir)
	fmt.Println("Module Successfully Installed")

}

func AddModToImport(pkg string, fpath string) {

	byteArray, err := ioutil.ReadFile(Concat(fpath, "/src/imports.go"))

	die(err)

	text := string(byteArray)

	index := strings.Index(text, ")")

	err = ioutil.WriteFile(Concat(fpath, "/src/imports.go"), []byte(Concat(text[:index-1], "\n _ \"", pkg, "\" \n", ")")), 0)

	die(err)
}

func RemoveModFromImport(pkg string, fpath string) {

	byteArray, err := ioutil.ReadFile(Concat(fpath, "/src/imports.go"))

	die(err)

	text := string(byteArray)

	newText := strings.Replace(text, Concat(" _ \"", pkg, "\" "), "", -1)
	err = ioutil.WriteFile(Concat(fpath, "/src/imports.go"), []byte(newText), 0)

	die(err)
}

func checkCurrDir() bool {
	currDir, _ := os.Getwd()

	_, err := os.Stat(Concat(currDir, "/src/main.go"))
	_, err1 := os.Stat(Concat(currDir, "/src/imports.go"))
	_, err2 := os.Stat(Concat(currDir, "/src/go.mod"))

	return !(os.IsNotExist(err) && os.IsNotExist(err1) && os.IsNotExist(err2))
}
