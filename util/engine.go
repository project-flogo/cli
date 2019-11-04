package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// PartialEngineDescriptor is the descriptor for a Flogo application
type PartialEngineDescriptor struct {
	Imports  []string               `json:"imports"`
	Services []*EngineServiceDetails `json:"services"`
}

type EngineServiceDetails struct {
	Ref string
}

type EngineImportDetails struct {
	Imp          Import
	TopLevel     bool // a toplevel import i.e. from imports section
	ServiceRef  bool // imports is used by a service

	HasAliasRef  bool // imports alias is used by a contrib reference
	HasDirectRef bool // a direct reference exists for this import
}

type EngineImports struct {
	imports     map[string]*EngineImportDetails
	orphanedRef map[string]void
	depManager      DepManager
}

func (ai *EngineImports) addImports(imports []string) error {
	for _, anImport := range imports {
		flogoImport, err := ParseImport(anImport)
		if err != nil {
			return err
		}

		if _, exists := ai.imports[flogoImport.GoImportPath()]; exists {
			//todo warn about duplicate import?
			continue
		}

		details := &EngineImportDetails{Imp: flogoImport, TopLevel:true}
		ai.imports[flogoImport.GoImportPath()] = details
	}

	return nil
}

func (ai *EngineImports) addReference(ref string, isService bool) error {
	cleanedRef := strings.TrimSpace(ref)

	if cleanedRef[0] == '#' {

		alias := cleanedRef[1:]
		found := false
		for _, importDetails := range ai.imports {

			if importDetails.TopLevel {
				// alias refs can only be to toplevel imports
				if importDetails.Imp.CanonicalAlias() == alias {
					importDetails.HasAliasRef = true
					importDetails.ServiceRef = isService
					found = true
					break
				}
			}
		}

		if !found {
			ai.orphanedRef[cleanedRef] = void{}
		}

	} else {
		flogoImport, err := ParseImport(ref)
		if err != nil {
			return err
		}

		if imp, exists := ai.imports[flogoImport.GoImportPath()]; exists {
			if !imp.TopLevel {
				//already accounted for
				return nil
			}

			imp.HasDirectRef = true
			imp.ServiceRef = true
			return nil
		}

		//doesn't exists so add new import
		details  := &EngineImportDetails{Imp: flogoImport, ServiceRef:isService}
		ai.imports[flogoImport.GoImportPath()] = details
	}

	return nil
}

func (ai *EngineImports) GetOrphanedReferences() []string {
	var refs []string
	for ref := range ai.orphanedRef {
		refs = append(refs, ref)
	}

	return refs
}

func (ai *EngineImports) GetAllImports() []Import {
	var allImports []Import
	for _, details := range ai.imports {
		allImports = append(allImports, details.Imp)
	}

	return allImports
}

func (ai *EngineImports) GetAllImportDetails() []*EngineImportDetails {

	var allImports []*EngineImportDetails
	for _, details := range ai.imports {
		allImports = append(allImports, details)
	}

	return allImports
}

func GetEngineImports(engJsonFile string, depManager DepManager) (*EngineImports, error) {
	engJson, err := os.Open(engJsonFile)
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(engJson)
	if err != nil {
		return nil, err
	}

	engDesc := &PartialEngineDescriptor{}
	err = json.Unmarshal(bytes, engDesc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to unmarshal flogo engine json: %s", engJsonFile)
		return nil, err
	}

	imports := make(map[string]void)

	if len(engDesc.Imports) > 0 {
		for _, imp := range engDesc.Imports {
			imports[imp] = void{}
		}
	}

	if len(engDesc.Services) > 0 {
		for _, service := range engDesc.Services {
			imports[service.Ref] = void{}
		}
	}


	ai := &EngineImports{depManager: depManager}
	ai.imports = make(map[string]*EngineImportDetails)
	ai.orphanedRef = make(map[string]void)

	err = ai.addImports(engDesc.Imports)
	if err != nil {
		return nil, err
	}

	// add service refs/imports
	for _, serviceDetails := range engDesc.Services {
		err := ai.addReference(serviceDetails.Ref, true)
		if err != nil {
			return nil, err
		}
	}

	return ai, err
}
