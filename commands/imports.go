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
	importsCmd.AddCommand(importsSyncCmd)
	importsCmd.AddCommand(importsResolveCmd)
	importsCmd.AddCommand(importsListCmd)
}

var importsCmd = &cobra.Command{
	Use:   "imports",
	Short: "manage project imports",
	Long:  `Manage project imports of contributions and dependencies.`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var importsSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "sync Go imports to project imports",
	Long:  `Synchronize Go imports to project imports.`,
	Run: func(cmd *cobra.Command, args []string) {

		err := api.SyncProjectImports(common.CurrentProject())

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error synchronzing imports: %v\n", err)
			os.Exit(1)
		}
	},
}

var importsResolveCmd = &cobra.Command{
	Use:   "resolve",
	Short: "resolve project imports to installed version",
	Long:  `Resolves all project imports to current installed version.`,
	Run: func(cmd *cobra.Command, args []string) {

		err := api.ResolveProjectImports(common.CurrentProject())

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error resolving import versions: %v\n", err)
			os.Exit(1)
		}
	},
}

var importsListCmd = &cobra.Command{
	Use:   "list",
	Short: "list project imports",
	Long:  `List all the project imports`,
	Run: func(cmd *cobra.Command, args []string) {

		err := api.ListProjectImports(common.CurrentProject())

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing imports: %v\n", err)
			os.Exit(1)
		}
	},
}
