package main

import (
	"os"

	"github.com/project-flogo/cli/commands"
)

func main() {
	//Initialize the commands
	os.Setenv("GO111MODULE", "on")
	commands.Initialize()
	commands.Execute()
}
