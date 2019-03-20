package api

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/project-flogo/cli/common"
	"github.com/project-flogo/cli/util"
)

const (
	fileEmbeddedAppGo string = "embeddedapp.go"
)

type BuildOptions struct {
	OptimizeImports bool
	EmbedConfig     bool
	Shim            string
}


func BuildProject(project common.AppProject, options BuildOptions) error {

	err := project.DepManager().AddLocalContribForBuild()
	if err != nil {
		return err
	}

	useShim := options.Shim != ""

	err = createEmbeddedAppGoFile(project, options.EmbedConfig || useShim)
	if err != nil {
		return err
	}

	err = createShimSupportGoFile(project, useShim)
	if err != nil {
		return err
	}

	err = initMain(project, useShim)
	if err != nil {
		return err
	}

	if useShim {
		buildExist, err := prepareShim(project, options.Shim)
		if err != nil {
			return err
		}
		if buildExist {
			return nil
		}

	}

	err = util.ExecCmd(exec.Command("go", "build"), project.SrcDir())
	if err != nil {
		fmt.Println("Error in building", project.SrcDir())
		return err
	}

	// assume linux/darwin env or cross platform by default
	exe := "main"

	if GOOSENV == "windows" || (runtime.GOOS == "windows" && GOOSENV == "") {
		// env or cross platform is windows
		exe = "main.exe"
	}

	exePath := filepath.Join(project.SrcDir(), exe)

	if common.Verbose() {
		fmt.Println("Path to exe is ", exePath)
	}
	if _, err := os.Stat(exePath); err == nil {
		finalExePath := project.Executable()
		err = os.MkdirAll(filepath.Dir(finalExePath), os.ModePerm)
		if err != nil {
			return err
		}
		err = os.Rename(exePath, project.Executable())
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("failed to build application, run with --verbose to see details")
	}

	return nil
}

func createEmbeddedAppGoFile(project common.AppProject, create bool) error {

	embedSrcPath := filepath.Join(project.SrcDir(), fileEmbeddedAppGo)

	if !create {
		if _, err := os.Stat(embedSrcPath); err == nil {
			err = os.Remove(embedSrcPath)
			if err != nil {
				return err
			}
		}
		return nil
	}

	buf, err := ioutil.ReadFile(filepath.Join(project.Dir(), fileFlogoJson))
	if err != nil {
		return err
	}

	flogoJSON := string(buf)

	data := struct {
		FlogoJSON string
	}{
		flogoJSON,
	}

	f, err := os.Create(embedSrcPath)
	if err != nil {
		return err
	}
	RenderTemplate(f, tplEmbeddedAppGoFile, &data)
	f.Close()

	return nil
}

var tplEmbeddedAppGoFile = `// Do not change this file, it has been generated using flogo-cli
// If you change it and rebuild the application your changes might get lost
package main

// embedded flogo app descriptor file
const flogoJSON string = ` + "`{{.FlogoJSON}}`" + `

func init () {
	cfgJson = flogoJSON
}
`

func copyFile(srcFilePath, destFilePath string) error {

	bytes, err := ioutil.ReadFile(srcFilePath)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(destFilePath, bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

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
