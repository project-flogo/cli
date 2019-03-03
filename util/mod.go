package util

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/msoap/byline"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type DepManager interface {
	Init() error
	AddDependency(flogoImport Import) error
	GetPath(flogoImport Import) (string, error)
	AddLocalContribForBuild() error
	InstallLocalPkg(string, string)
}

func NewDepManager(sourceDir string) DepManager {
	return &ModDepManager{srcDir: sourceDir, localMods: make(map[string]string)}
}

type ModDepManager struct {
	srcDir    string
	localMods map[string]string
}

func (m *ModDepManager) Init() error {

	err := ExecCmd(exec.Command("go", "mod", "init", "main"), m.srcDir)
	if err == nil {
		return err
	}

	return nil
}

func (m *ModDepManager) AddDependency(flogoImport Import) error {

	// use "go mod edit" (instead of "go get") as first method
	err := ExecCmd(exec.Command("go", "mod", "edit", "-require", flogoImport.GoModImportPath()), m.srcDir)
	if err != nil {
		return err
	}

	// force resolution
	// TODO: add a flag to skip download and perform download later (useful in 'flogo create' command for instance)
	err = ExecCmd(exec.Command("go", "mod", "download", flogoImport.ModulePath()), m.srcDir)

	if err != nil {
		// if the resolution fails and the Flogo import is "classic"
		// (meaning it does not separate module path from Go import path):
		// 1. remove the import manually ("go mod edit -droprequire") would fail
		// 2. try with "go get" instead
		if flogoImport.IsClassic() {
			m.RemoveImport(flogoImport)

			err = ExecCmd(exec.Command("go", "get", flogoImport.GoGetImportPath()), m.srcDir)
		}
	}

	if err != nil {
		return err
	}

	return nil
}

// GetPath gets the path of where the
func (m *ModDepManager) GetPath(flogoImport Import) (string, error) {

	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	pkg := flogoImport.ModulePath()

	path, ok := m.localMods[pkg]
	if ok && path != "" {

		return path, nil
	}
	defer os.Chdir(currentDir)

	os.Chdir(m.srcDir)

	file, err := os.Open(filepath.Join(m.srcDir, "go.mod"))
	defer file.Close()

	var pathForPartial string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		line := scanner.Text()
		reqComponents := strings.Fields(line)
		//It is the line in the go.mod which is not useful, so ignore.
		if len(reqComponents) < 2 || (reqComponents[0] == "require" && reqComponents[1] == "(") {
			continue
		}

		//typically package is 1st component and  version is the 2nd component
		reqPkg := reqComponents[0]
		version := reqComponents[1]
		if reqComponents[0] == "require" {
			//starts with require, so package is 2nd component and  version is the 3rd component
			reqPkg = reqComponents[1]
			version = reqComponents[2]
		}

		if strings.HasPrefix(pkg, reqPkg) {

			hasFull := strings.Contains(line, pkg)
			tempPath := strings.Split(reqPkg, "/")

			tempPath = toLower(tempPath)
			lastIdx := len(tempPath) - 1

			tempPath[lastIdx] = tempPath[lastIdx] + "@" + version

			pkgPath := filepath.Join(tempPath...)

			if !hasFull {
				remaining := pkg[len(reqPkg):]
				tempPath = strings.Split(remaining, "/")
				remainingPath := filepath.Join(tempPath...)

				pathForPartial = filepath.Join(os.Getenv("GOPATH"), "pkg", "mod", pkgPath, remainingPath)
			} else {
				return filepath.Join(os.Getenv("GOPATH"), "pkg", "mod", pkgPath, flogoImport.RelativeImportPath()), nil
			}
		}
	}
	return pathForPartial, nil
}

func (m *ModDepManager) RemoveImport(flogoImport Import) error {

	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	modulePath := flogoImport.ModulePath()

	defer os.Chdir(currentDir)

	os.Chdir(m.srcDir)

	file, err := os.Open(filepath.Join(m.srcDir, "go.mod"))
	if err != nil {
		return err
	}
	defer file.Close()

	modulePath = strings.Replace(modulePath, "/", "\\/", -1)
	modulePath = strings.Replace(modulePath, ".", "\\.", -1)
	importRegex := regexp.MustCompile(`\s*` + modulePath + `\s+` + flogoImport.Version() + `.*`)

	lr := byline.NewReader(file)

	lr.MapString(func(line string) string {
		if importRegex.MatchString(line) {
			return ""
		} else {
			return line
		}
	})

	updatedGoMod, err := lr.ReadAll()
	if err != nil {
		return err
	}

	file, err = os.Create(filepath.Join(m.srcDir, "go.mod"))
	if err != nil {
		return err
	}
	defer file.Close()

	file.Write(updatedGoMod)

	return nil
}

//This function converts capotal letters in package name
// to !(smallercase). Eg C => !c . As this is the way
// go.mod saves every repository in the $GOPATH/pkg/mod.
func toLower(s []string) []string {
	result := make([]string, len(s))
	for i := 0; i < len(s); i++ {
		var b bytes.Buffer
		for _, c := range s[i] {
			if c >= 65 && c <= 90 {
				b.WriteRune(33)
				b.WriteRune(c + 32)
			} else {
				b.WriteRune(c)
			}
		}
		result[i] = b.String()
	}
	return result
}

var verbose = false

func SetVerbose(enable bool) {
	verbose = enable
}

func Verbose() bool {
	return verbose
}

func ExecCmd(cmd *exec.Cmd, workingDir string) error {

	if workingDir != "" {
		cmd.Dir = workingDir
	}

	var out bytes.Buffer

	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		cmd.Stdout = nil
		cmd.Stderr = &out
	}

	err := cmd.Run()

	if err != nil {
		return fmt.Errorf(string(out.Bytes()))
	}

	return nil
}

func (m *ModDepManager) AddLocalContribForBuild() error {

	text, err := ioutil.ReadFile(filepath.Join(m.srcDir, "go.mod"))
	if err != nil {
		return err
	}
	data := string(text)

	index := strings.Index(data, "replace")
	if index != -1 {
		localModules := strings.Split(data[index-1:], "\n")

		for _, val := range localModules {
			if val != "" {
				mods := strings.Split(val, " ")
				//If the length of mods is more than 4 it contains the versions of package
				//so it is stating to use different version of pkg rather than
				// the local pkg.
				if len(mods) < 5 {

					m.localMods[mods[1]] = mods[3]
				}

			}

		}
		return nil
	}
	return nil
}

func (m *ModDepManager) InstallLocalPkg(pkg1 string, pkg2 string) {

	m.localMods[pkg1] = pkg2

	f, err := os.OpenFile(filepath.Join(m.srcDir, "go.mod"), os.O_APPEND|os.O_WRONLY, 0777)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err = f.WriteString(fmt.Sprintf("replace %v => %v", pkg1, pkg2)); err != nil {
		panic(err)
	}

}

type Resp struct {
	Name string `json:"name"`
}

func getLatestVersion(path string) string {

	//To get the latest version number use the  GitHub API.
	resp, err := http.Get("https://api.github.com/repos/TIBCOSoftware/flogo-contrib/releases/latest")
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var result Resp

	json.Unmarshal(body, &result)

	return result.Name

}
