package main

import (
	"fmt"

	"github.com/project-flogo/cli/commands"
)

func main() {
	fmt.Println("Cli App")

	//Initialize the commands

	commands.Initialize()
	commands.Execute()
}
