package commands

import (
	"fmt"
	"github.com/project-flogo/cli/api"
	"github.com/project-flogo/cli/common"
	"os"

	"github.com/spf13/cobra"
)

var verbose bool

//Root command
var rootCmd = &cobra.Command{
	Use:   "flogo [flags] [command]",
	Short: "flogo cli",
	Long:  `flogo command line interface for flogo applications`,
	//Args: cobra.MinimumNArgs(1),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		api.SetVerbose(verbose)
		common.SetVerbose(verbose)

		builtIn := cmd.Name() == "help" || cmd.Name() == "version" 

		if len(os.Args) > 1 && !builtIn {
			currentDir, err := os.Getwd()
			if err != nil {
				fmt.Fprintf(os.Stderr, "error - unable to determine working directory: %s\n", err)
				os.Exit(1)
			}
			appProject := api.NewAppProject(currentDir)

			err = appProject.Validate()
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}

			common.SetCurrentProject(appProject)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
	},
}

func Initialize() {
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "verbose output")

	//Add the current main commands
	rootCmd.AddCommand(versionCmd)

	//Get the list of commands from the registry of commands and add.
	commandList := common.GetPlugins()

	for _, command := range commandList {

		rootCmd.AddCommand(command)
	}
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "displays the version of flogo cli",
	Long:  `Get the current version number of the cli.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("flogo cli version 0.0.1")
	},
}

func Execute() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
