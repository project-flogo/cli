package commands

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var file bool
var core string

var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Flogo Cli lets you create with Flogo",
	Long:  `Flogo Cli create is great! `,
	Run: func(cmd *cobra.Command, args []string) {

		if file {
			if len(os.Args) <= 3 {
				fmt.Println("Enter file name")
				os.Exit(1)
			} else {
				fmt.Println("Building the Flogo.")
				//Check if file exists
				fileName := os.Args[len(os.Args)-1]
				CheckFile(fileName)

				CreateAppFolder(strings.Split(fileName, ".")[0])

				AddFiles(strings.Split(fileName, ".")[0], core)

				populateFilesFromSrc(strings.Split(fileName, ".")[0])

				listsOfRefs := GetRefsFromFile(fileName)

				for _, ref := range listsOfRefs {

					fmt.Println("Installing ", ref)
					currDir, _ := os.Getwd()
					//Move to the App folder.
					os.Chdir(Concat(currDir, "/", strings.Split(fileName, ".")[0]))
					currDir, _ = os.Getwd()
					//Edit imports file in the App Folder
					AddModToImport(ref, currDir)
				}
			}
		} else {

			CreateAppFolder(os.Args[len(os.Args)-1])

			AddFiles(os.Args[len(os.Args)-1], core)

			populateFilesFromCore(os.Args[len(os.Args)-1])

		}

	},
}

func init() {
	RootCmd.AddCommand(CreateCmd)
	CreateCmd.Flags().BoolVarP(&file, "file", "f", false, "Enter file")
	CreateCmd.Flags().StringVarP(&core, "core", "c", "", "Enter core version")
}
func CheckFile(args string) {
	if !strings.Contains(args, ".json") {
		fmt.Println("Please enter file name")
		os.Exit(1)
	}
}
func CreateAppFolder(args string) {

	dirName := strings.Split(args, ".")[0]

	err := os.Mkdir(dirName, os.ModePerm)

	if err != nil {
		log.Fatal(err)
	}
}

func GetRefsFromFile(args string) []string {

	var result []string

	file, err := os.Open(args)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "ref") {
			pkg := strings.Split(line, ":")[1]
			pkg = strings.TrimSpace(pkg)
			pkg = pkg[1 : len(pkg)-2]
			result = append(result, pkg)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return result
}
func AddFiles(dir string, core string) {
	//Get the Curr Dir.
	currDir, err := os.Getwd()

	//Add folders and files in the app folder
	err = os.Mkdir(Concat(currDir, "/", dir, "/bin"), os.ModePerm)
	err = os.Mkdir(Concat(currDir, "/", dir, "/src"), os.ModePerm)

	_, err = os.Create(Concat(currDir, "/", dir, "/src/imports.go"))

	_, err = os.Create(Concat(currDir, "/", dir, "/src/main.go"))

	//Move to src/ in App to initialize mod file.
	os.Chdir(Concat(currDir, "/", dir, "/src"))

	cliCmd, err := exec.Command("go", "mod", "init", "main").Output()
	cliCmd, err = exec.Command("go", "mod", "edit", "-require", "github.com/sirupsen/logrus@v1.1.1").Output()
	if len(core) > 1 {
		if core == "master" {
			cliCmd, err = exec.Command("go", "get", "github.com/project-flogo/core@master").CombinedOutput()
		} else {
			cliCmd, err = exec.Command("go", "mod", "edit", "-require", Concat("github.com/project-flogo/core@", core)).Output()
		}
	} else {
		cliCmd, err = exec.Command("go", "get", "github.com/project-flogo/core").CombinedOutput()
	}

	if err != nil {
		fmt.Println(string(cliCmd))

		log.Fatal(err)
	}
	os.Chdir(Concat(currDir, "/", dir))
	if err != nil {
		fmt.Println(string(cliCmd))
		log.Fatal(err)
	}
	err = ioutil.WriteFile(Concat(currDir, "/", dir, "/src/imports.go"), []byte("package main\n import (\n _ \"github.com/project-flogo/core/app\" \n )"), 0644)
	if err != nil {
		log.Fatal(err)
	}

}

func populateFilesFromCore(path string) {
	//Edit Import
	currDir, err := os.Getwd()

	filePath := Concat(os.Getenv("GOPATH"), "/pkg/mod/", getTruePath(currDir, "github.com/project-flogo/core"), "/examples/engine/imports.go")

	byteArray, err := ioutil.ReadFile(filePath)

	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(Concat(currDir, "/src/imports.go"), byteArray, 0644)
	if err != nil {
		log.Fatal(err)
	}
	//Edit Json
	filePath = Concat(os.Getenv("GOPATH"), "/pkg/mod/", getTruePath(currDir, "github.com/project-flogo/core"), "/examples/engine/flogo.json")

	byteArray, err = ioutil.ReadFile(filePath)

	_, err = os.Create(Concat(currDir, "/flogo.json"))

	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(Concat(currDir, "/flogo.json"), byteArray, 0644)
	if err != nil {
		log.Fatal(err)
	}
	//Edit main
	filePath = Concat(os.Getenv("GOPATH"), "/pkg/mod/", getTruePath(currDir, "github.com/project-flogo/core"), "/examples/engine/main.go")

	byteArray, err = ioutil.ReadFile(filePath)

	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(Concat(currDir, "/src/main.go"), byteArray, 0644)
	if err != nil {
		log.Fatal(err)
	}
	//Edit Mod
	f, err := os.OpenFile(Concat(currDir, "/src/go.mod"), os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString("\n replace github.com/Sirupsen/logrus v1.1.0 => github.com/sirupsen/logrus v1.1.0 \n replace github.com/TIBCOSoftware/flogo-lib v0.5.6 => github.com/TIBCOSoftware/flogo-lib v0.5.7-0.20181009194308-1fe2a7011501 \n"); err != nil {
		panic(err)
	}

}

func populateFilesFromSrc(path string) {
	currDir, err := os.Getwd()

	//Edit main
	filePath := Concat(os.Getenv("GOPATH"), "/pkg/mod/", getTruePath(currDir, "github.com/project-flogo/core"), "/examples/engine/main.go")

	byteArray, err := ioutil.ReadFile(filePath)

	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(Concat(currDir, "/src/main.go"), byteArray, 0644)
	if err != nil {
		log.Fatal(err)
	}
	//Edit Import

	err = ioutil.WriteFile(Concat(currDir, "/src/imports.go"), []byte("package main \n import ( \n _ \"os\" \n )"), 0644)
	if err != nil {
		log.Fatal(err)
	}
	//Edit Mod
	f, err := os.OpenFile(Concat(currDir, "/src/go.mod"), os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString("\n replace github.com/Sirupsen/logrus v1.1.0 => github.com/sirupsen/logrus v1.1.0 \n replace github.com/TIBCOSoftware/flogo-lib v0.5.6 => github.com/TIBCOSoftware/flogo-lib v0.5.7-0.20181009194308-1fe2a7011501 \n"); err != nil {
		panic(err)
	}

	//Copy Json
	err = os.Chdir("..")
	currDir, _ = os.Getwd()

	byteArray, err = ioutil.ReadFile(Concat("./", path, ".json"))
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(Concat(currDir, "/", path, "/", path, ".json"), byteArray, 0644)
	if err != nil {
		log.Fatal(err)
	}
	os.Chdir(currDir)
}
