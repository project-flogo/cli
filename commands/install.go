package commands

import (
	"github.com/project-flogo/cli/api"
	"github.com/project-flogo/cli/common"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install [flags] <contribution>",
	Short: "install a flogo contribution",
	Long:  "Installs a flogo contribution",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		api.InstallPackage(common.CurrentProject(), args[0])
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
