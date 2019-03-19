package commands

import (
	"fmt"
	"os"

	"github.com/project-flogo/cli/api"
	"github.com/project-flogo/cli/common"
	"github.com/spf13/cobra"
)

var localContrib string
var listFilter string

var installCmd = &cobra.Command{
	Use:   "install [flags] <contribution|dependency>",
	Short: "install a flogo contribution/dependency",
	Long:  "Installs a flogo contribution or dependency",
	Run: func(cmd *cobra.Command, args []string) {

		if listFilter != "" {
			err := api.InstallContribBundle(common.CurrentProject(), listFilter)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		}

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
	installCmd.Flags().StringVarP(&localContrib, "local", "l", "", "specify local contribution")
	installCmd.Flags().StringVarP(&listFilter, "file", "f", "", "specify contribution bundle")
	rootCmd.AddCommand(installCmd)
}
