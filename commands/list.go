package commands

import (
	"fmt"
	"os"

	"github.com/project-flogo/cli/api"
	"github.com/project-flogo/cli/common"
	"github.com/spf13/cobra"
)

var json bool
var orphaned bool
var listFilter string

func init() {
	listCmd.Flags().BoolVarP(&json, "json", "j", true, "print in json format")
	listCmd.Flags().BoolVarP(&orphaned, "orphaned", "", false, "list orphaned refs")
	listCmd.Flags().StringVarP(&listFilter, "filter", "", "", "apply list filter [used, unused]")
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list [flags]",
	Short: "list installed flogo contributions",
	Long:  "List installed flogo contributions",
	Run: func(cmd *cobra.Command, args []string) {

		if orphaned {
			err := api.ListOrphanedRefs(common.CurrentProject(), json)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting orphaned refs: %v\n", err)
				os.Exit(1)
			}

			return
		}

		err := api.ListContribs(common.CurrentProject(), json, listFilter)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting list of contributions: %v\n", err)
			os.Exit(1)
		}
	},
}
