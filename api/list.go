package api

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/project-flogo/cli/common"
	"github.com/project-flogo/cli/util"
)

type ListFilter int

func ListContribs(project common.AppProject, jsonFormat bool, filter string) error {

	ai, err := util.GetAppImports(filepath.Join(project.Dir(), fileFlogoJson), project.DepManager(), true)
	if err != nil {
		return err
	}

	var specs []*ContribSpec

	for _, details := range ai.GetAllImportDetails() {

		if !includeContrib(details, filter) {
			continue
		}

		specs = append(specs, getContribSpec(project, details))

	}

	for _, details := range ai.GetAllImportDetails() {

		if details.ContribDesc == nil {
			continue
		}

		if details.ContribDesc.Type == "flogo:function" {
			specs = append(specs, getContribSpec(project, details))
		}
	}

	if len(specs) == 0 {
		return nil
	}

	if jsonFormat {
		resp, err := json.MarshalIndent(specs, "", "  ")
		if err != nil {
			return err
		}

		fmt.Fprintf(os.Stdout, "%v \n", string(resp))
	} else {
		for _, spec := range specs {
			fmt.Println("Contrib: " + spec.Name)
			fmt.Println("  Type       : " + spec.Type)
			if spec.IsLegacy != nil {
				fmt.Println("  IsLegacy   : true")
			}
			fmt.Println("  Homepage   : " + spec.Homepage)
			fmt.Println("  Ref        : " + spec.Ref)
			fmt.Println("  Path       : " + spec.Path)
			fmt.Println("  Descriptor : " + spec.Path)
			fmt.Println("  Description: " + spec.Description)
			fmt.Println()
		}
	}

	return nil
}

func includeContrib(details *util.AppImportDetails, filter string) bool {

	if details.IsCoreContrib() {

		switch strings.ToLower(filter) {
		case "used":
			return details.Referenced()
		case "unused":
			return !details.Referenced()
		default:
			return true
		}
	}

	return false

}

type ContribSpec struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Description string      `json:"description"`
	Homepage    string      `json:"homepage"`
	Ref         string      `json:"ref"`
	Path        string      `json:"path"`
	Descriptor  string      `json:"descriptor"`
	IsLegacy    interface{} `json:"isLegacy,omitempty"`
}

func getContribSpec(project common.AppProject, details *util.AppImportDetails) *ContribSpec {
	path, err := project.GetPath(details.Imp)
	if err != nil {
		return nil
	}

	if Verbose() {
		fmt.Println("Path of contrib", path, "for contrib", details.Imp)
	}

	desc := details.ContribDesc

	spec := &ContribSpec{}
	spec.Name = desc.Name
	spec.Type = desc.Type
	spec.Description = desc.Description
	spec.Homepage = desc.Homepage
	spec.Ref = details.Imp.ModulePath()
	spec.Path = path

	if desc.IsLegacy {
		spec.IsLegacy = true
		spec.Descriptor = desc.GetContribType() + ".json"
	} else {
		spec.Descriptor = "descriptor.json"
	}

	return spec
}
func ListOrphanedRefs(project common.AppProject, jsonFormat bool) error {

	ai, err := util.GetAppImports(filepath.Join(project.Dir(), fileFlogoJson), project.DepManager(), true)
	if err != nil {
		return err
	}

	orphaned := ai.GetOrphanedReferences()

	if jsonFormat {
		resp, err := json.MarshalIndent(orphaned, "", "  ")
		if err != nil {
			return err
		}

		fmt.Fprintf(os.Stdout, "%v \n", string(resp))
	} else {
		for _, ref := range orphaned {
			fmt.Println(ref)
		}
	}

	return nil
}
