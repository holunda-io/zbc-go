package zbc

import (
	"bufio"
	"bytes"
	"errors"
	"log"
	"net"
	"time"

	"github.com/zeebe-io/zbc-go/zbc/zbmsgpack"
	"github.com/zeebe-io/zbc-go/zbc/zbsbe"
)

var (
	errTimeout     = errors.New("Request timeout")
	errSocketWrite = errors.New("Tried to write more bytes to socket")
)

// Client for one Zeebe broker
type Client struct {
	RequestHandler
	ResponseHandler

	Connection    net.Conn
	Cluster       *zbmsgpack.ClusterTopology
	closeCh       chan bool
	transactions  map[uint64]chan *Message
	subscriptions map[uint64]chan *Message
}

func (c *Client) sender(message *Message) error {
	if c.Cluster != nil {
		if time.Since(c.Cluster.UpdatedAt)  > TopologyRefreshInterval * time.Second {
			c.Topology()
		}
	}

	writer := NewMessageWriter(message)
	byteBuff := &bytes.Buffer{}
	writer.Write(byteBuff)

	n, err := c.Connection.Write(byteBuff.Bytes())
	if err != nil {
		return err
	}

	if n != len(byteBuff.Bytes()) {
		return errSocketWrite
	}
	return nil
}

func (c *Client) receiver() {
	for {
		select {
		case <-c.closeCh:
			log.Println("Closing client.")
			c.Connection.Close()
			return

		default:
			c.Connection.SetReadDeadline(time.Now().Add(time.Millisecond * 50))
			buffer := bufio.NewReaderSize(c.Connection, 20000)
			r := NewMessageReader(buffer)

			headers, tail, err := r.ReadHeaders()
			if err != nil {
				continue
			}
			message, err := r.ParseMessage(headers, tail)
			if err != nil && !headers.IsSingleMessage() {
				delete(c.transactions, headers.RequestResponseHeader.RequestID)
				continue
			}

			if !headers.IsSingleMessage() && message != nil {
				c.transactions[headers.RequestResponseHeader.RequestID] <- message
				continue
			}

			if err != nil && headers.IsSingleMessage() {
				// TODO: close and delete subscription
				continue
			}

			if headers.IsSingleMessage() && message != nil {
				subscriberKey := (*message.SbeMessage).(*zbsbe.SubscribedEvent).SubscriberKey
				c.subscriptions[subscriberKey] <- message
				continue
			}
		}
	}
}

// Responder implements synchronous way of sending ExecuteCommandRequest and waiting for ExecuteCommandResponse.
func (c *Client) Responder(message *Message) (*Message, error) {
	respCh := make(chan *Message)
	c.transactions[message.Headers.RequestResponseHeader.RequestID] = respCh

	if err := c.sender(message); err != nil {
		return nil, err
	}

	select {
	case resp := <-c.transactions[message.Headers.RequestResponseHeader.RequestID]:
		delete(c.transactions, message.Headers.RequestResponseHeader.RequestID)
		return resp, nil
	case <-time.After(time.Second * RequestTimeout):
		delete(c.transactions, message.Headers.RequestResponseHeader.RequestID)
		return nil, errTimeout
	}
}

// CreateTask will create new task on specified topic.
func (c *Client) CreateTask(topic string, m *zbmsgpack.Task) (*Message, error) {
	commandRequest := &zbsbe.ExecuteCommandRequest{
		PartitionId: c.Cluster.TopicLeaders[topic].PartitionID,
		Position:    0,
		TopicName:   []uint8(topic),
		Command:     []uint8{},
	}
	commandRequest.Key = commandRequest.KeyNullValue()
	message := c.createTaskRequest(commandRequest, m)

	return MessageRetry(func() (*Message, error) {
		return c.Responder(message)
	})
}

// CreateWorkflowInstance will create new workflow instance on the broker.
func (c *Client) CreateWorkflowInstance(topic string, m *zbmsgpack.WorkflowInstance) (*Message, error) {
	commandRequest := &zbsbe.ExecuteCommandRequest{
		PartitionId: c.Cluster.TopicLeaders[topic].PartitionID,
		Position:    0,
		TopicName:   []uint8(topic),
		Command:     []uint8{},
	}

	commandRequest.Key = commandRequest.KeyNullValue()

	return MessageRetry(func() (*Message, error) {
		return c.Responder(c.createWorkflowInstanceRequest(commandRequest, m))
	})
}

func (c *Client) Deploy(topic string, definition []byte) (*Message, error) {
	deployment := zbmsgpack.Deployment{
		State:   CreateDeployment,
		BPMNXML: definition,
	}
	commandRequest := &zbsbe.ExecuteCommandRequest{
		PartitionId: c.Cluster.TopicLeaders[topic].PartitionID,
		Position:    0,
		TopicName:   []uint8(topic),
		Command:     []uint8{},
	}
	commandRequest.Key = commandRequest.KeyNullValue()

	return MessageRetry(func() (*Message, error) {
		return c.Responder(c.newDeploymentRequest(commandRequest, &deployment))
	})
}

// TaskConsumer opens a subscription on task and returns a channel where all the SubscribedEvents will arrive.
func (c *Client) TaskConsumer(topic, lockOwner, taskType string) (chan *Message, error) {
	taskSub := &zbmsgpack.TaskSubscription{
		TopicName:     topic,
		PartitionID:   c.Cluster.TopicLeaders[topic].PartitionID,
		Credits:       32,
		LockDuration:  300000,
		LockOwner:     lockOwner,
		SubscriberKey: 0,
		TaskType:      taskType,
	}

	subscriptionCh := make(chan *Message, taskSub.Credits)
	msg := c.openTaskSubscriptionRequest(taskSub)

	response, err := MessageRetry(func() (*Message, error) { return c.Responder(msg) })
	if err != nil {
		return nil, err
	}

	d, _ := response.ParseToMap()
	subscriberKey := (*d)["subscriberKey"].(uint64) //t.(zbsbe.SubscribedEvent).SubscriberKey
	c.subscriptions[subscriberKey] = subscriptionCh

	return subscriptionCh, nil
}

// TopicConsumer opens a subscription on topic and returns a channel where all the SubscribedEvents will arrive.
func (c *Client) TopicConsumer(topic, subName string) (chan *Message, error) {
	topicSub := &zbmsgpack.TopicSubscription{
		StartPosition:    -1,
		Name:             subName,
		PrefetchCapacity: 0,
		ForceStart:       false,
		State:            TopicSubscriptionSubscribeState,
	}
	execCommandRequest := &zbsbe.ExecuteCommandRequest{
		PartitionId: c.Cluster.TopicLeaders[topic].PartitionID,
		Position:    0,
		EventType:   zbsbe.EventType.SUBSCRIBER_EVENT,
		TopicName:   []byte(topic),
	}
	execCommandRequest.Key = execCommandRequest.KeyNullValue()

	subscriptionCh := make(chan *Message)
	msg := c.openTopicSubscriptionRequest(execCommandRequest, topicSub)

	response, err := MessageRetry(func() (*Message, error) { return c.Responder(msg) })
	if err != nil {
		return nil, err
	}

	subscriberKey := (*response.SbeMessage).(*zbsbe.ExecuteCommandResponse).Key
	c.subscriptions[subscriberKey] = subscriptionCh

	return subscriptionCh, nil
}

// TopologyRequest will retrieve latest cluster topology information.
func (c *Client) Topology() (*zbmsgpack.ClusterTopology, error) {
	resp, err := MessageRetry(func() (*Message, error) { return c.Responder(c.topologyRequest()) })
	if err != nil {
		return nil, err
	}
	topology := c.newClusterTopologyResponse(resp)
	c.Cluster = topology
	return topology, nil
}

// Start will spinoff receiver in goroutine, which will make client effectively ready to communicate with the broker.
func (c *Client) Start() {
	go c.receiver()
}

func (c *Client) Close() {
	close(c.closeCh)
}

// NewClient is constructor for Client structure. It will resolve IP address and dial the provided tcp address.
func NewClient(addr string) (*Client, error) {
	tcpAddr, wrongAddr := net.ResolveTCPAddr("tcp4", addr)
	if wrongAddr != nil {
		return nil, wrongAddr
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}

	c := &Client{
		RequestHandler{},
		ResponseHandler{},
		conn,
		nil,
		make(chan bool),
		make(map[uint64]chan *Message),
		make(map[uint64]chan *Message),
	}
	c.Start()

	_, err = c.Topology()
	if err != nil {
		log.Printf("TopologyRequest err: %+v\n", err)
		return nil, err
	}

	return c, nil
}
