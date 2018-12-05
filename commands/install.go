package commands

import (
	"fmt"
	"os"

	"github.com/project-flogo/cli/api"
	"github.com/project-flogo/cli/common"
	"github.com/spf13/cobra"
)

var localContrib string
var installCmd = &cobra.Command{
	Use:   "install [flags] <contribution>",
	Short: "install a flogo contribution",
	Long:  "Installs a flogo contribution",
	Run: func(cmd *cobra.Command, args []string) {

		if localContrib != "" {
			err := api.InstallLocalPackage(common.CurrentProject(), localContrib, args[0])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		} else {
			for _, pkg := range args {
				err := api.InstallPackage(common.CurrentProject(), pkg)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error: %v\n", err)
					os.Exit(1)
				}
			}
		}

	},
}

func init() {
	installCmd.Flags().StringVarP(&localContrib, "localContrib", "l", "", "Specify local Contrib")
	rootCmd.AddCommand(installCmd)

}
