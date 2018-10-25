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
	Use:   "update",
	Short: "update the project packages",
	Long:  `update the project packages`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(os.Args) < 3 {
			fmt.Println("Enter package name")
			os.Exit(1)
		}
		err := api.UpdatePkg(common.CurrentProject(), os.Args[len(os.Args)-1])

		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	},
}
