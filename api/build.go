package api

import (
	"fmt"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/project-flogo/cli/common"
	"github.com/project-flogo/cli/util"
)

const (
	fileEmbeddedAppGo string = "embeddedapp.go"
)

func BuildProject(project common.AppProject, options common.BuildOptions) error {

	err := project.DepManager().AddReplacedContribForBuild()
	if err != nil {
		return err
	}

	buildPreProcessors := common.BuildPreProcessors()

	if len(buildPreProcessors) > 0 {
		for _, processor := range buildPreProcessors {
			err = processor.DoPreProcessing(project,options)
			if err != nil {
				return err
			}
		}
	}

	var builder common.Builder
	embedConfig := options.EmbedConfig

	if options.Shim != "" {
		builder = &ShimBuilder{shim:options.Shim}
		embedConfig = true
	} else {
		builder = &AppBuilder{}
	}

	if embedConfig {
		err = createEmbeddedAppGoFile(project)
		if err != nil {
			return err
		}
	} else {
		err = cleanupEmbeddedAppGoFile(project)
		if err != nil {
			return err
		}
	}

	if options.OptimizeImports {
		if Verbose() {
			fmt.Println("Optimizing imports...")
		}
		err := optimizeImports(project)
		defer restoreImports(project)

		if err != nil {
			return err
		}
	}

	err = builder.Build(project)
	if err != nil {
		return err
	}

	buildPostProcessors := common.BuildPostProcessors()

	if len(buildPostProcessors) > 0 {
		for _, processor := range buildPostProcessors {
			err = processor.DoPostProcessing(project)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func cleanupEmbeddedAppGoFile(project common.AppProject) error {
	embedSrcPath := filepath.Join(project.SrcDir(), fileEmbeddedAppGo)

		if _, err := os.Stat(embedSrcPath); err == nil {
			if Verbose() {
				fmt.Println("Removing embed configuration")
			}
			err = os.Remove(embedSrcPath)
			if err != nil {
				return err
			}
		}
		return nil
}

func createEmbeddedAppGoFile(project common.AppProject) error {

	embedSrcPath := filepath.Join(project.SrcDir(), fileEmbeddedAppGo)

	if Verbose() {
		fmt.Println("Embedding configuration in application...")
	}

	buf, err := ioutil.ReadFile(filepath.Join(project.Dir(), fileFlogoJson))
	if err != nil {
		return err
	}
	flogoJSON := string(buf)

	tplFile := tplEmbeddedAppGoFile
	if !isNewMain(project) {
		tplFile = tplEmbeddedAppOldGoFile
	}

	engineJSON := ""

	if util.FileExists(filepath.Join(project.Dir(), fileEngineJson)) {
		buf, err = ioutil.ReadFile(filepath.Join(project.Dir(), fileEngineJson))
		if err != nil {
			return err
		}

		engineJSON = string(buf)
	}

	data := struct {
		FlogoJSON string
		EngineJSON string
	}{
		flogoJSON,
		engineJSON,
	}

	f, err := os.Create(embedSrcPath)
	if err != nil {
		return err
	}
	RenderTemplate(f, tplFile, &data)
	_ = f.Close()

	return nil
}

func isNewMain(project common.AppProject) bool {
	mainGo := filepath.Join(project.SrcDir(), fileMainGo)
	buf, err := ioutil.ReadFile(mainGo)
	if err == nil {
		mainCode := string(buf)
		return strings.Contains(mainCode, "cfgEngine")

	}

	return false
}


var tplEmbeddedAppGoFile = `// Do not change this file, it has been generated using flogo-cli
// If you change it and rebuild the application your changes might get lost
package main

// embedded flogo app descriptor file
const flogoJSON string = ` + "`{{.FlogoJSON}}`" + `
const engineJSON string = ` + "`{{.EngineJSON}}`" + `

func init () {
	cfgJson = flogoJSON
	cfgEngine = engineJSON
}
`

var tplEmbeddedAppOldGoFile = `// Do not change this file, it has been generated using flogo-cli
// If you change it and rebuild the application your changes might get lost
package main

// embedded flogo app descriptor file
const flogoJSON string = ` + "`{{.FlogoJSON}}`" + `

func init () {
	cfgJson = flogoJSON
}
`

func initMain(project common.AppProject, backupMain bool) error {

	//backup main if it exists
	mainGo := filepath.Join(project.SrcDir(), fileMainGo)
	mainGoBak := filepath.Join(project.SrcDir(), fileMainGo+".bak")

	if backupMain {
		if _, err := os.Stat(mainGo); err == nil {
			err = os.Rename(mainGo, mainGoBak)
			if err != nil {
				return err
			}
		} else if _, err := os.Stat(mainGoBak); err != nil {
			return fmt.Errorf("project corrupt, main missing")
		}
	} else {
		if _, err := os.Stat(mainGoBak); err == nil {
			err = os.Rename(mainGoBak, mainGo)
			if err != nil {
				return err
			}
		} else if _, err := os.Stat(mainGo); err != nil {
			return fmt.Errorf("project corrupt, main missing")
		}
	}

	return nil
}

func optimizeImports(project common.AppProject) error {

	appImports, err := util.GetAppImports(filepath.Join(project.Dir(), fileFlogoJson), project.DepManager(), true)
	if err != nil {
		return err
	}

	var unused []util.Import
	appImports.GetAllImports()
	for _, impDetails := range appImports.GetAllImportDetails() {
		if !impDetails.Referenced() && impDetails.IsCoreContrib() {
			unused = append(unused, impDetails.Imp)
		}
	}

	importsFile := filepath.Join(project.SrcDir(), fileImportsGo)
	importsFileOrig := filepath.Join(project.SrcDir(), fileImportsGo+".orig")

	err = util.CopyFile(importsFile, importsFileOrig)
	if err != nil {
		return err
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, importsFile, nil, parser.ImportsOnly)
	if err != nil {
		return err
	}

	for _, i := range unused {
		if Verbose() {
			fmt.Printf("  Removing Import: %s\n", i.GoImportPath())
		}
		util.DeleteImport(fset, file, i.GoImportPath())
	}

	f, err := os.Create(importsFile)
	defer f.Close()
	if err := printer.Fprint(f, fset, file); err != nil {
		return err
	}

	return nil
}

func restoreImports(project common.AppProject) {

	importsFile := filepath.Join(project.SrcDir(), fileImportsGo)
	importsFileOrig := filepath.Join(project.SrcDir(), fileImportsGo+".orig")

	if _, err := os.Stat(importsFileOrig); err == nil {
		err = util.CopyFile(importsFileOrig, importsFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error restoring imports file '%s': %v\n", importsFile, err)
			return
		}

		var err = os.Remove(importsFileOrig)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error removing backup imports file '%s': %v\n", importsFileOrig, err)
			fmt.Fprintf(os.Stderr, "Manually remove backup imports file '%s'\n", importsFileOrig)
		}
	}
}
