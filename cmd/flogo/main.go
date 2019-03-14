package main

import (
	"os"

	"github.com/project-flogo/cli/commands"
)

// Not set by default, will be filled by init() function in "./currentversion.go" file, if it exists.
// This latter file is generated with a "go generate" command.
var Version string = ""

//go:generate go run gen/version.go
func main() {
	//Initialize the commands
	os.Setenv("GO111MODULE", "on")
	commands.Initialize(Version)
	commands.Execute()
}
