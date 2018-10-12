package commands

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Flogo Cli also lets you update your packages ",
	Long:  `Flogo Cli update is great! `,
	Run: func(cmd *cobra.Command, args []string) {

		if len(os.Args) == 2 {
			fmt.Println("Enter package name")
			os.Exit(1)
		} else {
			path := os.Getenv("GOPATH")
			os.Chdir(Concat(path, "/src/github.com/project-flogo/cli"))
			cliCmd, err := exec.Command("go", "get", os.Args[2]).CombinedOutput()
			if err != nil {

				fmt.Println(string(cliCmd))

				log.Fatal(err)

			}
			BuildModule(os.Args[2], true)
		}

	},
}
