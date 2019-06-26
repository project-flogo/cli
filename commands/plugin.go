package commands

import (
	"errors"
	"fmt"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/project-flogo/cli/common"
	"github.com/project-flogo/cli/util"
	"github.com/spf13/cobra"
)

const (
	fileImportsGo = "imports.go"
	add           = false
	remove        = true
)

func init() {
	pluginCmd.AddCommand(pluginInstallCmd)
	pluginCmd.AddCommand(pluginListCmd)
	pluginCmd.AddCommand(pluginUpdateCmd)
	pluginCmd.AddCommand(pluginRemoveCmd)
	rootCmd.AddCommand(pluginCmd)
}

var (
	goPath     = os.Getenv("GOPATH")
	cliPath    = filepath.Join(goPath, filepath.Join("src", "github.com", "project-flogo", "cli"))
	cliCmdPath = filepath.Join(cliPath, "cmd", "flogo")
)

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

		err := useBuildGoMod()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		defer restoreGoMod()

		pluginPkg := args[0]

		fmt.Printf("Installing plugin: %s\n", pluginPkg)

		added, err := updatePlugin(pluginPkg, add)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error adding plugin: %v\n", err)
			os.Exit(1)
		}

		if added {
			err = updateCLI()

			if err != nil {
				fmt.Fprintf(os.Stderr, "Error updating CLI: %v\n", err)
				//remove plugin import on failure

				modifyPluginImports(pluginPkg, true)

				os.Exit(1)
			}

			fmt.Printf("Installed plugin\n")
		} else {
			fmt.Printf("Plugin '%s' already installed\n", pluginPkg)
		}
	},
}

var pluginListCmd = &cobra.Command{
	Use:   "list",
	Short: "list installed plugins",
	Long:  "Lists installed CLI plugins",
	Run: func(cmd *cobra.Command, args []string) {
		for _, cmd := range common.GetPlugins() {
			fmt.Println(cmd.Name())
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
		removed, err := updatePlugin(pluginPkg, remove)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error adding plugin: %v\n", err)
			os.Exit(1)
		}

		if removed {
			err = updateCLI()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error updating CLI: %v\n", err)
				//remove plugin import on failure
				os.Exit(1)
			}
			fmt.Printf("Removed plugin %v \n", pluginPkg)
		}

	},
}

var pluginUpdateCmd = &cobra.Command{
	Use:   "update <plugin>",
	Short: "update plugin",
	Long:  "Updates the specified installed CLI plugin",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		err := useBuildGoMod()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		defer restoreGoMod()

		plugin := args[0]
		fmt.Printf("Updating plugin: %s\n", plugin)

		err = util.ExecCmd(exec.Command("go", "get", "-u", plugin), cliCmdPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error updating plugin: %v\n", err)
			os.Exit(1)
		}

		err = updateCLI()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error updating CLI: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Updated plugin\n")
	},
}

func useBuildGoMod() error {

	baseGoMod := filepath.Join(cliPath, "go.mod")
	bakGoMod := filepath.Join(cliPath, "go.mod.bak")
	buildGoMod := filepath.Join(cliPath, "go.mod.build")

	if _, err := os.Stat(buildGoMod); err != nil {

		if verbose {
			fmt.Printf("Creating plugin build go.mod")
		}

		err := util.CopyFile(baseGoMod, buildGoMod)
		if err != nil {
			return err
		}
	}

	if verbose {
		fmt.Printf("Switching to plugin build go.mod")
	}

	err := os.Rename(baseGoMod, bakGoMod)
	if err != nil {
		return err
	}

	err = os.Rename(buildGoMod, baseGoMod)
	if err != nil {
		return err
	}

	return nil
}

func restoreGoMod() {

	if verbose {
		fmt.Printf("Restoring default CLI go.mod")
	}
	baseGoMod := filepath.Join(cliPath, "go.mod")
	bakGoMod := filepath.Join(cliPath, "go.mod.bak")
	buildGoMod := filepath.Join(cliPath, "go.mod.build")

	err := os.Rename(baseGoMod, buildGoMod)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
	err = os.Rename(bakGoMod, baseGoMod)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
}

func updatePlugin(pluginPkg string, opt bool) (bool, error) {

	err := util.ExecCmd(exec.Command("go", "get", pluginPkg), cliCmdPath)
	if err != nil {
		return false, err
	}

	added, err := modifyPluginImports(pluginPkg, opt)
	if err != nil {
		return added, err
	}

	if added {
		//Download all the modules. This is just to ensure all packages are downloaded before go build.
		err := util.ExecCmd(exec.Command("go", "mod", "download"), cliCmdPath)
		if err != nil {
			modifyPluginImports(pluginPkg, true)
			return false, err
		}
	}

	return added, nil
}

func updateCLI() error {

	exe, err := os.Executable()
	if err != nil {
		return err
	}

	backupExe := exe + ".bak"
	if _, err := os.Stat(exe); err == nil {
		err = os.Rename(exe, backupExe)
		if err != nil {
			return err
		}
	}

	err = util.ExecCmd(exec.Command("go", "build"), cliCmdPath)
	if err != nil {
		osErr := os.Rename(backupExe, exe)
		fmt.Fprintf(os.Stderr, "Error: %v\n", osErr)
		return err
	}

	err = os.Rename(filepath.Join(cliCmdPath, "flogo"), exe)
	if err != nil {
		return err
	}

	err = os.Remove(backupExe)
	if err != nil {
		return err
	}

	return nil
}

func modifyPluginImports(pkg string, remove bool) (bool, error) {

	importsFile := filepath.Join(cliCmdPath, fileImportsGo)

	fset := token.NewFileSet()
	file, _ := parser.ParseFile(fset, importsFile, nil, parser.ImportsOnly)

	if file.Imports == nil {
		return false, errors.New("No Imports found.")
	}

	successful := false

	if remove {

		successful = util.DeleteImport(fset, file, pkg)
	} else {
		successful = util.AddImport(fset, file, pkg)
	}

	if successful {
		f, err := os.Create(importsFile)
		if err != nil {
			return false, err
		}
		defer f.Close()
		if err := printer.Fprint(f, fset, file); err != nil {
			return false, err
		}
	}

	return successful, nil
}
