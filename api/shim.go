package api

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
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

	for _, trgCfg := range descriptor.Triggers {
		if trgCfg.Id == shim {

			ref := trgCfg.Ref

			if trgCfg.Ref == "" {
				found := false
				ref, found = GetAliasRef("flogo:trigger", trgCfg.Type)
				if !found {
					return false, fmt.Errorf("unable to determine ref for trigger: %s", trgCfg.Id)
				}
			}

			refImport, err := util.NewFlogoImportFromPath(ref)
			if err != nil {
				return false, err
			}

			impPath, err := project.GetPath(refImport)
			if err != nil {
				return false, err
			}
			var shimFilePath string

			shimFilePath = filepath.Join(impPath, dirShim, fileShimGo)

			if _, err := os.Stat(shimFilePath); err == nil {

				err = copyFile(shimFilePath, filepath.Join(project.SrcDir(), fileShimGo))
				if err != nil {
					return false, err
				}

				// Check if this shim based trigger has a gobuild file. If the trigger has a gobuild
				// execute that file, otherwise check if there is a Makefile to execute
				goBuildFilePath := filepath.Join(impPath, dirShim, fileBuildGo)

				makefilePath := filepath.Join(shimFilePath, dirShim, fileMakefile)

				if _, err := os.Stat(goBuildFilePath); err == nil {
					fmt.Println("This trigger makes use of a go build file...")

					err = copyFile(goBuildFilePath, filepath.Join(project.SrcDir(), fileBuildGo))
					if err != nil {
						return false, err
					}

					// Execute go run gobuild.go
					err = util.ExecCmd(exec.Command("go", "run", fileBuildGo), project.SrcDir())
					if err != nil {
						return false, err
					}
				} else if _, err := os.Stat(makefilePath); err == nil {
					//look for Makefile and execute it
					fmt.Println("Make File:", makefilePath)

					err = copyFile(makefilePath, filepath.Join(project.SrcDir(), fileMakefile))
					if err != nil {
						return false, err
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
				} else {
					return false, nil
				}
			}

			break
		}
	}

	return true, nil
}

func createShimSupportGoFile(project common.AppProject, create bool) error {

	shimSrcPath := filepath.Join(project.SrcDir(), fileShimSupportGo)

	if !create {
		if _, err := os.Stat(shimSrcPath); err == nil {
			err = os.Remove(shimSrcPath)
			if err != nil {
				return err
			}
		}

		//
		//shimSrcPath := filepath.Join(project.SrcDir(), fileShimSupportGo)
		//
		//if _, err := os.Stat(shimSrcPath); err == nil {
		//	os.Remove(shimSrcPath)
		//	if err != nil {
		//		return err
		//	}
		//}
		return nil
	}

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

	parts := strings.Split(anImport, " ")

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

	ct, err := getContribType(project, ref)
	if err != nil {
		return err
	}

	if ct == "" {
		return fmt.Errorf("unable to determine contribution type for import: %s", anImport)
	}

	RegisterAlias(ct, alias, ref)
	return nil
}

func getContribType(project common.AppProject, ref string) (string, error) {

	refAsFlogoImport, err := util.NewFlogoImportFromPath(ref)
	if err != nil {
		return "", err
	}

	impPath, err := project.GetPath(refAsFlogoImport)
	if err != nil {
		return "", err
	}
	var descriptorPath string

	if _, err := os.Stat(filepath.Join(impPath, fileDescriptorJson)); err == nil {
		descriptorPath = filepath.Join(impPath, fileDescriptorJson)

	} else if _, err := os.Stat(filepath.Join(impPath, "activity.json")); err == nil {
		descriptorPath = filepath.Join(impPath, "activity.json")
	} else if _, err := os.Stat(filepath.Join(impPath, "trigger.json")); err == nil {
		descriptorPath = filepath.Join(impPath, "trigger.json")
	} else if _, err := os.Stat(filepath.Join(impPath, "action.json")); err == nil {
		descriptorPath = filepath.Join(impPath, "action.json")
	}

	if _, err := os.Stat(descriptorPath); descriptorPath != "" && err == nil {

		desc, err := util.ReadContribDescriptor(descriptorPath)
		if err != nil {
			return "", err
		}

		return desc.Type, nil
	}

	return "", nil
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
