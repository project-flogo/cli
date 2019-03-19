package walrus

import (
	"fmt"

	"github.com/project-flogo/cli/common"
	"github.com/spf13/cobra"
)

func GetWalrus() {
	//fmt.Println("Log Example")

	fmt.Println("Log A Walrus")
}

var helloCmd = &cobra.Command{
	Use:              "walrus",
	Short:            "says walrus",
	Long:             `This subcommand says walrus`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {},
	Run: func(cmd *cobra.Command, args []string) {
		GetWalrus()
	},
}

func init() {
	common.RegisterPlugin(helloCmd)

	helloCmd.AddCommand(sayCmd)
}

var sayCmd = &cobra.Command{
	Use:   "say",
	Short: "says walrus",
	Long:  `This subcommand says walrus`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("This is sub command")
	},
}
