package commands

import (
	"fmt"
	"os"

	"github.com/project-flogo/cli/api"
	"github.com/spf13/cobra"
)

var flogoJsonPath string
var coreVersion string

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
			fmt.Fprintf(os.Stderr, "Error: unable to determine working directory - %v\n", err)
			os.Exit(1)
		}

		_, err = api.CreateProject(currentDir, appName, flogoJsonPath, coreVersion)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	CreateCmd.Flags().StringVarP(&flogoJsonPath, "file", "f", "", "path to flogo.json file")
	CreateCmd.Flags().StringVarP(&coreVersion, "cv", "", "", "specify core library version")
	rootCmd.AddCommand(CreateCmd)
}
