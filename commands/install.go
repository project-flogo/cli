package commands

import (
	"fmt"
	"os"

	"github.com/project-flogo/cli/api"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the version module",
	Long:  `All software has versions. This is Hugo's`,
	Run: func(cmd *cobra.Command, args []string) {
		if api.CheckCurrDir() {
			currDir, _ := os.Getwd()
			api.InstallPackage(os.Args[2], currDir)
		} else {
			fmt.Println("Error in detecting app")
		}

	},
}

func init() {
	RootCmd.AddCommand(installCmd)
}
