# Zeebe Client
Beware of the little gnomes, flying cats and zebras with wings!


## DISCLAIMER
This is a work in progress and **NOT** meant for any production use!


## Installation

to use as a library, the usual ...

```go get github.com/jsam/zbc-go```

to use zbctl (make sure you are on $GOPATH)

```
git clone git@github.com:jsam/zbc-go.git
cd zbc-go
make build
cd target/bin
./zbctl create <resource>
```


## Usage 

Read the code, it's awesome!

```zbctl``` can be installed by ```make build``` (output will be inside target/bin)

To execute a command, first describe a resource as a yaml file and then:

```
cd target/bin
./zbctl create ../../examples/create-task.yaml
```

To point your zbctl to some other broker edit ```config.toml``` which can be find in the target/bin as well.
## Credits
* Special thanks to Sebastian Menski and Daniel Meyer.
* Big thanks to the entire team of Zeebe.