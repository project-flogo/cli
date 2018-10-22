package commands

import (
	"github.com/project-flogo/cli/api"
	"github.com/spf13/cobra"
)

//Build the project.
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the App module",
	Long:  `All software has versions. This is Hugo's`,
	Run: func(cmd *cobra.Command, args []string) {
		api.BuildProject()
	},
}

func init() {
	RootCmd.AddCommand(buildCmd)
}
