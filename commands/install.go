package commands

import (
	"fmt"
	"os"

	"github.com/project-flogo/cli/api"
	"github.com/project-flogo/cli/common"
	"github.com/spf13/cobra"
)

var (
	localContrib string
	palette      string
	forceInstall bool
)

var installCmd = &cobra.Command{
	Use:   "install [flags] <contribution>",
	Short: "install a flogo contribution",
	Long:  "Installs a flogo contribution",
	Run: func(cmd *cobra.Command, args []string) {
		if palette != "" {
			err := api.InstallPalette(common.CurrentProject(), palette, forceInstall)
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
				err := api.InstallPackage(common.CurrentProject(), pkg, forceInstall)
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
	installCmd.Flags().StringVarP(&palette, "palette", "p", "", "Specify Palette")
	installCmd.Flags().BoolVar(&forceInstall, "force", false, "force install when go get fails")
	rootCmd.AddCommand(installCmd)

}
