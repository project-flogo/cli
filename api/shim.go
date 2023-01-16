package api

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/project-flogo/cli/common"
	"github.com/project-flogo/cli/util"
)

const (
	fileShimSupportGo string = "shim_support.go"
	fileShimGo        string = "shim.go"
	fileBuildGo       string = "build.go"
	fileMakefile      string = "Makefile"
	dirShim           string = "shim"
)

var fileSampleShimSupport = filepath.Join("examples", "engine", "shim", fileShimSupportGo)

var flogoImportPattern = regexp.MustCompile(`^(([^ ]*)[ ]+)?([^@:]*)@?([^:]*)?:?(.*)?$`)


type ShimBuilder struct {
	appBuilder common.Builder
	shim string
}

func (sb *ShimBuilder) Build(project common.AppProject) error {

	err := backupMain(project)
	if err != nil {
		return err
	}

	defer shimCleanup(project)

	err = createShimSupportGoFile(project)
	if err != nil {
		return err
	}

	if Verbose() {
		fmt.Println("Preparing shim...")
	}
	built, err := prepareShim(project, sb.shim)
	if err != nil {
		return err
	}

	if !built {
		fmt.Println("Using go build to build shim...")

		err := simpleGoBuild(project)
		if err != nil {
			return err
		}
	}

	return nil
}

func prepareShim(project common.AppProject, shim string) (bool, error) {

	buf, err := ioutil.ReadFile(filepath.Join(project.Dir(), fileFlogoJson))
	if err != nil {
		return false, err
	}

	flogoJSON := string(buf)

	descriptor, err := util.ParseAppDescriptor(flogoJSON)
	if err != nil {
		return false, err
	}

	err = registerImports(project, descriptor)
	if err != nil {
		return false, err
	}

	var ref string = ""
	var triggerFound bool = false

	for _, trgCfg := range descriptor.Triggers {

		if trgCfg.Id != shim {
			continue 
		}

		ref, triggerFound = GetAliasRef("flogo:trigger", trgCfg.Ref)

		if !triggerFound {
			break
		} 
		
		return prepareShimBuildOrMake(project, ref)

	}

	if !triggerFound {

		for _, resourceCfg := range descriptor.Resources {

			for _, taskConfig := range resourceCfg.Data.Tasks {

				if taskConfig.Id != shim {
					continue
				}

				var found bool
				ref, found = GetAliasRef("flogo:activity", taskConfig.Activity.Ref)
				
				if !found {
					break
				} 

				return prepareShimBuildOrMake(project, ref)

			}

		}

	}

	return false, fmt.Errorf("unable to to find shim trigger or activity: %s", shim)
}

func prepareShimBuildOrMake(project common.AppProject, ref string) (bool, error) {

	refImport, err := util.NewFlogoImportFromPath(ref)
	if err != nil {
		return false, err
	}

	impPath, err := project.GetPath(refImport)
	if err != nil {
		return false, err
	}

	// var shimFilePath string

	// shimFilePath := filepath.Join(impPath, dirShim, fileShimGo)
	shimDir := filepath.Join(impPath, dirShim)

	// if _, err := os.Stat(shimFilePath); err == nil {
	if _, err := os.Stat(shimDir); err == nil {

		// err = util.CopyFile(shimFilePath, filepath.Join(project.SrcDir(), fileShimGo))
		// if err != nil {
		// 	return false, err
		// }

		goBuildFilePath := filepath.Join(shimDir, fileBuildGo)

		makefilePath := filepath.Join(shimDir, fileMakefile)

		if _, err := os.Stat(goBuildFilePath); err == nil {
			fmt.Println("Using build.go to build shim......")

			err = util.CopyFile(goBuildFilePath, filepath.Join(project.SrcDir(), fileBuildGo))
			if err != nil {
				return false, err
			}

			// Execute go run gobuild.go
			err = util.ExecCmd(exec.Command("go", "run", fileBuildGo), project.SrcDir())
			if err != nil {
				return false, err
			}

			return true, nil

		} else if _, err := os.Stat(makefilePath); err == nil {

			//look for Makefile and execute it
			fmt.Println("Using make file to build shim...")

			err = util.CopyFile(makefilePath, filepath.Join(project.SrcDir(), fileMakefile))
			if err != nil {
				return false, err
			}

			if Verbose() {
				fmt.Println("Make File:", makefilePath)
			}

			// Execute make
			cmd := exec.Command("make", "-C", project.SrcDir())
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Env = util.ReplaceEnvValue(os.Environ(), "GOPATH", project.Dir())

			err = cmd.Run()
			if err != nil {
				return false, err
			}

			return true, nil
			
		} 
	}

	return  false, nil
}

func shimCleanup(project common.AppProject) {

	if Verbose() {
		fmt.Println("Cleaning up shim support files...")
	}

	err := util.DeleteFile(filepath.Join(project.SrcDir(), fileShimSupportGo))
	if err != nil {
		fmt.Printf("Unable to delete: %s", fileShimSupportGo)
	}
	err = util.DeleteFile(filepath.Join(project.SrcDir(), fileShimGo))
	if err != nil {
		fmt.Printf("Unable to delete: %s", fileShimGo)
	}
	err = util.DeleteFile(filepath.Join(project.SrcDir(), fileBuildGo))
	if err != nil {
		fmt.Printf("Unable to delete: %s", fileBuildGo)
	}
}

func createShimSupportGoFile(project common.AppProject) error {

	shimSrcPath := filepath.Join(project.SrcDir(), fileShimSupportGo)

	if Verbose() {
		fmt.Println("Creating shim support files...")
	}

	flogoCoreImport, err := util.NewFlogoImportFromPath(flogoCoreRepo)
	if err != nil {
		return err
	}

	corePath, err := project.GetPath(flogoCoreImport)
	if err != nil {
		return err
	}

	bytes, err := ioutil.ReadFile(filepath.Join(corePath, fileSampleShimSupport))
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(shimSrcPath, bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func registerImports(project common.AppProject, appDesc *util.FlogoAppDescriptor) error {

	for _, anImport := range appDesc.Imports {
		err := registerImport(project, anImport)
		if err != nil {
			return err
		}
	}

	return nil
}

func registerImport(project common.AppProject, anImport string) error {

	matches := flogoImportPattern.FindStringSubmatch(anImport)

	parts := strings.Split(matches[3], " ")

	var alias string
	var ref string
	numParts := len(parts)
	if numParts == 1 {
		ref = parts[0]
		alias = path.Base(ref)

	} else if numParts == 2 {
		alias = parts[0]
		ref = parts[1]
	} else {
		return fmt.Errorf("invalid import %s", anImport)
	}

	if alias == "" || ref == "" {
		return fmt.Errorf("invalid import %s", anImport)
	}

	ct, err := util.GetContribType(project.DepManager(), ref)
	if err != nil {
		return err
	}

	if ct != "" {
		RegisterAlias(ct, alias, ref)
	}

	return nil
}

var aliases = make(map[string]map[string]string)

func RegisterAlias(contribType string, alias, ref string) {

	aliasToRefMap, exists := aliases[contribType]
	if !exists {
		aliasToRefMap = make(map[string]string)
		aliases[contribType] = aliasToRefMap
	}

	aliasToRefMap[alias] = ref
}

func GetAliasRef(contribType string, alias string) (string, bool) {
	if alias == "" {
		return "", false
	}

	if alias[0] == '#' {
		alias = alias[1:]
	}
	aliasToRefMap, exists := aliases[contribType]
	if !exists {
		return "", false
	}

	ref, exists := aliasToRefMap[alias]
	if !exists {
		return "", false
	}

	return ref, true
}
