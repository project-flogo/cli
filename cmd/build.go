package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

//Build the project.
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the App module",
	Long:  `All software has versions. This is Hugo's`,
	Run: func(cmd *cobra.Command, args []string) {
		if checkCurrDir() {
			path, _ := os.Getwd()
			os.Chdir(Concat(path, "/src"))
			cliCmd, err := exec.Command("go", "build").CombinedOutput()
			if err != nil {
				fmt.Println(string(cliCmd))
			}
			die(err)
			_, err = exec.Command("cp", "main", "../bin/").Output()
			die(err)
			os.Chdir(path)
		}

	},
}

func init() {
	RootCmd.AddCommand(buildCmd)
}
