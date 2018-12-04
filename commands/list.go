package commands

import (
	"fmt"
	"os"

	"github.com/project-flogo/cli/api"
	"github.com/project-flogo/cli/common"
	"github.com/spf13/cobra"
)

var json bool
var listCmd = &cobra.Command{
	Use:   "list [flags] <contribution>",
	Short: "list all flogo contribution",
	Long:  "lists a flogo contribution",
	Run: func(cmd *cobra.Command, args []string) {

		err := api.ListPackages(common.CurrentProject(), json)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		os.Exit(1)
	},
}

func init() {
	listCmd.Flags().BoolVarP(&json, "json", "j", true, "print in json format")
	rootCmd.AddCommand(listCmd)
}
