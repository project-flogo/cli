package commands

import (
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
)

var (
	goPath     = os.Getenv("GOPATH")
	cliPath    = filepath.Join(goPath, filepath.Join("src", "github.com", "project-flogo", "cli"))
	cliCmdPath = filepath.Join(cliPath, "cmd", "flogo")
)

var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "manage cli plugins",
	Long:  "Manage cli plugins",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		common.SetVerbose(verbose)
	},
}

var pluginInstall = &cobra.Command{
	Use:   "install <plugin>",
	Short: "install plugin",
	Long:  "Installs a cli plugin",
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
		added, err := addPlugin(pluginPkg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if added {
			err = updateCLI()
			if err != nil {
				//remove plugin import on failure
				modifyPluginImports(pluginPkg, true)
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Installed plugin\n")
		} else {
			fmt.Printf("Plugin '%s' already installed\n", pluginPkg)
		}
	},
}

var pluginList = &cobra.Command{
	Use:   "list",
	Short: "list installed plugins",
	Long:  "Lists installed cli plugins",
	Run: func(cmd *cobra.Command, args []string) {
		for _, cmd := range common.GetPlugins() {
			fmt.Println(cmd.Name())
		}
	},
}
var pluginUpdate = &cobra.Command{
	Use:   "update <plugin>",
	Short: "update plugin",
	Long:  "Updates the specified installed cli plugin",
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
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		err = updateCLI()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Updated plugin\n")
	},
}

func init() {
	rootCmd.AddCommand(pluginCmd)
	pluginCmd.AddCommand(pluginInstall)
	pluginCmd.AddCommand(pluginList)
	pluginCmd.AddCommand(pluginUpdate)
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

	os.Rename(baseGoMod, bakGoMod)
	os.Rename(buildGoMod, baseGoMod)

	return nil
}

func restoreGoMod() error {

	if verbose {
		fmt.Printf("Restoring default cli go.mod")
	}
	baseGoMod := filepath.Join(cliPath, "go.mod")
	bakGoMod := filepath.Join(cliPath, "go.mod.bak")
	buildGoMod := filepath.Join(cliPath, "go.mod.build")

	os.Rename(baseGoMod, buildGoMod)
	os.Rename(bakGoMod, baseGoMod)

	return nil
}

func addPlugin(pluginPkg string) (bool, error) {

	added, err := modifyPluginImports(pluginPkg, false)
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
		os.Rename(exe, backupExe)
	}

	err = util.ExecCmd(exec.Command("go", "build"), cliCmdPath)
	if err != nil {
		os.Rename(backupExe, exe)
		return err
	}

	os.Rename(filepath.Join(cliCmdPath, "flogo"), exe)
	os.Remove(backupExe)

	return nil
}

func modifyPluginImports(pkg string, remove bool) (bool, error) {

	importsFile := filepath.Join(cliCmdPath, fileImportsGo)

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, importsFile, nil, parser.ImportsOnly)
	if err != nil {
		return false, err
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
