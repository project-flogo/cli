<!--
title: plugins
weight: 5020
pre: "<i class=\"fa fa-terminal\" aria-hidden=\"true\"></i> "
-->

# Plugins

The Flogo CLI has support for plugins.  These plugins can be used to extend the Flogo CLI command.

## Creating a CLI plugin

First lets setup the go project:

```bash
# Create a directory for your plugin project
$ mkdir myplugin

# Go to the directory
$ cd myplugin

# Initialize the Go module information
$ go mod init github.com/myuser/myplugin

# Edit/Create the plugin code
$ vi myplugin.go
```

Next lets create the code for our simple plugin:

```go
package myplugin

import (
	"fmt"
	"github.com/project-flogo/cli/common" // Flogo CLI support code
	"github.com/spf13/cobra"
)

func init() {
	common.RegisterPlugin(myCmd)
}

var myCmd = &cobra.Command{
	Use:	"mycmd",
	Short:	"says hello world",
	Long:	"This plugin command says hello world",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello World")
	},
}
```
Once you save the code, we need to fix up the Go Module dependencies.

```bash
$ go mod tidy
```

Now you are ready to test out your plugin.  First you must host your plugin in your git repository.  Then you are ready to install and run your plugin

```
# Install your plugin
$ flogo plugin install github.com/myuser/myplugin

# Run your new plugin command
$ flogo mycmd
```
