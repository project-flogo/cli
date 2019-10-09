package common

import (
	"github.com/spf13/cobra"
)

var commands []*cobra.Command
var pluginPkgs []string

func RegisterPlugin(command *cobra.Command) {
	commands = append(commands, command)
}

func GetPlugins() []*cobra.Command {

	tmp := make([]*cobra.Command, len(commands))
	copy(tmp, commands)

	return tmp
}

func GetPluginPkgs() []string{
	return pluginPkgs
}
