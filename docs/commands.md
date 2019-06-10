<!--
title: flogo
weight: 5020
pre: "<i class=\"fa fa-terminal\" aria-hidden=\"true\"></i> "
-->

# Commands

- [build](#build) - Build the flogo application
- [create](#create) - Create a flogo application project
- [help](#help)  - Help about any command
- [imports](#imports) - Manage project dependency imports
- [install](#install) - Install a flogo contribution/dependency
- [list](#list) - List installed flogo contributions
- [plugin](#plugin) - Manage CLI plugins
- [update](#update) - Update an application contribution/dependency

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
  -f, --file string   specify a flogo.json to build
  -o, --optimize      optimize build
      --shim string   use shim trigger   
```
_**Note:** the optimize flag removes unused trigger, acitons and activites from the built binary._


### Examples
Build the current project application

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
      --cv string     specify core library version (ex. master)
  -f, --file string   specify a flogo.json to create project from
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

This command helps manage project imports of contributions and dependencies.

```
Usage:
  flogo imports [command]

Available Commands:
  sync     sync Go imports to project imports
  resolve  resolve project imports to installed version
  list     list project imports
```   

## install

This command is used to install a flogo contribution or dependency.

```
Usage:
  flogo install [flags] <contribution|dependency>

Flags:
  -f, --file string      specify contribution bundle
  -r, --replace string   specify path to replacement contribution/dependency
```
      
### Examples
Install the basic REST trigger:

```bash
$ flogo install github.com/project-flogo/contrib/trigger/rest
```
Install a contribution that you are currently developing on your computer:

```bash
$ flogo install -r /tmp/dev/myactivity github.com/myuser/myactivity
```

Install a contribution that is being developed by different person on their fork:

```bash
$ flogo install -r github.com/otherusr/myactivity@master github.com/myuser/myactivity
```

## list

This command lists installed contributions in your application

```
Usage:
  flogo list [flags]

Flags:
      --filter string   apply list filter [used, unused]
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
$ flogo list --filter used
```
_**Note:** the results of this command are the only contributions that will be compiled into your application when using `flogo build` with the optimize flag_


## plugin

This command is used to install a plugin to the Flogo CLI.

```
Usage:
  flogo plugin [command]

Available Commands:
  install     install CLI plugin
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
<br>
More information on Flogo CLI plugins can be found [here](plugins.md)

## update

This command updates a contribution or dependency in the project.

```
Usage:
  flogo update [flags] <contribution|dependency>

```   
### Examples
Update you log activity to master:

```bash
$ flogo update github.com/project-flogo/contrib/activity/log@master
```

Update your flogo core library to latest master:

```bash
$ flogo update github.com/project-flogo/core@master
```
