package main

import (
	"fmt"

	"github.com/project-flogo/cli/cmd"
)

func main() {
	fmt.Println("Cli App")

	//Initialize the commands

	cmd.Initialize()
	cmd.Execute()
}
