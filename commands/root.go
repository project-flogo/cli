package commands

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/project-flogo/cli/registry"

	"github.com/spf13/cobra"
)

var legacySupport bool

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

func die(err error) {
	if err != nil {
		fmt.Println("Error in module installtion")
		log.Fatal(err)
	}
}

func Concat(path ...string) string {
	var b bytes.Buffer

	for _, p := range path {
		b.WriteString(p)
	}

	return b.String()
}
