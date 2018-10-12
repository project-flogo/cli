package commands

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/project-flogo/cli/registry"
	"github.com/spf13/cobra"
)

var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "Flogo Cli lets you explore plugin ",
	Long:  `Flogo Cli create is great! `,
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
		InstallPluginHelper(os.Args[3])

		BuildModule(os.Args[3], true)

	},
}

var pluginList = &cobra.Command{
	Use:   "list",
	Short: "Flogo Cli lets you lists all the plugins installed ",
	Long:  `Flogo Cli create is great! `,
	Run: func(cmd *cobra.Command, args []string) {
		for _, cmd := range registry.GetCommands() {
			fmt.Println(cmd.Use)
		}
	},
}

func init() {
	RootCmd.AddCommand(pluginCmd)
	pluginCmd.AddCommand(pluginInstall)
	pluginCmd.AddCommand(pluginList)
}

func InstallPluginHelper(args string) error {

	//Get the current GOPATH.
	path := os.Getenv("GOPATH")
	//Change the Dir.
	err := os.Chdir(Concat(path, "/src/cli/cmd/flogo"))

	die(err)

	currdir, _ := os.Getwd()

	AddModToImportPlugin(args, currdir) //Edit the imports.go file

	//Download all the modules. This is just to ensure all packages are downloaded before go build.
	cliCmd, err := exec.Command("go", "mod", "download").CombinedOutput()
	if err != nil {
		RemoveModFromImportPlugin(args, currdir)

		fmt.Println(string(cliCmd))

		log.Fatal(err)
	}
	die(err)

	return err
}

func BuildModule(args string, flag bool) {

	//Get the current GOPATH.
	path := os.Getenv("GOPATH")
	//Change the Dir.
	err := os.Chdir(Concat(path, "/src/cli/cmd/flogo"))

	die(err)

	currdir, _ := os.Getwd()

	//Build the modules.
	cliCmd, err := exec.Command("go", "build").CombinedOutput()

	if err != nil {
		RemoveModFromImportPlugin(args, currdir)

		fmt.Println(string(cliCmd))

		log.Fatal(err)
	}
	cliCmd, err = exec.Command("cp", Concat(os.Getenv("GOPATH"), "/src/cli/cmd/flogo/flogo"), Concat(os.Getenv("GOPATH"), "/bin")).CombinedOutput()
	if err != nil {
		RemoveModFromImportPlugin(args, currdir)

		fmt.Println(string(cliCmd))

		log.Fatal(err)
	}
	if flag {
		//Done.
		fmt.Println("Module Successfully Installed")
	}

}

func AddModToImportPlugin(pkg string, fpath string) {

	byteArray, err := ioutil.ReadFile(Concat(fpath, "/imports.go"))

	die(err)

	text := string(byteArray)

	index := strings.Index(text, ")")

	err = ioutil.WriteFile(Concat(fpath, "/imports.go"), []byte(Concat(text[:index-1], "\n _ \"", pkg, "\" \n", ")")), 0)

	die(err)
}

func RemoveModFromImportPlugin(pkg string, fpath string) {

	byteArray, err := ioutil.ReadFile(Concat(fpath, "/imports.go"))

	die(err)

	text := string(byteArray)

	newText := strings.Replace(text, Concat(" _ \"", pkg, "\" "), "", -1)
	err = ioutil.WriteFile(Concat(fpath, "/imports.go"), []byte(newText), 0)

	die(err)
}
