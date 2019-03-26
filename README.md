<p align="center">
  <img src ="https://raw.githubusercontent.com/TIBCOSoftware/flogo/master/images/projectflogo.png" />
</p>

<p align="center" >
  <b>Serverless functions and edge microservices made painless</b>
</p>

<p align="center">
  <img src="https://travis-ci.org/TIBCOSoftware/flogo-cli.svg"/>
  <img src="https://img.shields.io/badge/dependencies-up%20to%20date-green.svg"/>
  <img src="https://img.shields.io/badge/license-BSD%20style-blue.svg"/>
  <a href="https://gitter.im/project-flogo/Lobby?utm_source=share-link&utm_medium=link&utm_campaign=share-link"><img src="https://badges.gitter.im/Join%20Chat.svg"/></a>
</p>

<p align="center">
  <a href="#Installation">Installation</a> | <a href="#getting-started">Getting Started</a> | <a href="#documentation">Documentation</a> | <a href="#contributing">Contributing</a> | <a href="#license">License</a>
</p>

<br/>
Project Flogo is an open source framework to simplify building efficient & modern serverless functions and edge microservices and this is the cli that makes it all happen. 

FLOGO CLI
======================

The Flogo CLI is the primary tool to use when working with a Flogo application.  It is used to create, modify and build Flogo applications
## Installation
### Prerequisites
To get started with the Flogo CLI you'll need to have a few things
* The Go programming language version 1.11 or later should be [installed](https://golang.org/doc/install).
* In order to simplify dependency management, we're using **go mod**. For more information on **go mod**, visit the [Go Modules Wiki](https://github.com/golang/go/wiki/Modules).

### Install the cli
To install the CLI, simply open a terminal and enter the below command

```
$ go get -u github.com/project-flogo/cli/...
```

_Note that the -u parameter automatically updates the cli if it exists_

### Build the CLI from source
You can build the cli from source code as well, which is convenient if you're developing new features for it! To do that, follow these easy steps

```bash
# Get the flogo CLI from GitHub
$ git clone https://github.com/project-flogo/cli.git

# Go to the directory
$ cd cli

# Optionally check out the branch you want to use 
$ git checkout test_branch

# Run the install command
$ go install ./... 
```

## Getting started
Getting started should be easy and fun, and so is getting started with the Flogo cli. 

First, create a file called `flogo.json` and with the below content (which is a simple app with an [HTTP trigger](https://tibcosoftware.github.io/flogo/development/webui/triggers/rest/))

```json
{
  "name": "SampleApp",
  "type": "flogo:app",
  "version": "0.0.1",
  "appModel": "1.1.0",
  "imports": [
  	"github.com/project-flogo/contrib/trigger/rest",
  	"github.com/project-flogo/flow",
  	"github.com/project-flogo/contrib/activity/log"
  ],
  "triggers": [
    {
      "id": "receive_http_message",
      "ref": "#rest",
      "name": "Receive HTTP Message",
      "description": "Simple REST Trigger",
      "settings": {
        "port": 9233
      },
      "handlers": [
        {
          "settings": {
            "method": "GET",
            "path": "/test"
          },
          "action": {
            "ref": "#flow",
            "settings": {
              "flowURI": "res://flow:sample_flow"
            }
          }
        }
      ]
    }
  ],
  "resources": [
    {
      "id": "flow:sample_flow",
      "data": {
        "name": "SampleFlow",
        "tasks": [
          {
            "id": "log_message",
            "name": "Log Message",
            "description": "Simple Log Activity",
            "activity": {
              "ref": "#log",
              "input": {
                "message": "Simple Log",
                "addDetails": "false"
              }
            }
          }
        ]
      }
    }
  ]
}
```

Based on this file we'll create a new flogo app

```bash
$ flogo create -f flogo.json myApp
```

From the app folder we can build the executable

```bash
$ cd myApp
$ flogo build -e
```

Now that there is an executable we can run it!

```bash
$ cd bin
$ ./myApp
```

The above commands will start the REST server and wait for messages to be sent to `http://localhost:9233/test`. To send a message you can use your browser, or a new terminal window and run

```bash
$ curl http://localhost:9233/test
```

_For more tutorials check out the [Labs](https://tibcosoftware.github.io/flogo/labs/) section in our documentation_

## Documentation

There is documentation also available for [CLI Commands](docs/commands.md) and [CLI Plugins](docs/plugins.md).

## Contributing
Want to contribute to Project Flogo? We've made it easy, all you need to do is fork the repository you intend to contribute to, make your changes and create a Pull Request! Once the pull request has been created, you'll be prompted to sign the CLA (Contributor License Agreement) online.

Not sure where to start? No problem, you can browse the Project Flogo repos and look for issues tagged `kind/help-wanted` or `good first issue`. To make this even easier, we've added the links right here too!
* Project Flogo: [kind/help-wanted](https://github.com/TIBCOSoftware/flogo/issues?q=is%3Aissue+is%3Aopen+label%3A%22help+wanted%22) and [good first issue](https://github.com/TIBCOSoftware/flogo/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22)
* flogo cli: [kind/help-wanted](https://github.com/project-flogo/cli/issues?q=is%3Aissue+is%3Aopen+label%3A%22help+wanted%22) and [good first issue](https://github.com/project-flogo/cli/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22)
* flogo core: [kind/help-wanted](https://github.com/project-flogo/core/issues?q=is%3Aissue+is%3Aopen+label%3A%22help+wanted%22) and [good first issue](https://github.com/project-flogo/core/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22)
* flogo contrib: [kind/help-wanted](https://github.com/project-flogo/contrib/issues?q=is%3Aissue+is%3Aopen+label%3A%22help+wanted%22) and [good first issue](https://github.com/project-flogo/contrib/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22)

Another great way to contribute to Project Flogo is to check [flogo-contrib](https://github.com/project-flogo/contrib). That repository contains some basic contributions, such as activities, triggers, etc. Perhaps there is something missing? Create a new activity or trigger or fix a bug in an existing activity or trigger.

If you have any questions, feel free to post an issue and tag it as a question, email flogo-oss@tibco.com or chat with the team and community:

* The [project-flogo/Lobby](https://gitter.im/project-flogo/Lobby?utm_source=share-link&utm_medium=link&utm_campaign=share-link) Gitter channel should be used for general discussions, start here for all things Flogo!
* The [project-flogo/developers](https://gitter.im/project-flogo/developers?utm_source=share-link&utm_medium=link&utm_campaign=share-link) Gitter channel should be used for developer/contributor focused conversations. 

For additional details, refer to the [Contribution Guidelines](https://github.com/TIBCOSoftware/flogo/blob/master/CONTRIBUTING.md).

## License 
Flogo source code in [this](https://github.com/project-flogo/cli) repository is under a BSD-style license, refer to [LICENSE](https://github.com/project-flogo/cli/blob/master/LICENSE) 
