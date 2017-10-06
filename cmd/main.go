package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"

	"github.com/BurntSushi/toml"
	"github.com/urfave/cli"
	"github.com/zeebe-io/zbc-go/zbc"
	"github.com/zeebe-io/zbc-go/zbc/zbmsgpack"
	"os/signal"
	"text/tabwriter"
)

const (
	defaultConfiguration = "config.toml"
)

var (
	errResourceNotFound = errors.New("Resource at the given path not found")
	version             = "dev"
	commit              = "none"
)

func isFatal(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type contact struct {
	Address string `toml:"address"`
	Port    string `toml:"port"`
}

func (c *contact) String() string {
	return fmt.Sprintf("%s:%s", c.Address, c.Port)
}

type config struct {
	Version string  `toml:"version"`
	Broker  contact `toml:"broker"`
}

func (cf *config) String() string {
	return fmt.Sprintf("version: %s\tBroker: %s", cf.Version, cf.Broker.String())

}

func loadCommandYaml(path string, command interface{}) error {
	yamlFile, _ := loadFile(path)

	err := yaml.Unmarshal(yamlFile, command)
	if err != nil {
		return err
	}
	return nil
}

func loadFile(path string) ([]byte, error) {
	log.Printf("Loading resource at %s\n", path)
	if len(path) == 0 {
		return nil, errResourceNotFound
	}

	filename, _ := filepath.Abs(path)
	return ioutil.ReadFile(filename)
}

func loadConfig(path string, c *config) {
	if _, err := toml.DecodeFile(path, c); err != nil {
		log.Printf("Reading configuration failed. Expecting to found configuration file at %s\n", path)
		log.Printf("HINT: Configuration file is not in place. Try setting configuration path with:")
		log.Printf(" zbctl --config <path to config.toml>")
		log.Printf("Using 0.0.0.0:51015 as default.")
		c.Broker = contact{
			Address: "0.0.0.0",
			Port:    "51015",
		}
	}
}

func main() {
	var conf config
	log.SetOutput(ioutil.Discard)

	app := cli.NewApp()
	app.Usage = "Zeebe control client application"
	app.Version = fmt.Sprintf("%s (%s)", version, commit)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config, cfg",
			Value:  defaultConfiguration,
			Usage:  "Location of the configuration file.",
			EnvVar: "ZBC_CONFIG",
		},
	}
	app.Before = cli.BeforeFunc(func(c *cli.Context) error {
		loadConfig(c.String("config"), &conf)
		log.Println(conf.String())
		return nil
	})

	app.Authors = []cli.Author{
		{Name: "Zeebe Team", Email: "info@zeebe.io"},
	}

	app.Commands = []cli.Command{
		{
			Name:    "create",
			Aliases: []string{"c"},
			Usage:   "create a resource",
			Subcommands: []cli.Command{
				{
					Name:    "topic",
					Aliases: []string{"tp"},
					Usage:   "Create a new topic",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:   "partitions, p",
							Value:  "1",
							Usage:  "Specify the number of partitions for new topic",
							EnvVar: "",
						},
						cli.StringFlag{
							Name:   "name, n",
							Value:  "",
							Usage:  "Specify the name of the new topic",
							EnvVar: "",
						},
					},
					Action: func(c *cli.Context) error {
						client, err := zbc.NewClient(conf.Broker.String())
						isFatal(err)
						log.Println("Connected to Zeebe.")
						// TODO: check if there is any value of those flags
						topic, err := client.CreateTopic(c.String("name"), c.Int("partitions"))
						isFatal(err)

						fmt.Println(topic.State)
						return nil
					},
				},
				{
					Name:    "task",
					Aliases: []string{"ts"},
					Usage:   "Create a new task using the given YAML file",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:   "topic, t",
							Value:  "default-topic",
							Usage:  "Executing command request on specific topic.",
							EnvVar: "ZB_TOPIC_NAME",
						},
					},
					Action: func(c *cli.Context) error {
						var task zbmsgpack.Task
						err := loadCommandYaml(c.Args().First(), &task)
						isFatal(err)
						client, err := zbc.NewClient(conf.Broker.String())
						isFatal(err)
						log.Println("Connected to Zeebe.")

						response, err := client.CreateTask(c.String("topic"), &task)
						isFatal(err)

						if response.Data != nil {
							task := response.Task()
							fmt.Println(task.State)
						} else {
							fmt.Println("err: received nil response")
						}

						return nil
					},
				},
				{
					Name:    "workflow",
					Aliases: []string{"wf"},
					Usage:   "Create a new workflow using the given bpmn file",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:   "topic, t",
							Value:  "default-topic",
							Usage:  "Executing command request on specific topic.",
							EnvVar: "ZB_TOPIC_NAME",
						},
					},
					Action: func(c *cli.Context) error {
						filename := c.Args().First()

						var resourceType string
						switch filepath.Ext(filename) {
						case ".yaml", ".yml":
							resourceType = zbc.YamlWorkflow
						default:
							resourceType = zbc.BpmnXml
						}

						client, err := zbc.NewClient(conf.Broker.String())
						isFatal(err)

						response, err := client.CreateWorkflowFromFile(c.String("topic"), resourceType, filename)
						isFatal(err)

						if response.Data != nil {
							m, _ := response.ParseToMap()
							if state, ok := (*m)["state"]; ok {
								fmt.Println(state)
							}
						} else {
							fmt.Println("err: received nil response")
						}
						return nil
					},
				},
				{
					Name:    "instance",
					Aliases: []string{"wfi"},
					Usage:   "Create a new workflow instance using the given YAML file",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:   "topic, t",
							Value:  "default-topic",
							Usage:  "Executing command request on specific topic.",
							EnvVar: "ZB_TOPIC_NAME",
						},
					},
					Action: func(c *cli.Context) error {
						var workflowInstance zbmsgpack.WorkflowInstance
						err := loadCommandYaml(c.Args().First(), &workflowInstance)
						isFatal(err)

						client, err := zbc.NewClient(conf.Broker.String())
						isFatal(err)
						log.Println("Connected to Zeebe.")

						response, err := client.CreateWorkflowInstance(c.String("topic"), &workflowInstance)
						isFatal(err)

						log.Println("Success. Received response:")

						if response.Data != nil {
							m, _ := response.ParseToMap()
							state, ok := (*m)["state"]
							if ok {
								fmt.Println(state)
							}
							if state == zbc.WorkflowInstanceRejected {
								fmt.Println("\nDid you deploy the workflow?")
							}
						} else {
							fmt.Println("err: received nil response")
						}
						return nil
					},
				},
			},
		},
		{
			Name:    "subscribe",
			Aliases: []string{"sub"},
			Usage:   "subscribe to a task or topic",
			Subcommands: []cli.Command{
				{
					Name:    "task",
					Aliases: []string{"ts"},
					Usage:   "open a task subscription",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:   "topic, t",
							Value:  "default-topic",
							Usage:  "Executing command request on specific topic.",
							EnvVar: "ZB_TOPIC_NAME",
						},
						cli.StringFlag{
							Name:   "lock-owner, l",
							Value:  "zbc",
							Usage:  "Specify lock owner.",
							EnvVar: "ZB_LOCK_OWNER",
						},
						cli.StringFlag{
							Name:   "task-type, tt",
							Value:  "foo",
							Usage:  "Specify task type.",
							EnvVar: "ZB_TASK_TYPE",
						},
					},
					Action: func(c *cli.Context) error {
						client, err := zbc.NewClient(conf.Broker.String())
						isFatal(err)

						subscriptionCh, subscription, err := client.TaskConsumer(c.String("topic"), c.String("lock-owner"), c.String("task-type"))
						isFatal(err)

						osCh := make(chan os.Signal, 1)
						signal.Notify(osCh, os.Interrupt)
						go func() {
							<-osCh
							fmt.Println("Closing subscription.")
							_, err := client.CloseTaskSubscription(subscription)
							if err != nil {
								fmt.Println("failed to close subscription: ", err)
							} else {
								fmt.Println("Subscription closed.")
							}
							os.Exit(0)
						}()

						log.Println("Waiting for events ....")
						for {
							message := <-subscriptionCh
							fmt.Println(message.String())
						}
					},
				},
				{
					Name:    "topic",
					Aliases: []string{"tp"},
					Usage:   "open a topic subscription",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:   "topic, t",
							Value:  "default-topic",
							Usage:  "",
							EnvVar: "ZB_TOPIC_NAME",
						},
						cli.StringFlag{
							Name:   "subscription-name, sn",
							Value:  "default-name",
							Usage:  "",
							EnvVar: "ZB_TOPIC_NAME",
						},
						cli.Int64Flag{
							Name:   "start-position, sp",
							Value:  -1,
							Usage:  "",
							EnvVar: "ZB_TOPIC_SUB_START_POSITION",
						},
					},
					Action: func(c *cli.Context) error {
						client, err := zbc.NewClient(conf.Broker.String())
						isFatal(err)
						log.Println("Connected to Zeebe.")
						subscriptionCh, sub, err := client.TopicConsumer(c.String("topic"), c.String("subscription-name"), c.Int64("start-position"))
						isFatal(err)

						osCh := make(chan os.Signal, 1)
						signal.Notify(osCh, os.Interrupt)
						go func() {
							<-osCh
							fmt.Println("Closing subscription.")
							_, err := client.CloseTopicSubscription(sub)
							if err != nil {
								fmt.Println("failed to close subscription: ", err)
							} else {
								fmt.Println("Subscription closed.")
							}
							os.Exit(0)
						}()

						for {
							message := <-subscriptionCh
							fmt.Println(message.String())
						}

					},
				},
			},
		},
		{
			Name:    "describe",
			Aliases: []string{"desc"},
			Usage:   "describe a resource",
			Subcommands: []cli.Command{
				{
					Name:    "topology",
					Aliases: []string{"t"},
					Usage:   "check cluster topology",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:   "topic, t",
							Value:  "default-topic",
							Usage:  "Executing command request on specific topic.",
							EnvVar: "ZB_TOPIC_NAME",
						},
					},
					Action: func(c *cli.Context) error {
						client, err := zbc.NewClient(conf.Broker.String())
						isFatal(err)

						topology, err := client.Topology()
						isFatal(err)

						w := tabwriter.NewWriter(os.Stdout, 0, 0, 10, ' ', tabwriter.TabIndent)

						fmt.Fprintln(w, "Topic Name\tBroker\tPartitionID")
						for key, value := range (*topology).TopicLeaders {
							for _, leader := range value {
								fmt.Fprintln(w, fmt.Sprintf("%s\t%s\t%d", key, leader.Addr(), leader.PartitionID))
							}
						}

						fmt.Fprintln(w, "")
						fmt.Fprintln(w, "Brokers")
						for _, broker := range (*topology).Brokers {
							fmt.Fprintln(w, broker.Addr())
						}
						w.Flush()
						return nil
					},
				},
			},
		},
	}
	app.Run(os.Args)
}
