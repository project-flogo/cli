package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/project-flogo/cli/api"
	"github.com/project-flogo/cli/common"
	"github.com/spf13/cobra"
)

var buildShim string
var buildOptimize bool
var buildEmbed bool
var flogoJsonFile string

func init() {
	buildCmd.Flags().StringVarP(&buildShim, "shim", "", "", "use shim trigger")
	buildCmd.Flags().BoolVarP(&buildOptimize, "optimize", "o", false, "optimize build")
	buildCmd.Flags().BoolVarP(&buildEmbed, "embed", "e", false, "embed configuration in binary")
	buildCmd.Flags().StringVarP(&flogoJsonFile, "file", "f", "", "specify a flogo.json to build")
	rootCmd.AddCommand(buildCmd)
}

//Build the project.
var buildCmd = &cobra.Command{
	Use:              "build",
	Short:            "build the flogo application",
	Long:             `Build the flogo application.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {},
	Run: func(cmd *cobra.Command, args []string) {

		if flogoJsonFile == "" {
			preRun(cmd, args, verbose)
			options := api.BuildOptions{Shim: buildShim, OptimizeImports: buildOptimize, EmbedConfig: buildEmbed}

			err := api.BuildProject(common.CurrentProject(), options)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error building project: %v\n", err)
				os.Exit(1)
			}
		} else {
			//If a jsonFile is specified in the build.
			//Create a new project in the temp folder and copy the bin.

			tempDir, err := api.GetTempDir()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting temp dir: %v\n", err)
				os.Exit(1)
			}

			api.SetVerbose(verbose)
			tempProject, err := api.CreateProject(tempDir, "", flogoJsonFile, "latest")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating temp project: %v\n", err)
				os.Exit(1)
			}

			common.SetCurrentProject(tempProject)

			options := api.BuildOptions{Shim: buildShim, OptimizeImports: buildOptimize, EmbedConfig: buildEmbed}

			err = api.BuildProject(common.CurrentProject(), options)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error building temp project: %v\n", err)
				os.Exit(1)
			}

			copyBin(verbose, tempProject)
		}
	},
}

func copyBin(verbose bool, tempProject common.AppProject) {

	currDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error determining working directory: %v\n", err)
		os.Exit(1)
	}

	if verbose {
		fmt.Printf("Copying the binary from  %s to %s \n", tempProject.BinDir(), currDir)
	}

	if runtime.GOOS == "windows" || api.GOOSENV == "windows" {
		err = os.Rename(tempProject.Executable(), filepath.Join(currDir, "main.exe"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error renaming executable: %v\n", err)
			os.Exit(1)
		}
	} else {
		err = os.Rename(tempProject.Executable(), filepath.Join(currDir, tempProject.Name()))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error renaming executable: %v\n", err)
			os.Exit(1)
		}
	}

	if verbose {
		fmt.Printf("Removing the temp dir: %s\n ", tempProject.Dir())
	}

	err = os.RemoveAll(tempProject.Dir())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error removing temp dir: %v\n", err)
		os.Exit(1)
	}
}
