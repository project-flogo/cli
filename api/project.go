package api

import (
	"fmt"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
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

	if runtime.GOOS == "windows" || GOOSENV == "windows" {
		execPath = filepath.Join(p.binDir, p.appName+".exe")
	} else {
		execPath = filepath.Join(p.binDir, p.appName)
	}

	return execPath
}

func (p *appProjectImpl) GetPath(pkg string) (string, error) {

	return p.dm.GetPath(pkg)
}

func (p *appProjectImpl) AddImports(ignoreError bool, imports ...string) error {

	importsFile := filepath.Join(p.SrcDir(), fileImportsGo)

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, importsFile, nil, parser.ImportsOnly)
	if err != nil {
		return err
	}

	for _, impPath := range imports {
		err := p.DepManager().AddDependency(impPath, "", true)
		if err != nil {
			if ignoreError {
				fmt.Printf("Warning: unable to install %s\n", impPath)
				continue
			}
			return err
		}
		util.AddImport(fset, file, impPath)
	}

	f, err := os.Create(importsFile)
	defer f.Close()
	if err := printer.Fprint(f, fset, file); err != nil {
		return err
	}

	//p.dm.Finalize()

	return nil
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
