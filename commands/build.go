package commands

import (
	"fmt"
	"os"

	"github.com/project-flogo/cli/api"
	"github.com/project-flogo/cli/common"
	"github.com/spf13/cobra"
)

var buildShim string

var buildOptimize bool
var buildEmbed bool

func init() {
	buildCmd.Flags().StringVarP(&buildShim, "shim", "", "", "trigger shim")
	buildCmd.Flags().BoolVarP(&buildOptimize, "optimize", "o", false, "optimize build")
	buildCmd.Flags().BoolVarP(&buildEmbed, "embed", "e", false, "embed config")
	rootCmd.AddCommand(buildCmd)
}

//Build the project.
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "build the flogo application",
	Long:  `Build the flogo application.`,
	Run: func(cmd *cobra.Command, args []string) {

		options := api.BuildOptions{Shim: buildShim, OptimizeImports: buildOptimize, EmbedConfig: buildEmbed}

		err := api.BuildProject(common.CurrentProject(), options)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}
