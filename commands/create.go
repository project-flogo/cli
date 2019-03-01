package commands

import (
	"fmt"
	"os"

	"github.com/project-flogo/cli/api"
	"github.com/spf13/cobra"
)

var (
	flogoJSONPath string
	coreVersion   string
	forceCreate   bool
)

var createCmd = &cobra.Command{
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

		_, err = api.CreateProject(currentDir, appName, flogoJSONPath, coreVersion, forceCreate)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	createCmd.Flags().StringVarP(&flogoJSONPath, "file", "f", "", "path to flogo.json file")
	createCmd.Flags().StringVarP(&coreVersion, "core", "c", "", "specify core library version")
	createCmd.Flags().BoolVar(&forceCreate, "force", false, "force install when go get fails")
	rootCmd.AddCommand(createCmd)
}
