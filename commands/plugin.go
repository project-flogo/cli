package commands

import (
	"fmt"
	"os"

	"github.com/project-flogo/cli/api"
	"github.com/project-flogo/cli/common"
	"github.com/spf13/cobra"
)

var pluginCmd = &cobra.Command{
	Use:              "plugin",
	Short:            "Flogo Cli lets you explore plugin ",
	Long:             `Flogo Cli create is great! `,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("This is root Plugin Command")
	},
}

var pluginInstall = &cobra.Command{
	Use:   "install",
	Short: "Flogo Cli lets you install plugin ",
	Long:  `Flogo Cli create is great! `,
	Run: func(cmd *cobra.Command, args []string) {
		if len(os.Args) <= 3 {
			fmt.Println("Enter the package name")
			os.Exit(1)
		}
		api.InstallPluginHelper(os.Args[3])

		api.BuildModule(os.Args[3], true)

	},
}

var pluginList = &cobra.Command{
	Use:   "list",
	Short: "Flogo Cli lets you lists all the plugins installed ",
	Long:  `Flogo Cli create is great! `,
	Run: func(cmd *cobra.Command, args []string) {
		for _, cmd := range common.GetPlugins() {
			fmt.Println(cmd.Use)
		}
	},
}

func init() {
	rootCmd.AddCommand(pluginCmd)
	pluginCmd.AddCommand(pluginInstall)
	pluginCmd.AddCommand(pluginList)
}
