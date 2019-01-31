package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/project-flogo/cli/api"
	"github.com/project-flogo/cli/common"
	"github.com/spf13/cobra"
)

var buildShim string

var buildOptimize bool
var buildEmbed bool
var jsonFile string

func init() {
	buildCmd.Flags().StringVarP(&buildShim, "shim", "", "", "trigger shim")
	buildCmd.Flags().BoolVarP(&buildOptimize, "optimize", "o", false, "optimize build")
	buildCmd.Flags().BoolVarP(&buildEmbed, "embed", "e", false, "embed config")
	buildCmd.Flags().StringVarP(&jsonFile, "file", "f", "", "json file")
	rootCmd.AddCommand(buildCmd)
}

//Build the project.
var buildCmd = &cobra.Command{
	Use:              "build",
	Short:            "build the flogo application",
	Long:             `Build the flogo application.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {},
	Run: func(cmd *cobra.Command, args []string) {

		if jsonFile == "" {
			preRun(cmd, args, verbose)
			options := api.BuildOptions{Shim: buildShim, OptimizeImports: buildOptimize, EmbedConfig: buildEmbed}

			err := api.BuildProject(common.CurrentProject(), options)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

		} else {
			//If a jsonFile is specified in the build.
			//Create a new project in the temp folder and copy the bin.
			currDir, err := os.Getwd()

			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			tempDir, err := api.GetTempDir()

			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			api.SetVerbose(verbose)
			tempProject, err := api.CreateProject(tempDir, "", jsonFile, "master")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			common.SetCurrentProject(tempProject)

			options := api.BuildOptions{Shim: buildShim, OptimizeImports: buildOptimize, EmbedConfig: buildEmbed}

			err = api.BuildProject(common.CurrentProject(), options)

			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			if verbose {
				fmt.Printf("Copying the binary from  %s to %s \n", filepath.Join(tempDir, currProject.Name(), "bin"), currDir)
			}
			_, err = exec.Command("cp", filepath.Join(tempDir, currProject.Name(), "bin", currProject.Name()), currDir).Output()

			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			if verbose {
				fmt.Printf("Removing the temp dir  %s  \n ", tempDir)
			}
			err = os.RemoveAll(tempDir)

			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		}

	},
}
