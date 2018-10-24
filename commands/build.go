package commands

import (
	"fmt"
	"github.com/project-flogo/cli/common"
	"os"

	"github.com/project-flogo/cli/api"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(buildCmd)
}

//Build the project.
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "build the flogo application",
	Long:  `Build the flogo application.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := api.BuildProject(common.CurrentProject())

		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	},
}
