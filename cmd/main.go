package main

import (
	"errors"
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

var (
	configPaths      = []string{"./config.toml", "/etc/zeebe/config.toml"}
	ResourceNotFound = errors.New("Resource at the given path not found.")
	ConfigNotFound   = errors.New("Configuration file not found.")
)

func logAndDie(err error) {
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

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

func (cf *Config) String() string {
	return fmt.Sprintf("Version: %s\tBroker: %s", cf.Version, cf.Broker.String())

}

func sendCreateTask(client *zbc.Client, topic string, m *zbc.CreateTask) (*zbc.Message, error) {
	commandRequest := zbc.NewCreateTaskMessage(&sbe.ExecuteCommandRequest{
		PartitionId: 0,
		Key:         0,
		EventType:   sbe.EventTypeEnum(0),
		TopicName:   []uint8(topic),
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
	log.Printf("Loading resource at %s\n", path)
	if len(path) == 0 {
		return nil, ResourceNotFound
	}

	filename, _ := filepath.Abs(path)
	yamlFile, _ := ioutil.ReadFile(filename)

	var command zbc.CreateTask
	err := yaml.Unmarshal(yamlFile, &command)
	if err != nil {
		return nil, err
	}
	return &command, nil
}

func loadConfig() (*Config, error) {
	for i := 0; i < len(configPaths); i++ {
		log.Printf("Loading configuration at %s\n", configPaths[i])

		var config Config
		if _, err := toml.DecodeFile(configPaths[i], &config); err != nil {
			continue
		}
		return &config, nil
	}
	return nil, ConfigNotFound
}

func main() {
	config, err := loadConfig()
	logAndDie(err)
	log.Println(config.String())

	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name:    "create",
			Aliases: []string{"c"},
			Usage:   "create a resource",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "topic, t",
					Value:  "default-topic",
					Usage:  "Executing command request on specific topic.",
					EnvVar: "ZB_TOPIC_NAME",
				},
			},
			Action: func(c *cli.Context) error {
				createTask, err := loadCommandYaml(c.Args().First())
				logAndDie(err)

				client, err := zbc.NewClient(config.Broker.String())
				logAndDie(err)
				log.Println("Connected to Zeebe.")

				response, err := sendCreateTask(client, c.String("topic"), createTask)
				logAndDie(err)

				log.Println("Success. Received response:")
				log.Println(*response.Data)
				return nil
			},
		},
	}
	app.Run(os.Args)
}
