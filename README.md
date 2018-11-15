## Installation

### Prerequisites
To get started with the Project Flogo cli you'll need 
* The Go programming language version 1.11 or later 

### Install the cli
To install the cli, simply open a terminal and enter the below command
```
$ go get -u github.com/project-flogo/cli/...
```
_Note that the -u parameter automatically updates the cli if it exists_

## Commands

### build

This command is used to build the application.

```
$ flogo build
```

options
```
-e, --embed         embed config
-h, --help          help for build
-o, --optimize      optimize build
      --shim string   trigger shim
```

### create

This command is used to create a flogo application project.

_Create the base sample project with a specific name_ 
```
$ flogo create my_app
```

_Create a flogo application project from an existing flogo application descriptor_

```
$ flogo create -f myapp.json
```

options
```
-c, --core string   specify core library version [master| v0.0.1]
-f, --file string   path to flogo.json file

```
### help

This commad shows help for any flogo commands.

```
$ flogo help build
```

### install

This command is used to install a contribution to your project.

```
$ flogo install github.com/skothari-tibco/csvtimer
```

### plugin

This command is used to install a plugin to your cli.

```
$ flogo plugin install github.com/skothari-tibco/walrus

$ flogo walrus
```

### Global Flags
```
-- verbose verbose output
```

## Creating a Flogo Cli Plugin.

* To get started, create a sample `cobra` program.
* Import the `github.com/project-flogo/cli/common` packages
* Implement `init()` function to register to the flogo cli. You can also register subcommands.

Sample
```go
package walrus

import (
	"fmt"
    "github.com/project-flogo/cli/common" //Import Flogo Cli 
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func GetWalrus() {
	log.WithFields(log.Fields{
		"animal": "walrus",
	}).Info("A walrus appears")
}

var helloCmd = &cobra.Command{
	Use:              "walrus",
	Short:            "says walrus",
	Long:             `This subcommand says walrus`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {}, // Add any functions you want to run before running the command. If not leave blank.
	Run: func(cmd *cobra.Command, args []string) {
        //Your Main Function
		GetWalrus()
	},
}

func init() {
	common.RegisterPlugin(helloCmd) // Register your main command

	helloCmd.AddCommand(sayCmd) // Register your sub commands
}

var sayCmd = &cobra.Command{
	Use:   "say",
	Short: "says walrus",
	Long:  `This subcommand says walrus`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("This is sub command")
	},
}
```

* Host your Repo.
* Install the plugin using following command:
```
$ flogo plugin install github.com/skothari-tibco/walrus
```

