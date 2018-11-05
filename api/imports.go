package api

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/project-flogo/cli/common"
	"github.com/project-flogo/cli/util"
)

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

	path, err := project.GetPath(ref)
	if err != nil {
		return "", err
	}

	var descriptorPath string

	if _, err := os.Stat(filepath.Join(path, fileDescriptorJson)); err == nil {
		descriptorPath = filepath.Join(path, fileDescriptorJson)
	} else if _, err := os.Stat(filepath.Join(path, "activity.json")); err == nil {
		descriptorPath = filepath.Join(path, "activity.json")
	} else if _, err := os.Stat(filepath.Join(path, "trigger.json")); err == nil {
		descriptorPath = filepath.Join(path, "trigger.json")
	} else if _, err := os.Stat(filepath.Join(path, "action.json")); err == nil {
		descriptorPath = filepath.Join(path, "action.json")
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
