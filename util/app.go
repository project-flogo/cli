package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// PartialAppDescriptor is the descriptor for a Flogo application
type PartialAppDescriptor struct {
	AppModel  string        `json:"appModel"`
	Imports   []string      `json:"imports"`
	Triggers  []interface{} `json:"triggers"`
	Resources []interface{} `json:"resources"`
	Actions   []interface{} `json:"actions"`
}

type void struct{}

type AppImportDetails struct {
	Imp          Import
	ContribDesc  *FlogoContribDescriptor
	TopLevel     bool // a toplevel import i.e. from imports section
	HasAliasRef  bool // imports alias is used by a contrib reference
	HasDirectRef bool // a direct reference exists for this import
}

func (d *AppImportDetails) Referenced() bool {
	return d.HasAliasRef || d.HasDirectRef
}

func (d *AppImportDetails) IsCoreContrib() bool {

	if d.ContribDesc == nil {
		return false
	}
	ct := d.ContribDesc.GetContribType()

	switch ct {
	case "action", "trigger", "activity":
		return true
	default:
		return false
	}
}

type AppImports struct {
	imports     map[string]*AppImportDetails
	orphanedRef map[string]void

	resolveContribs bool
	depManager      DepManager
}

func (ai *AppImports) addImports(imports []string) error {
	for _, anImport := range imports {
		flogoImport, err := ParseImport(anImport)
		if err != nil {
			return err
		}

		if _, exists := ai.imports[flogoImport.GoImportPath()]; exists {
			//todo warn about duplicate import?
			continue
		}

		details, err := ai.newImportDetails(flogoImport)
		if err != nil {
			return err
		}
		details.TopLevel = true

		ai.imports[flogoImport.GoImportPath()] = details
	}

	return nil
}

func (ai *AppImports) addReference(ref string, contribType string) error {
	cleanedRef := strings.TrimSpace(ref)

	if cleanedRef[0] == '#' {
		if !ai.resolveContribs {
			// wont be able to determine contribTypes for existing imports, so just return
			return nil
		}

		alias := cleanedRef[1:]
		found := false
		for _, importDetails := range ai.imports {

			if importDetails.TopLevel {
				// alias refs can only be to toplevel imports
				if importDetails.Imp.CanonicalAlias() == alias && importDetails.ContribDesc != nil &&
					importDetails.ContribDesc.GetContribType() == contribType {
					importDetails.HasAliasRef = true
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
			return nil
		}

		//doesn't exists so add new import
		details, err := ai.newImportDetails(flogoImport)
		if err != nil {
			return err
		}

		ai.imports[flogoImport.GoImportPath()] = details
	}

	return nil
}

func (ai *AppImports) newImportDetails(anImport Import) (*AppImportDetails, error) {
	details := &AppImportDetails{Imp: anImport}

	if ai.resolveContribs {
		desc, err := GetContribDescriptorFromImport(ai.depManager, anImport)
		if err != nil {
			return nil, err
		}
		details.ContribDesc = desc
	}

	return details, nil
}

func (ai *AppImports) GetOrphanedReferences() []string {
	var refs []string
	for ref := range ai.orphanedRef {
		refs = append(refs, ref)
	}

	return refs
}

func (ai *AppImports) GetAllImports() []Import {
	var allImports []Import
	for _, details := range ai.imports {
		allImports = append(allImports, details.Imp)
	}

	return allImports
}

func (ai *AppImports) GetAllImportDetails() []*AppImportDetails {

	var allImports []*AppImportDetails
	for _, details := range ai.imports {
		allImports = append(allImports, details)
	}

	return allImports
}

func GetAppImports(appJsonFile string, depManager DepManager, resolveContribs bool) (*AppImports, error) {
	appJson, err := os.Open(appJsonFile)
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(appJson)
	if err != nil {
		return nil, err
	}

	appDesc := &PartialAppDescriptor{}
	err = json.Unmarshal(bytes, appDesc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to unmarshal flogo app json: %s", appJsonFile)
		return nil, err
	}

	ai := &AppImports{depManager: depManager, resolveContribs: resolveContribs}
	ai.imports = make(map[string]*AppImportDetails)
	ai.orphanedRef = make(map[string]void)

	err = ai.addImports(appDesc.Imports)
	if err != nil {
		return nil, err
	}

	err = extractAppReferences(ai, appDesc)

	return ai, err
}

func extractAppReferences(ai *AppImports, appDesc *PartialAppDescriptor) error {

	//triggers
	for _, trg := range appDesc.Triggers {
		if trgMap, ok := trg.(map[string]interface{}); ok {

			// a ref should exists for every trigger
			if refVal, ok := trgMap["ref"]; ok {
				if strVal, ok := refVal.(string); ok {
					err := ai.addReference(strVal, "trigger")
					if err != nil {
						return err
					}
				}
			}

			// actions are under handlers, so assume an action contribType
			err := extractReferences(ai, trgMap["handlers"], "action")
			if err != nil {
				return err
			}
		}
	}

	//in actions section, refs should be to actions
	err := extractReferences(ai, appDesc.Actions, "action") //action
	if err != nil {
		return err
	}

	//in resources section, refs should be to activities
	err = extractReferences(ai, appDesc.Resources, "activity") //activity
	if err != nil {
		return err
	}

	return nil
}

func extractReferences(ai *AppImports, item interface{}, contribType string) error {
	switch t := item.(type) {
	case map[string]interface{}:
		for key, val := range t {
			if strVal, ok := val.(string); ok {
				if key == "ref" {
					err := ai.addReference(strVal, contribType)
					if err != nil {
						return err
					}
				}
			} else {
				err := extractReferences(ai, val, contribType)
				if err != nil {
					return err
				}
			}
		}
	case []interface{}:
		for _, val := range t {
			err := extractReferences(ai, val, contribType)
			if err != nil {
				return err
			}
		}
	default:
	}

	return nil
}
