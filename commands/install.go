package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/project-flogo/cli/api"
	"github.com/project-flogo/cli/common"
	"github.com/spf13/cobra"
)

var replaceContrib string
var contribBundleFile string

func init() {
	installCmd.Flags().StringVarP(&replaceContrib, "replace", "r", "", "specify path to replacement contribution/dependency")
	installCmd.Flags().StringVarP(&contribBundleFile, "file", "f", "", "specify contribution bundle")
	rootCmd.AddCommand(installCmd)
}

var installCmd = &cobra.Command{
	Use:   "install [flags] <contribution|dependency>",
	Short: "install a flogo contribution/dependency",
	Long:  "Installs a flogo contribution or dependency",
	Run: func(cmd *cobra.Command, args []string) {

		if contribBundleFile != "" {
			err := api.InstallContribBundle(common.CurrentProject(), contribBundleFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error installing contribution bundle: %v\n", err)
				os.Exit(1)
			}
		}

		if replaceContrib != "" {
			replaceContrib = strings.Replace(replaceContrib, "@", " ", -1)
			err := api.InstallReplacedPackage(common.CurrentProject(), replaceContrib, args[0])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error installing contribution/dependency: %v\n", err)
				os.Exit(1)
			}
		} else {
			for _, pkg := range args {
				err := api.InstallPackage(common.CurrentProject(), pkg)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error installing contribution/dependency: %v\n", err)
					os.Exit(1)
				}
			}
		}
	},
}
