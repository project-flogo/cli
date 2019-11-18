package main

import (
	"fmt"
	"os"

	"github.com/project-flogo/cli/commands"
	"github.com/project-flogo/cli/util"
)

// Not set by default, will be filled by init() function in "./currentversion.go" file, if it exists.
// This latter file is generated with a "go generate" command.
var Version string = ""

//go:generate go run gen/version.go
func main() {

	if util.GetGoPath() == "" {
		_, _ = fmt.Fprintf(os.Stderr, "Error: GOPATH must be set before running flogo cli\n")
		os.Exit(1)
	}

	//Initialize the commands
	_ = os.Setenv("GO111MODULE", "on")
	commands.Initialize(Version)
	commands.Execute()
}
