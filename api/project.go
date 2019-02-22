package api

import (
	"encoding/json"
	"fmt"
	"github.com/project-flogo/core/app" // dependency to core ensures the CLI always uses an up-to-date struct for JSON manipulation (this dependency already exists implicitly in the "flogo create" command)
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/project-flogo/cli/common"
	"github.com/project-flogo/cli/util"
)

const (
	flogoCoreRepo = "github.com/project-flogo/core"
	fileFlogoJson = "flogo.json"
	fileMainGo    = "main.go"
	fileImportsGo = "imports.go"
	dirSrc        = "src"
	dirBin        = "bin"
)

var GOOSENV = os.Getenv("GOOS")

type appProjectImpl struct {
	appDir  string
	appName string
	srcDir  string
	binDir  string
	dm      util.DepManager
}

func NewAppProject(appDir string) common.AppProject {
	project := &appProjectImpl{appDir: appDir}
	project.srcDir = filepath.Join(appDir, dirSrc)
	project.binDir = filepath.Join(appDir, dirBin)
	project.dm = util.NewDepManager(project.srcDir)
	project.appName = filepath.Base(appDir)
	return project
}

func (p *appProjectImpl) Validate() error {
	_, err := os.Stat(filepath.Join(p.appDir, fileFlogoJson))
	if os.IsNotExist(err) {
		return fmt.Errorf("not a valid flogo app project directory, missing flogo.json")
	}

	_, err = os.Stat(p.srcDir)
	if os.IsNotExist(err) {
		return fmt.Errorf("not a valid flogo app project directory, missing 'src' diretory")
	}

	_, err = os.Stat(filepath.Join(p.srcDir, fileImportsGo))
	if os.IsNotExist(err) {
		return fmt.Errorf("flogo app directory corrupt, missing 'src/imports.go' file")
	}

	_, err = os.Stat(filepath.Join(p.srcDir, "go.mod"))
	if os.IsNotExist(err) {
		return fmt.Errorf("flogo app directory corrupt, missing 'src/go.mod' file")
	}

	return nil
}

func (p *appProjectImpl) Name() string {
	return p.appName
}

func (p *appProjectImpl) Dir() string {
	return p.appDir
}

func (p *appProjectImpl) BinDir() string {
	return p.binDir
}

func (p *appProjectImpl) SrcDir() string {
	return p.srcDir
}

func (p *appProjectImpl) DepManager() util.DepManager {
	return p.dm
}

func (p *appProjectImpl) Executable() string {

	var execPath string

	execPath = filepath.Join(p.binDir, p.appName)

	if GOOSENV == "windows" || (runtime.GOOS == "windows" && GOOSENV == "") {
		// env or cross platform is windows
		execPath = filepath.Join(p.binDir, p.appName+".exe")
	}

	return execPath
}

func (p *appProjectImpl) GetPath(flogoImport util.Import) (string, error) {
	return p.dm.GetPath(flogoImport)
}

func (p *appProjectImpl) addImportsInGo(ignoreError bool, imports ...util.Import) error {
	importsFile := filepath.Join(p.SrcDir(), fileImportsGo)

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, importsFile, nil, parser.ImportsOnly)
	if err != nil {
		return err
	}

	for _, i := range imports {
		err := p.DepManager().AddDependency(i, true)
		if err != nil {
			if ignoreError {
				fmt.Printf("Warning: unable to install %s\n", i)
				continue
			}
			return err
		}
		util.AddImport(fset, file, i.ImportPath())
	}

	f, err := os.Create(importsFile)
	defer f.Close()
	if err := printer.Fprint(f, fset, file); err != nil {
		return err
	}

	//p.dm.Finalize()

	// using "go mod verify" can solve some dependencies conflicts (when contributions are dependent of others)
	err = util.ExecCmd(exec.Command("go", "mod", "verify"), p.srcDir)
	if err != nil {
		return err
	}

	return nil
}

func (p *appProjectImpl) addImportsInJson(ignoreError bool, imports ...util.Import) error {
	appDescriptorFile := filepath.Join(p.appDir, fileFlogoJson)
	appDescriptorJsonFile, err := os.Open(appDescriptorFile)
	if err != nil {
		return err
	}
	defer appDescriptorJsonFile.Close()

	appDescriptorData, err := ioutil.ReadAll(appDescriptorJsonFile)
	if err != nil {
		return err
	}

	var appDescriptor app.Config
	json.Unmarshal([]byte(appDescriptorData), &appDescriptor)

	// list existing imports in JSON to avoid duplicates
	existingImports := make(map[string]bool)
	jsonImports, _ := util.ParseImports(appDescriptor.Imports)
	for _, e := range jsonImports {
		existingImports[e.CanonicalImport()] = true
	}

	for _, i := range imports {
		if _, ok := existingImports[i.CanonicalImport()]; !ok {
			appDescriptor.Imports = append(appDescriptor.Imports, i.CanonicalImport())
		}
	}

	appDescriptorUpdated, err := json.MarshalIndent(appDescriptor, "", "  ")
	if err != nil {
		return err
	}

	appDescriptorUpdatedJson := string(appDescriptorUpdated)

	err = ioutil.WriteFile(appDescriptorFile, []byte(appDescriptorUpdatedJson), 0644)
	if err != nil {
		return err
	}

	return nil
}

func (p *appProjectImpl) AddImports(ignoreError bool, imports ...util.Import) error {
	err := p.addImportsInGo(ignoreError, imports...) // begin with Go imports as they are more likely to fail
	if err != nil {
		return err
	}
	err = p.addImportsInJson(ignoreError, imports...) // adding imports in JSON after Go imports ensure the flogo.json is self-sufficient

	return err
}

func (p *appProjectImpl) RemoveImports(imports ...string) error {

	importsFile := filepath.Join(p.SrcDir(), fileImportsGo)

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, importsFile, nil, parser.ImportsOnly)
	if err != nil {
		return err
	}

	for _, impPath := range imports {
		util.DeleteImport(fset, file, impPath)
	}

	f, err := os.Create(importsFile)
	defer f.Close()
	if err := printer.Fprint(f, fset, file); err != nil {
		return err
	}

	return nil
}
