// +build ignore

package main

import (
	"github.com/project-flogo/cli/util"
	"os"
)

// This Go program is aimed at being called by go:generate from "cmd/flogo/main.go" to create a "currentversion.go" file
// in the "cmd/flogo" subdirectory.
//
// Once this file is created, it will set at runtime the version of the CLI without relying on GOPATH nor Git command to
// parse the version from the source. Hence it is possible to distribute the CLI as a fully static binary.
//
// Users getting the CLI with a classic "go get" command will still have the version retrieved from the directory
// $GOPATH/src/github.com/project-flogo/cli
func main() {

	_, currentVersion, _ := util.GetCLIInfo()
	wd, err := os.Getwd()
	if err != nil {
		return
	}

	util.CreateVersionFile(wd, currentVersion)
}
