package util

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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
	description string   `json:"description"`
	Contribs    []string `json:"contributions"`
}

func (d *FlogoContribDescriptor) GetContribType() string {
	return strings.Split(d.Type, ":")[1]
}

func GetContribDescriptor(path string) (*FlogoContribDescriptor, error) {

	files, err := ioutil.ReadDir(path)

	if err != nil {
		fmt.Errorf("Unable to find %v", path)
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
func GetImports(appJsonPath string) ([]string, error) {

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

	for key, _ := range importSet {
		allImports = append(allImports, key)
	}

	return allImports, nil
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

	for key, _ := range importSet {
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

func GetAllImports(path string) ([]string, error) {

	var results []string
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	data := string(bytes)

	pkgs := strings.Split(data[strings.Index(data, "(")+1:], "\n") //Get indivual rows containing pkgs.

	for _, pkg := range pkgs {

		// Remove last line containing ")" and any empty rows
		if !strings.Contains(pkg, ")") && len(pkg) != 0 {
			result := strings.Replace(strings.TrimSpace(pkg), "\"", "", -1)
			results = append(results, strings.TrimSpace(strings.Replace(result, "_", "", -1)))

		}
	}

	return results, nil
}
func ParseImportPath(path string) (string, string) {

	// If @ is speicified split
	if strings.Contains(path, "@") {

		results := strings.Split(path, "@")

		return results[0], results[1]

	}
	return path, ""

}
