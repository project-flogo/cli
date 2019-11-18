package commands

import (
	"fmt"
	"os"

	"github.com/project-flogo/cli/api"
	"github.com/spf13/cobra"
)

var flogoJsonPath string
var coreVersion string

func init() {
	CreateCmd.Flags().StringVarP(&flogoJsonPath, "file", "f", "", "specify a flogo.json to create project from")
	CreateCmd.Flags().StringVarP(&coreVersion, "cv", "", "", "specify core library version (ex. master)")
	rootCmd.AddCommand(CreateCmd)
}

var CreateCmd = &cobra.Command{
	Use:              "create [flags] [appName]",
	Short:            "create a flogo application project",
	Long:             `Creates a flogo application project.`,
	Args:             cobra.RangeArgs(0, 1),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {},
	Run: func(cmd *cobra.Command, args []string) {

		api.SetVerbose(verbose)
		appName := ""
		if len(args) > 0 {
			appName = args[0]
		}

		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error determining working directory: %v\n", err)
			os.Exit(1)
		}
		_, err = api.CreateProject(currentDir, appName, flogoJsonPath, coreVersion)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating project: %v\n", err)
			os.Exit(1)
		}
	},
}
