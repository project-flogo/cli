package registry

import (
	"github.com/spf13/cobra"
)

var commands []*cobra.Command

func RegisterCommands(command *cobra.Command) {
	commands = append(commands, command)
}

func GetCommands() []*cobra.Command {

	return commands
}
