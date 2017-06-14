package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/jsam/zbc-go/zbc"
	"github.com/jsam/zbc-go/zbc/sbe"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type Contact struct {
	Address string `toml:"address"`
	Port    string `toml:"port"`
}

func (c *Contact) String() string {
	return fmt.Sprintf("%s:%s", c.Address, c.Port)
}

type Config struct {
	Version string  `toml:"version"`
	Broker  Contact `toml:"broker"`
}

var configPath = "./config.toml"

func sendCreateTask(client *zbc.Client, m *zbc.CreateTask) (*zbc.Message, error) {
	commandRequest := zbc.NewCreateTaskMessage(&sbe.ExecuteCommandRequest{
		PartitionId: 0,
		Key:         0,
		EventType:   sbe.EventTypeEnum(0),
		TopicName:   []uint8("default-topic"),
		Command:     []uint8{},
	}, m)

	response, err := client.Responder(commandRequest)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return response, nil
}

func loadCommandYaml(path string) (*zbc.CreateTask, error) {
	filename, _ := filepath.Abs(path)
	yamlFile, _ := ioutil.ReadFile(filename)

	var command zbc.CreateTask
	err := yaml.Unmarshal(yamlFile, &command)
	if err != nil {
		return nil, err
	}
	return &command, nil
}

func loadConfig(configPath string) (*Config, error) {
	var config Config
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func main() {
	config, err := loadConfig(configPath)
	if err != nil {
		log.Printf("cannot find config: %+#v", err)
	}

	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name:    "create",
			Aliases: []string{"c"},
			Usage:   "create a resource",
			Action: func(c *cli.Context) error {
				log.Printf("added something %+#v", c.Args())

				client, err := zbc.NewClient(config.Broker.String())
				if err != nil {
					log.Println(err)
				}

				log.Println("Connected to Zeebe.")
				createTask, _ := loadCommandYaml(c.Args().First())

				response, err := sendCreateTask(client, createTask)
				if err != nil {
					log.Printf("something really bad happend: %+#v", err)
				}

				log.Println("Success. Received response:")
				log.Println(*response.Data)
				return nil
			},
		},
	}
	app.Run(os.Args)
}
