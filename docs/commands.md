# Commands

- [build](#build) - Build the flogo application
- [create](#create) - Create a flogo application project
- [help](#help)  - Help about any command
- [imports](#imports) - Manage project dependency imports
- [install](#install) - Install a flogo contribution/dependency
- [list](#list) - List installed flogo contributions
- [plugin](#plugin) - Manage CLI plugins
- [update](#update) - Update an application dependency

### Global Flags
```
  --verbose   verbose output
```

  
## build

This command is used to build the application.

```
Usage:
  flogo build [flags]

Flags:
  -e, --embed         embed configuration in binary
  -f, --file string   specify flogo.json file
  -o, --optimize      optimize build
      --shim string   specify shim trigger   
```
_**Note:** the optimize flag removes unused trigger, acitons and activites from the built binary._


### Examples
Build the working application

```bash
$ flogo build
```
Build an application directly from a flogo.json

```bash
$ flogo build -f flogo.json
```
_**Note:** this command will only generate the application binary for the specified json and can be run outside of a flogo application project_

## create

This command is used to create a flogo application project.

```
Usage:
  flogo create [flags] [appName]

Flags:
      --cv string     core library version (ex. master)
  -f, --file string   specify flogo.json file
```

_**Note:** when using the --cv flag to specify a version, the exact version specified might not be used the project.  The application will install the version that satisfies all the dependency constraints.  Typically this flag is used when trying to use the master version of the core library._

### Examples

Create a base sample project with a specific name:

```
$ flogo create my_app
```

Create a project from an existing flogo application descriptor:

```
$ flogo create -f myapp.json
```

## help

This command shows help for any flogo commands.

```
Usage:
  flogo help [command]
```  

### Examples
Get help for the build command:

```bash
$ flogo help build
```
## imports

This command helps manage with project depedency imports

```
Usage:
  flogo imports [command]

Available Commands:
  sync     sync all go imports to project imports
  resolve  resolve project import versions
  list     list all project imports
```   

## install

This command is used to install a flogo contribution or dependency.

```
Usage:
  flogo install [flags] <contribution|dependency>

Flags:
  -f, --file string    specify contribution bundle
  -l, --local string   local path to contribution
```
      
### Examples
Install the basic REST trigger:

```bash
$ flogo install github.com/project-flogo/contrib/trigger/rest
```
Install a contribution that you are currently developing on your computer:

```bash
$ flogo install -l /tmp/dev/myactivity github.com/myuser/myactivity
```

## list

This command lists installed contributions in your application

```
Usage:
  flogo list [flags]

Flags:
  -f, --filter string   apply list filter [used, unused]
  -j, --json            print in json format (default true)
      --orphaned        list orphaned refs
```  
_**Note** orphaned refs are `ref` entries that use an import alias (ex. `"ref": "#log"`) which has no corresponding import._

### Examples
List all installed contributions:

```bash
$ flogo list
```
List all contributions directly used by the application:

```bash
$ flogo list -f used
```
_**Note:** the results of this command are the only contributions that will be compiled into your application when using `flogo build` with the optimize flag_


## plugin

This command is used to install a plugin to the Flogo CLI.

```
Usage:
  flogo plugin [command]

Available Commands:
  install     install cli plugin
  list        list installed plugins
  update      update plugin
```      

### Examples
List all installed plugins:

```bash
$ flogo plugin list
```
Install the legacy support plugin:

```bash
$ flogo plugin install github.com/project-flogo/legacybridge/cli`
```
_**Note:** more information on the legacy support plugin can be found [here](https://github.com/project-flogo/legacybridge/tree/master/cli)_

Install and use custom plugin:

```
$ flogo plugin install github.com/myuser/myplugin

$ flogo `your_command`
```

More information on Flogo CLI plugins can be found [here](plugins.md)
