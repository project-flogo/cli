package commands

import (
	"fmt"
	"github.com/project-flogo/cli/common"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	pluginCmd.AddCommand(pluginInstallCmd)
	pluginCmd.AddCommand(pluginListCmd)
	pluginCmd.AddCommand(pluginUpdateCmd)
	pluginCmd.AddCommand(pluginRemoveCmd)
	rootCmd.AddCommand(pluginCmd)
}

var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "manage CLI plugins",
	Long:  "Manage CLI plugins",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		common.SetVerbose(verbose)
	},
}

var pluginInstallCmd = &cobra.Command{
	Use:   "install <plugin>",
	Short: "install CLI plugin",
	Long:  "Installs a CLI plugin",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		pluginPkg := args[0]

		fmt.Printf("Installing plugin: %s\n", pluginPkg)

		err := UpdateCLI(pluginPkg, UpdateOptAdd)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error adding plugin: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Installed plugin: %s\n", pluginPkg)
	},
}

var pluginListCmd = &cobra.Command{
	Use:   "list",
	Short: "list installed plugins",
	Long:  "Lists installed CLI plugins",
	Run: func(cmd *cobra.Command, args []string) {

		for _, pluginPkg := range common.GetPluginPkgs() {
			fmt.Println(pluginPkg)
		}
	},
}

var pluginRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "remove installed plugins",
	Long:  "Remove installed CLI plugins",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		pluginPkg := args[0]

		fmt.Printf("Removing plugin: %s\n", pluginPkg)

		err := UpdateCLI(pluginPkg, UpdateOptRemove)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error adding plugin: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Removed plugin: %s\n", pluginPkg)
	},
}

var pluginUpdateCmd = &cobra.Command{
	Use:   "update <plugin>",
	Short: "update plugin",
	Long:  "Updates the specified installed CLI plugin",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		pluginPkg := args[0]

		fmt.Printf("Updating plugin: %s\n", pluginPkg)

		err := UpdateCLI(pluginPkg, UpdateOptUpdate)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error updating plugin: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Updated plugin: %s\n", pluginPkg)
	},
}
