# ZBC

[![Go Report Card](https://goreportcard.com/badge/github.com/jsam/zbc-go)](https://goreportcard.com/report/github.com/jsam/zbc-go)
[![Build Status](https://travis-ci.org/jsam/zbc-go.svg?branch=master)](https://travis-ci.org/jsam/zbc-go)
[![GoDoc](https://godoc.org/github.com/jsam/zbc-go?status.svg)](https://godoc.org/github.com/jsam/zbc-go)

Beware of the little gnomes, flying cats and zebras with wings!


## DISCLAIMER
This is a work in progress and **NOT** meant for any production use!


## Installation

To use as a library, the usual ...

```go get github.com/jsam/zbc-go```

or 

```git clone git@github.com:jsam/zbc-go.git```

### Building ```zbctl```

```
cd $GOPATH/src/github.com/jsam/zbc-go 
make build
```
This will build ```zbctl```. You can find it inside ```target/bin```.

or

you can install it globally with.
```
cd $GOPATH/src/github.com/jsam/zbc-go
make build
sudo make install
```


## Usage 

To execute a command, first describe a resource as a yaml file (look at examples folder for examples) and then:

```
zbctl create --topic default-topic examples/create-task.yaml
```

To point your ```zbctl``` to some other broker edit ```config.toml``` which can be find in the ```/etc/zeebe/config.toml```.


## Contributing

  * Get started by checking our [contribution guidelines](https://github.com/jsam/zbc-go/blob/master/CONTRIBUTING.md).
  * Make sure you follow our [code of conduct](https://github.com/jsam/zbc-go/blob/master/CODE_OF_CONDUCT.md).
  * Read the zbc-go [wiki](https://github.com/jsam/zbc-go/wiki) for more technical and design details.
  * The [Zeebe Protocol Specification](http://www.zeebe.io/) contains a wealth of useful information.
  * If you have any questions, just [ask](https://github.com/jsam/zbc-go/issues)!


## Credits
* Special thanks to Sebastian Menski and Daniel Meyer.
* Big thanks to everyone involved in Zeebe.
