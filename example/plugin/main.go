package walrus

import (
	"fmt"

	"github.com/project-flogo/cli/common"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func GetWalrus() {
	//fmt.Println("Log Example")

	log.WithFields(log.Fields{
		"animal": "walrus",
	}).Info("A walrus appears")
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
