package api

import (
	"github.com/project-flogo/cli/util"
)

const (
	fileDescriptorJson string = "descriptor.json"
)

var verbose = false
var scaffold = false

func SetVerbose(enable bool) {
	verbose = enable
	util.SetVerbose(enable)
}

func SetScaffold(val bool) {
	scaffold = val
}

func Verbose() bool {
	return verbose
}

func Scaffold() bool {
	return scaffold
}

//TODO use a logger like struct for API that can be used to log or console output
