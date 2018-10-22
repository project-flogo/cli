package commands

import (
	"fmt"
	"os"

	"github.com/project-flogo/cli/registry"

	"github.com/spf13/cobra"
)

//Root command
var RootCmd = &cobra.Command{
	Use:   "Flogo Cli",
	Short: "Flogo Cli lets you work with Flogo",
	Long:  `Flogo Cli is great! `,
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("Welcome to Flogo")

		//fmt.Println("Hugo Static Site Generator v0.9 -- HEAD")
	},
}

func Initialize() {
	//Add the current main commands
	RootCmd.AddCommand(versionCmd)

	//Get the list of commands from the registry of commands and add.
	commandList := registry.GetCommands()

	for _, command := range commandList {

		RootCmd.AddCommand(command)
	}
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Flogo Cli",
	Long:  `Flogo Version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Flogo CLi version v0.0.1")
	},
}

func Execute() {

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
