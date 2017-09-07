# ZBC

[![Go Report Card](https://goreportcard.com/badge/github.com/zeebe-io/zbc-go)](https://goreportcard.com/report/github.com/zeebe-io/zbc-go)
[![codecov](https://codecov.io/gh/zeebe-io/zbc-go/branch/master/graph/badge.svg)](https://codecov.io/gh/zeebe-io/zbc-go)
[![Build Status](https://travis-ci.org/zeebe-io/zbc-go.svg?branch=master)](https://travis-ci.org/zeebe-io/zbc-go)
[![GoDoc - Client](http://godoc.org/github.com/zeebe-io/zbc-go/zbc?status.svg)](https://godoc.org/github.com/zeebe-io/zbc-go/zbc)
[![Slack](https://zeebe-slackin.herokuapp.com/badge.svg)](https://zeebe-slackin.herokuapp.com/)

Beware of the little gnomes, flying cats and zebras with wings!


## DISCLAIMER
This is a work in progress and **NOT** meant for any production use!


## Installation

To use as a library, the usual ...

```go get github.com/zeebe-io/zbc-go```

or

```git clone git@github.com:zeebe-io/zbc-go.git```

### Building ```zbctl```

```
cd $GOPATH/src/github.com/zeebe-io/zbc-go
make build
```
This will build ```zbctl```. You can find it inside ```target/bin```.

or

you can install it globally with.
```
cd $GOPATH/src/github.com/zeebe-io/zbc-go
make build
sudo make install
```


## Usage

```
NAME:
   zbctl - Zeebe control client application

USAGE:
   zbctl [global options] command [command options] [arguments...]

VERSION:
   0.2.0

AUTHOR:
   Zeebe Team <info@zeebe.io>

COMMANDS:
     create, c       create a resource
     subscribe, sub  subscribe to a task or topic
     describe, desc  describe a resource
     help, h         Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config value, --cfg value  Location of the configuration file. (default: "/etc/zeebe/config.toml") [$ZBC_CONFIG]
   --help, -h                   show help
   --version, -v                print the version
```

To execute a command, first describe a resource as a yaml file (look at examples folder for examples).

To create a task:

```
$ zbctl create task examples/create-task.yaml
CREATED
```

To deploy a workflow:

```
$ zbctl create workflow examples/demoProcess.bpmn
DEPLOYMENT_CREATED
```

To create a workflow instance:

```
$ zbctl create instance --topic default-topic examples/create-workflow-instance.yaml
WORKFLOW_INSTANCE_CREATED
```

To point your ```zbctl``` to some other broker edit ```config.toml``` which can be find in the ```/etc/zeebe/config.toml```.


## Contributing

  * Get started by checking our [contribution guidelines](https://github.com/zeebe-io/zbc-go/blob/master/CONTRIBUTING.md).
  * Make sure you follow our [code of conduct](https://github.com/zeebe-io/zbc-go/blob/master/CODE_OF_CONDUCT.md).
  * Checkout client documentation for [developers](http://godoc.org/github.com/zeebe-io/zbc-go/zbc)
  * Checkout zbctl documentation for [ops](http://godoc.org/github.com/zeebe-io/zbc-go/cmd)
  * Read the zbc-go [wiki](https://github.com/zeebe-io/zbc-go/wiki) for more technical and design details.
  * The [Zeebe Protocol Specification](http://www.zeebe.io/) contains a wealth of useful information.
  * If you have any questions, just [ask](https://github.com/zeebe-io/zbc-go/issues)!


## Credits
* Special thanks to Sebastian Menski and Daniel Meyer.
* Big thanks to everyone involved in Zeebe.
