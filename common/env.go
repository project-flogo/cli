package common

import (
	"github.com/project-flogo/cli/util"
)

var verbose = false
var scaffold = false
var appProject AppProject

func SetVerbose(enable bool) {
	verbose = enable
	util.SetVerbose(enable)
}

func Verbose() bool {
	return verbose
}

func CurrentProject() AppProject {
	return appProject
}

func SetCurrentProject(project AppProject) {
	appProject = project
}
