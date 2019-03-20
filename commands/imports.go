package commands

import (
	"fmt"
	"os"

	"github.com/project-flogo/cli/api"
	"github.com/project-flogo/cli/common"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(importsCmd)
	importsCmd.AddCommand(syncCmd)
	importsCmd.AddCommand(resolveCmd)
}

var importsCmd = &cobra.Command{
	Use:   "imports",
	Short: "Manage Imports in the project",
	Long:  `Manage Imports in the project`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}
var syncCmd = &cobra.Command{
	Use:   "sync",
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
var resolveCmd = &cobra.Command{
	Use:   "resolve",
	Short: "resolve a project package",
	Long:  `Resolve the packages in the project`,
	Run: func(cmd *cobra.Command, args []string) {

		err := api.ResolvePkg(common.CurrentProject())

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}
