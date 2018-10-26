package commands

import (
	"fmt"
	"os"

	"github.com/project-flogo/cli/api"
	"github.com/project-flogo/cli/common"
	"github.com/spf13/cobra"
)

const (
	Version    = "0.0.1"
	VersionTpl = `{{with .Name}}{{printf "%s " .}}{{end}}{{printf "cli version %s" .Version}}
`
)

var verbose bool

//Root command
var rootCmd = &cobra.Command{
	Use:     "flogo [flags] [command]",
	Short:   "flogo cli",
	Long:    `flogo command line interface for flogo applications`,
	Version: Version,

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		api.SetVerbose(verbose)
		common.SetVerbose(verbose)

		builtIn := cmd.Name() == "help" || cmd.Name() == "version"

		if len(os.Args) > 1 && !builtIn {
			currentDir, err := os.Getwd()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: unable to determine working directory - %s\n", err)
				os.Exit(1)
			}
			appProject := api.NewAppProject(currentDir)

			err = appProject.Validate()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			common.SetCurrentProject(appProject)
		}
	},
}

func Initialize() {
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "verbose output")

	rootCmd.SetVersionTemplate(VersionTpl)

	//Get the list of commands from the registry of commands and add.
	commandList := common.GetPlugins()

	for _, command := range commandList {

		rootCmd.AddCommand(command)
	}
}

func Execute() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
