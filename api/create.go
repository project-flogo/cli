package api

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

func CreateProject(file string, flag bool, corev string, appPath string) error {

	if appPath == "." {
		appPath, _ = os.Getwd()

	}

	return createProject(file, flag, corev, appPath)

}

func createProject(fileName string, flag bool, coreVersion string, appPath string) error {
	if flag {

		fmt.Println("Building the Flogo.")
		//Check if file exists

		err := CheckFile(fileName)
		if err != nil {
			return err
		}

		err = CreateAppFolder(strings.Split(fileName, ".")[0], appPath)
		if err != nil {
			return err
		}

		err = AddFiles(strings.Split(fileName, ".")[0], coreVersion, appPath)
		if err != nil {
			return err
		}

		err = populateFiles(strings.Split(fileName, ".")[0], Concat(appPath, "/", strings.Split(fileName, ".")[0]), true)
		if err != nil {
			return err
		}

		listsOfRefs := GetRefsFromFile(Concat(appPath, "/", fileName))

		for _, ref := range listsOfRefs {

			fmt.Println("Installing ", ref)
			//Move to the App folder.
			os.Chdir(Concat(appPath, "/", strings.Split(fileName, ".")[0]))
			currDir, _ := os.Getwd()
			//Edit imports file in the App Folder
			AddModToImport(ref, currDir)
		}

	} else {

		err := CreateAppFolder(fileName, appPath)
		if err != nil {
			return err
		}

		err = AddFiles(fileName, coreVersion, appPath)
		if err != nil {
			return err
		}

		err = populateFiles(fileName, Concat(appPath, "/", fileName), false)
		if err != nil {
			return err
		}

	}
	return nil
}

func CheckFile(args string) error {
	if !strings.Contains(args, ".json") {
		fmt.Println("Please enter file name")
		return errors.New("Please enter file name")
	}
	return nil
}
func CreateAppFolder(args string, path string) error {

	dirName := strings.Split(args, ".")[0]

	err := os.Mkdir(Concat(path, "/", dirName), os.ModePerm)

	if err != nil {
		return err
	}
	return nil
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
func AddFiles(dir string, core string, currDir string) error {
	//Get the Curr Dir.

	//Add folders and files in the app folder
	err := os.Mkdir(Concat(currDir, "/", dir, "/bin"), os.ModePerm)
	if err != nil {
		return err
	}
	err = os.Mkdir(Concat(currDir, "/", dir, "/src"), os.ModePerm)
	if err != nil {
		return err
	}

	_, err = os.Create(Concat(currDir, "/", dir, "/src/imports.go"))
	if err != nil {
		return err
	}

	_, err = os.Create(Concat(currDir, "/", dir, "/src/main.go"))
	if err != nil {
		return err
	}

	//Move to src/ in App to initialize mod file.
	os.Chdir(Concat(currDir, "/", dir, "/src"))

	cliCmd, err := exec.Command("go", "mod", "init", "main").Output()
	if err != nil {
		return err
	}

	cliCmd, err = exec.Command("go", "mod", "edit", "-require", "github.com/sirupsen/logrus@v1.1.1").Output()

	if len(core) > 1 {
		if core == "master" {
			cliCmd, err = exec.Command("go", "get", "github.com/project-flogo/core@master").CombinedOutput()
			if err != nil {
				return err
			}

		} else {
			cliCmd, err = exec.Command("go", "mod", "edit", "-require", Concat("github.com/project-flogo/core@", core)).Output()
			if err != nil {
				return err
			}

		}
		cliCmd, err = exec.Command("go", "get", "github.com/project-flogo/core").CombinedOutput()
		if err != nil {
			return err
		}
	}

	if err != nil {
		fmt.Println(string(cliCmd))

		return err
	}
	os.Chdir(Concat(currDir, "/", dir))
	if err != nil {
		fmt.Println(string(cliCmd))
		return err
	}
	err = ioutil.WriteFile(Concat(currDir, "/", dir, "/src/imports.go"), []byte("package main\n import (\n _ \"github.com/project-flogo/core/app\" \n )"), 0644)
	if err != nil {
		return err
	}

	return nil
}

func populateFiles(file string, currDir string, src bool) (err error) {

	//Is source present
	if src {
		//Edit Import
		AddModToImport("os", currDir)

		//Edit Json

		cliCmd, err := exec.Command("cp", Concat("../", file, ".json"), Concat(currDir, "/")).Output()
		if err != nil {
			fmt.Println(string(cliCmd))
			return err
		}

	} else {
		//Edit Import
		filePath := Concat(os.Getenv("GOPATH"), "/pkg/mod/", getTruePath(currDir, "github.com/project-flogo/core"), "/examples/engine/imports.go")

		byteArray, err := ioutil.ReadFile(filePath)

		if err != nil {
			return err
		}

		err = ioutil.WriteFile(Concat(currDir, "/src/imports.go"), byteArray, 0644)
		if err != nil {
			return err
		}

		//Edit Json
		filePath = Concat(os.Getenv("GOPATH"), "/pkg/mod/", getTruePath(currDir, "github.com/project-flogo/core"), "/examples/engine/flogo.json")

		byteArray, err = ioutil.ReadFile(filePath)

		_, err = os.Create(Concat(currDir, "/flogo.json"))

		if err != nil {
			return err
		}
		err = ioutil.WriteFile(Concat(currDir, "/flogo.json"), byteArray, 0644)
		if err != nil {
			return err
		}
	}

	//Edit main

	filePath := Concat(os.Getenv("GOPATH"), "/pkg/mod/", getTruePath(currDir, "github.com/project-flogo/core"), "/examples/engine/main.go")

	byteArray, err := ioutil.ReadFile(filePath)

	if err != nil {
		return err
	}
	err = ioutil.WriteFile(Concat(currDir, "/src/main.go"), byteArray, 0644)
	if err != nil {
		return err
	}

	//Edit Mod
	f, err := os.OpenFile(Concat(currDir, "/src/go.mod"), os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err = f.WriteString("\n replace github.com/Sirupsen/logrus v1.1.0 => github.com/sirupsen/logrus v1.1.0 \n replace github.com/TIBCOSoftware/flogo-lib v0.5.6 => github.com/TIBCOSoftware/flogo-lib v0.5.7-0.20181009194308-1fe2a7011501 \n"); err != nil {
		return err
	}
	return nil
}
