package commands

import (
	"fmt"
	"os"

	"github.com/project-flogo/cli/api"
	"github.com/project-flogo/cli/common"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:   "update [flags] <contribution|dependency>",
	Short: "update a project contribution/dependency",
	Long:  `Updates a contribution or dependency in the project`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		err := api.UpdatePkg(common.CurrentProject(), args[0])

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error updating contribution/dependency: %v\n", err)
			os.Exit(1)
		}
	},
}
