package commands

import (
	"fmt"
	"os"
	"go/parser"
	"go/printer"
	"go/token"
	"path/filepath"
	"os/exec"
	"github.com/project-flogo/cli/util"
	"github.com/project-flogo/cli/common"
	"github.com/spf13/cobra"
)
var (
	cliPath =  filepath.Join("src","github.com","project-flogo","cli","cmd","flogo") //"/src/github.com/project-flogo/cli/cmd/flogo"
	path = os.Getenv("GOPATH")
)

var pluginCmd = &cobra.Command{
	Use:              "plugin",
	Short:            "manage your plugins ",
	Long:             `manage your plugins`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("You can install, list and update your plugins for flogo Cli.")
	},
}

var pluginInstall = &cobra.Command{
	Use:   "install",
	Short: "install the plugins to cli ",
	Long:  `install the plugins to cli `,
	Run: func(cmd *cobra.Command, args []string) {
		if len(os.Args) <= 3 {
			fmt.Println("Enter the package name")
			os.Exit(1)
		}
		err := addPlugin(os.Args[3])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		buildModule(os.Args[3], true)

	},
}

var pluginList = &cobra.Command{
	Use:   "list",
	Short: "lists all the installed plugins",
	Long:  `list all the plugins of cli `,
	Run: func(cmd *cobra.Command, args []string) {
		for _, cmd := range common.GetPlugins() {
			fmt.Println(cmd.Use)
		}
	},
}
var pluginUpdate = &cobra.Command{
	Use:   "update",
	Short: "update all the installed plugins",
	Long:  `update all the installed plugins`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if len(os.Args) == 2 {
			fmt.Println("Enter package name")
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
	
		os.Chdir(filepath.Join(path,cliPath))
		
		cliCmd, err := exec.Command("go", "get", "-u",os.Args[2]).CombinedOutput()
		if err != nil {

			fmt.Println(string(cliCmd))

			os.Exit(1)

		}
		buildModule(os.Args[2], true)
		
	},
}

func init() {
	rootCmd.AddCommand(pluginCmd)
	pluginCmd.AddCommand(pluginInstall)
	pluginCmd.AddCommand(pluginList)
	pluginCmd.AddCommand(pluginUpdate)
}



func addPlugin(args string) error {

	
	err := os.Chdir(filepath.Join(path, cliPath))

	currdir, _ := os.Getwd()

	AddModToImportPlugin(args, currdir) //Edit the imports.go file

	//Download all the modules. This is just to ensure all packages are downloaded before go build.
	cliCmd, err := exec.Command("go", "mod", "download").CombinedOutput()
	if err != nil {
		RemoveModFromImportPlugin(args, currdir)

		fmt.Println(string(cliCmd))

		return err
	}

	return nil
}

func buildModule(args string, flag bool) error {


	err := os.Chdir(filepath.Join(path, cliPath))

	if err!=nil{
		return err
	}

	currdir, _ := os.Getwd()

	//Build the modules.
	cliCmd, err := exec.Command("go", "build").CombinedOutput()

	if err != nil {
		RemoveModFromImportPlugin(args, currdir)

		fmt.Println(string(cliCmd))

		return err
	}
	cliCmd, err = exec.Command("cp", filepath.Join(path, cliPath), filepath.Join(path,"bin")).CombinedOutput()
	if err != nil {
		RemoveModFromImportPlugin(args, currdir)

		fmt.Println(string(cliCmd))

		return err
	}
	if flag {
		//Done.
		fmt.Println("Module Successfully Installed")
	}
	return nil

}

func AddModToImportPlugin(pkg string, fpath string) error {

	importsFile := filepath.Join(fpath, "imports.go")

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, importsFile, nil, parser.ImportsOnly)
	if err != nil {
		return err
	}

	util.AddImport(fset, file, pkg)

	f, err := os.Create(importsFile)
	defer f.Close()
	if err := printer.Fprint(f, fset, file); err != nil {
		return err
	}

	return nil
}

func RemoveModFromImportPlugin(pkg string, fpath string) error {
	importsFile := filepath.Join(fpath, "imports.go")

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, importsFile, nil, parser.ImportsOnly)
	if err != nil {
		return err
	}

	util.DeleteImport(fset, file, pkg)

	f, err := os.Create(importsFile)
	defer f.Close()
	if err := printer.Fprint(f, fset, file); err != nil {
		return err
	}

	return nil
}
