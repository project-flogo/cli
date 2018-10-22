package commands

import (
	"fmt"
	"os"

	"github.com/project-flogo/cli/api"
	"github.com/spf13/cobra"
)

var file bool
var core string

var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Flogo Cli lets you create with Flogo",
	Long:  `Flogo Cli create is great! `,
	Run: func(cmd *cobra.Command, args []string) {
		if file {
			if len(os.Args) <= 3 {
				fmt.Println("Enter file name")
				os.Exit(1)
			}
		}

		currDir, err := os.Getwd()
		if err != nil {
			fmt.Println("Error in determinng Dir")
			os.Exit(1)
		}

		err = api.CreateProject(os.Args[len(os.Args)-1], file, core, currDir)
		fmt.Println(err)
	},
}

func init() {
	RootCmd.AddCommand(CreateCmd)
	CreateCmd.Flags().BoolVarP(&file, "file", "f", false, "Enter file")
	CreateCmd.Flags().StringVarP(&core, "core", "c", "", "Enter core version")
}
