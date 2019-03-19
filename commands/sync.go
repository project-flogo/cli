package commands

import (
	"fmt"
	"os"

	"github.com/project-flogo/cli/api"
	"github.com/project-flogo/cli/common"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(syncCmd)
}

var syncCmd = &cobra.Command{
	Use:   "sync [flags] <package>",
	Short: "sync a project package",
	Long:  `Sync a package in the project`,
	Run: func(cmd *cobra.Command, args []string) {

		err := api.SyncPkg(common.CurrentProject())

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}
