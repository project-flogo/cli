package main

import (
	"fmt"
	"os"

	"github.com/project-flogo/cli/commands"
)

func main() {
	fmt.Println("Cli App")

	//Initialize the commands
	os.Setenv("GO111MODULE", "on")
	commands.Initialize()
	commands.Execute()
}
