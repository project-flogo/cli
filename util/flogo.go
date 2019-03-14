package util

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	fpath "path"
	"path/filepath"
	"strings"
)

var exists = struct{}{}

// ParseAppDescriptor parse the application descriptor
func ParseAppDescriptor(appJson string) (*FlogoAppDescriptor, error) {
	descriptor := &FlogoAppDescriptor{}

	err := json.Unmarshal([]byte(appJson), descriptor)

	if err != nil {
		return nil, err
	}

	return descriptor, nil
}

// FlogoAppDescriptor is the descriptor for a Flogo application
type FlogoAppDescriptor struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Version     string   `json:"version"`
	Description string   `json:"description"`
	AppModel    string   `json:"appModel,omitempty"`
	Imports     []string `json:"imports"`

	Triggers []*FlogoTriggerConfig `json:"triggers"`
}

type FlogoTriggerConfig struct {
	Id   string `json:"id"`
	Ref  string `json:"ref"`
	Type string `json:"type"`
}

type AppConfig struct {
	Imports   []string      `json:"imports,omitempty"`
	Triggers  []interface{} `json:"triggers"`
	Resources []interface{} `json:"resources,omitempty"`
	Actions   []interface{} `json:"actions,omitempty"`
}

// FlogoAppDescriptor is the descriptor for a Flogo application
type FlogoContribDescriptor struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Homepage    string `json:"homepage"`
	Ref         string `json:"ref"`
	Shim        string `json:"shim"`
}

type FlogoContribBundleDescriptor struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Contribs    []string `json:"contributions"`
}

func (d *FlogoContribDescriptor) GetContribType() string {
	parts := strings.Split(d.Type, ":")
	if len(parts) > 1 {
		return parts[1]
	}
	return ""
}

func GetContribDescriptor(path string) (*FlogoContribDescriptor, error) {

	files, err := ioutil.ReadDir(path)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to find %v", path)
		return nil, err
	}

	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".json") {

			desc, err := ReadContribDescriptor(filepath.Join(path, f.Name()))
			if err == nil {
				return desc, nil
			}
		}
	}
	return nil, nil
}

// ParseAppDescriptor parse the application descriptor
func GetImports(appJsonPath string) (Imports, error) {

	importSet := make(map[string]struct{})

	imports, err := getImports(appJsonPath)
	if err != nil {
		return nil, err
	}

	for _, value := range imports {
		importSet[value] = exists
	}

	if len(imports) == 0 {
		imports, err = getImportsLegacy(appJsonPath)
		if err != nil {
			return nil, err
		}

		for _, value := range imports {
			importSet[value] = exists
		}
	}

	var allImports []string

	for key := range importSet {
		allImports = append(allImports, key)
	}

	var result Imports
	for _, i := range allImports {
		parsedImport, err := ParseImport(i)
		if err != nil {
			return nil, err
		}
		result = append(result, parsedImport)
	}

	return result, nil
}

func getImports(appJsonPath string) ([]string, error) {
	appJsonFile, err := os.Open(appJsonPath)
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(appJsonFile)
	if err != nil {
		return nil, err
	}

	descriptor := &FlogoAppDescriptor{}

	err = json.Unmarshal(bytes, descriptor)
	if err != nil {
		return nil, err
	}

	return descriptor.Imports, nil
}

func getImportsLegacy(appJsonPath string) ([]string, error) {

	importSet := make(map[string]struct{})

	file, err := os.Open(appJsonPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if idx := strings.Index(line, "\"ref\""); idx > -1 {
			startPkgIdx := strings.Index(line[idx+6:], "\"")
			pkg := strings.Split(line[idx+6+startPkgIdx:], "\"")[1]

			importSet[pkg] = exists
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	var imports []string

	for key := range importSet {
		imports = append(imports, key)
	}

	return imports, nil
}

func ReadContribDescriptor(descriptorFile string) (*FlogoContribDescriptor, error) {

	descriptorJson, err := os.Open(descriptorFile)
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(descriptorJson)
	if err != nil {
		return nil, err
	}

	descriptor := &FlogoContribDescriptor{}

	err = json.Unmarshal(bytes, descriptor)
	if err != nil {
		return nil, fmt.Errorf("failed to parse descriptor '%s': %s", descriptorFile, err.Error())
	}

	return descriptor, nil
}

func ParseImportPath(path string) (string, string) {

	// If @ is specified split
	if strings.Contains(path, "@") {

		results := strings.Split(path, "@")

		return results[0], results[1]

	}
	return path, ""
}

func GetImportsFromJSON(path string) (Imports, error) {

	appConfig := &AppConfig{}
	//fmt.Println("Path is", path)
	descriptorJson, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(descriptorJson)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, appConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to marshal ")
		return nil, err
	}

	refs := getRefsFromConfig(appConfig)
	var result Imports

	for _, key := range refs {
		found := false

		for index, contrib := range appConfig.Imports {
			flogoImport, err := ParseImport(contrib)
			if err != nil {
				return nil, err
			}
			if fpath.Base(contrib) == key || flogoImport.Alias() == key {
				found = true

				result = append(result, flogoImport)
				//delete the found contrib. Reduces the number of iteration.
				appConfig.Imports = append(appConfig.Imports[:index], appConfig.Imports[index+1:]...)
			}
		}
		//
		if !found {
			flogoImport, err := ParseImport(key)
			if err != nil {
				return nil, err
			}
			result = append(result, flogoImport)
		}

	}
	return result, nil
}

func getRefsFromConfig(appConfig *AppConfig) []string {
	var results []string
	triggers := extractDependencies(appConfig.Triggers)
	for _, val := range triggers {
		if val != "" {
			results = append(results, val)
		}
	}
	resources := extractDependencies(appConfig.Resources)
	for _, val := range resources {
		if val != "" {
			results = append(results, val)
		}
	}
	actions := extractDependencies(appConfig.Actions)
	for _, val := range actions {
		if val != "" {
			results = append(results, val)
		}
	}

}

func extractDependencies(resource interface{}) []string {
	var refs []string
	switch resource.(type) {
	case map[string]interface{}:

		for key, val := range resource.(map[string]interface{}) {
			//Type is deprecated use ref instead.
			if key == "ref" {
				val = strings.Trim(val.(string), "#")
				refs = append(refs, val.(string))
				return refs
			}
			refs = append(refs, extractDependencies(resource.(map[string]interface{})[key])...)
		}
	case []interface{}:

		for i := 0; i < len(resource.([]interface{})); i++ {
			refs = append(refs, extractDependencies(resource.([]interface{})[i])...)
		}
	default:
		return append(refs, "")
	}
	return refs
}
