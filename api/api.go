package api

import (
	"github.com/project-flogo/cli/util"
)

const (
	fileDescriptorJson string = "descriptor.json"
)

var verbose = false

func SetVerbose(enable bool) {
	verbose = enable
	util.SetVerbose(enable)
}

func Verbose() bool {
	return verbose
}
