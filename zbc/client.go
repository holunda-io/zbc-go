package zbc

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"net"
	"path/filepath"
	"time"

	"math/rand"
	"sync"

	"github.com/zeebe-io/zbc-go/zbc/zbmsgpack"
	"github.com/zeebe-io/zbc-go/zbc/zbsbe"
	"github.com/vmihailenco/msgpack"
	"fmt"
)

var (
	errTimeout             = errors.New("request timeout")
	errSocketWrite         = errors.New("tried to write more bytes to socket")
	errTopicLeaderNotFound = errors.New("topic leader not found")
	errResourceNotFound    = errors.New("resource not found")
)

// Client for Zeebe broker with support for clustered deployment.
type Client struct {
	*sync.Mutex

	requestHandler
	responseHandler

	Connection    net.Conn
	Cluster       *zbmsgpack.ClusterTopology
	closeCh       chan bool
	transactions  SafeMap //map[uint64]chan *Message
	subscriptions SafeMap //map[uint64]chan *SubscriptionEvent
}

func (c *Client) partitionID(topic string) (uint16, error) {
	c.Lock()
	leaders, ok := c.Cluster.TopicLeaders[topic]
	c.Unlock()

	if !ok {
		c.Topology()

		c.Lock()
		leaders, ok = c.Cluster.TopicLeaders[topic]
		c.Unlock()

		if !ok {
			return 0, errTopicLeaderNotFound
		}
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	index := rnd.Intn(len(leaders))
	return leaders[index].PartitionID, nil
}

func (c *Client) sender(message *Message) error {
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

	reader := NewMessageReader(c.Connection)

	for {
		select {
		case <-c.closeCh:
			c.Connection.Close()
			return

		default:
			headers, tail, err := reader.readHeaders()
			if err != nil {
				continue
			}
			message, err := reader.parseMessage(headers, tail)

			if err != nil && !headers.IsSingleMessage() {
				c.transactions.Remove(fmt.Sprintf("%d", headers.RequestResponseHeader.RequestID))
				continue
			}

			if !headers.IsSingleMessage() && message != nil && len(message.Data) > 0 {
				transactionKey := fmt.Sprintf("%d", headers.RequestResponseHeader.RequestID)
				if ch, ok := c.transactions.Get(transactionKey); ok {
					requestCh := ch.(chan *Message)
					requestCh <- message
				}

				continue
			}

			if err != nil && headers.IsSingleMessage() {
				continue
			}

			if headers.IsSingleMessage() && message != nil && len(message.Data) > 0 {
				event := (*message.SbeMessage).(*zbsbe.SubscribedEvent)
				subscriberKey := fmt.Sprintf("%d", event.SubscriberKey)
				if ch, ok := c.subscriptions.Get(subscriberKey); ok {
					requestCh := ch.(chan *SubscriptionEvent)
					requestCh <- &SubscriptionEvent{Task: c.unmarshalTask(message), Event: event}
				}
			}
		}
	}
	fmt.Println("omg receiver died. help me!!!!!")
}

// responder implements synchronous way of sending ExecuteCommandRequest and waiting for ExecuteCommandResponse.
func (c *Client) responder(message *Message) (*Message, error) {
	respCh := make(chan *Message, 10)
	c.transactions.Set(fmt.Sprintf("%d", message.Headers.RequestResponseHeader.RequestID), respCh)
	if err := c.sender(message); err != nil {
		return nil, err
	}

	select {
	case resp := <-respCh:
		c.transactions.Remove(fmt.Sprintf("%d", message.Headers.RequestResponseHeader.RequestID))
		return resp, nil
	case <-time.After(time.Second * RequestTimeout):
		c.transactions.Remove(fmt.Sprintf("%d", message.Headers.RequestResponseHeader.RequestID))
		return nil, errTimeout
	}
}

// CreateTask will create new task on specified topic.
func (c *Client) CreateTask(topic string, m *zbmsgpack.Task) (*zbmsgpack.Task, error) {
	partitionID, err := c.partitionID(topic)
	if err != nil {
		return nil, err
	}

	commandRequest := &zbsbe.ExecuteCommandRequest{
		PartitionId: partitionID,
		Position:    0,
		TopicName:   []uint8(topic),
		Command:     []uint8{},
	}
	commandRequest.Key = commandRequest.KeyNullValue()
	message := c.createTaskRequest(commandRequest, m)

	msg, err := MessageRetry(func() (*Message, error) {
		return c.responder(message)
	})
	return c.unmarshalTask(msg), err
}

// CreateWorkflow will deploy process to the broker.
func (c *Client) CreateWorkflow(topic string, resourceType string, resource []byte) (*zbmsgpack.Workflow, error) {
	partitionID, err := c.partitionID(topic)
	if err != nil {
		return nil, err
	}

	deployment := zbmsgpack.Workflow{
		State:        CreateDeployment,
		ResourceType: resourceType,
		Resource:     resource,
	}
	commandRequest := &zbsbe.ExecuteCommandRequest{
		PartitionId: partitionID,
		Position:    0,
		TopicName:   []uint8(topic),
		Command:     []uint8{},
	}
	commandRequest.Key = commandRequest.KeyNullValue()

	msg, err := MessageRetry(func() (*Message, error) {
		return c.responder(c.newWorkflowRequest(commandRequest, &deployment))
	})
	return c.unmarshalWorkflow(msg), err
}

func (c *Client) CreateWorkflowFromFile(topic, resourceType, path string) (*zbmsgpack.Workflow, error) {
	if len(path) == 0 {
		return nil, errResourceNotFound
	}

	filename, _ := filepath.Abs(path)
	definition, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errResourceNotFound
	}
	return c.CreateWorkflow(topic, resourceType, definition)
}

// CreateWorkflowInstance will create new workflow instance on the broker.
func (c *Client) CreateWorkflowInstance(topic string, m *zbmsgpack.WorkflowInstance) (*zbmsgpack.WorkflowInstance, error) {
	partitionID, err := c.partitionID(topic)
	if err != nil {
		return nil, err
	}

	commandRequest := &zbsbe.ExecuteCommandRequest{
		PartitionId: partitionID,
		Position:    0,
		TopicName:   []uint8(topic),
		Command:     []uint8{},
	}

	commandRequest.Key = commandRequest.KeyNullValue()

	msg, err := MessageRetry(func() (*Message, error) {
		return c.responder(c.createWorkflowInstanceRequest(commandRequest, m))
	})
	return c.unmarshalWorkflowInstance(msg), err
}

// TaskConsumer opens a subscription on task and returns a channel where all the SubscribedEvents will arrive.
func (c *Client) TaskConsumer(topic, lockOwner, taskType string) (chan *SubscriptionEvent, *zbmsgpack.TaskSubscription, error) {
	partitionID, err := c.partitionID(topic)
	if err != nil {
		return nil, nil, err
	}

	taskSub := &zbmsgpack.TaskSubscription{
		TopicName:     topic,
		PartitionID:   partitionID,
		Credits:       32,
		LockDuration:  300000,
		LockOwner:     lockOwner,
		SubscriberKey: 0,
		TaskType:      taskType,
	}

	subscriptionCh := make(chan *SubscriptionEvent, taskSub.Credits)

	msg := c.openTaskSubscriptionRequest(taskSub)
	response, err := MessageRetry(func() (*Message, error) { return c.responder(msg) })

	if err != nil {
		return nil, nil, err
	}

	d := c.unmarshalTaskSubscription(response)
	c.subscriptions.Set(fmt.Sprintf("%d", d.SubscriberKey), subscriptionCh)

	return subscriptionCh, d, nil
}

// CompleteTask will notify broker about finished task.
func (c *Client) CompleteTask(task *SubscriptionEvent) (*zbmsgpack.Task, error) {
	msg, err := MessageRetry(func() (*Message, error) { return c.responder(c.completeTaskRequest(task)) })
	return c.unmarshalTask(msg), err
}

//  IncreaseTaskSubscriptionCredits will increase the current credits of the task subscription.
func (c *Client) IncreaseTaskSubscriptionCredits(task *zbmsgpack.TaskSubscription) (*zbmsgpack.TaskSubscription, error) {
	msg := c.increaseTaskSubscriptionCreditsRequest(task)
	response, err := MessageRetry(func() (*Message, error) { return c.responder(msg) })
	if err != nil {
		return nil, err
	}

	d := c.unmarshalTaskSubscription(response)
	return d, nil
}

// CloseTaskSubscription will close currently active task subscription.
func (c *Client) CloseTaskSubscription(task *zbmsgpack.TaskSubscription) (*Message, error) {
	return MessageRetry(func() (*Message, error) { return c.responder(c.closeTaskSubscriptionRequest(task)) })
}

// CloseTopicSubscription will close currently active topic subscription.
func (c *Client) CloseTopicSubscription(task *zbmsgpack.TopicSubscription) (*Message, error) {
	return MessageRetry(func() (*Message, error) {
		msg, err := c.responder(c.closeTopicSubscriptionRequest(task))
		return msg, err
	})
}

// TopicConsumer opens a subscription on topic and returns a channel where all the SubscribedEvents will arrive.
func (c *Client) TopicConsumer(topic, subName string, startPosition int64) (chan *SubscriptionEvent, *zbmsgpack.TopicSubscription, error) {
	partitionID, err := c.partitionID(topic)
	if err != nil {
		return nil, nil, err
	}

	topicSub := &zbmsgpack.OpenTopicSubscription{
		StartPosition:    startPosition,
		Name:             subName,
		PrefetchCapacity: 0,
		ForceStart:       true,
		State:            TopicSubscriptionSubscribeState,
	}
	execCommandRequest := &zbsbe.ExecuteCommandRequest{
		PartitionId: partitionID,
		Position:    0,
		EventType:   zbsbe.EventType.SUBSCRIBER_EVENT,
		TopicName:   []byte(topic),
	}
	execCommandRequest.Key = execCommandRequest.KeyNullValue()

	subscriptionCh := make(chan *SubscriptionEvent, 1000)
	msg := c.openTopicSubscriptionRequest(execCommandRequest, topicSub)

	response, err := MessageRetry(func() (*Message, error) { return c.responder(msg) })
	if err != nil {
		return nil, nil, err
	}

	cmdResponse := (*response.SbeMessage).(*zbsbe.ExecuteCommandResponse)
	subscriberKey := cmdResponse.Key
	c.subscriptions.Set(fmt.Sprintf("%d", subscriberKey), subscriptionCh)

	subscriptionInfo := &zbmsgpack.TopicSubscription{
		TopicName:     topic,
		PartitionID:   partitionID,
		SubscriberKey: subscriberKey,
	}

	return subscriptionCh, subscriptionInfo, nil
}

// TopologyRequest will retrieve latest cluster topology information.
func (c *Client) Topology() (*zbmsgpack.ClusterTopology, error) {
	resp, err := MessageRetry(func() (*Message, error) { return c.responder(c.topologyRequest()) })
	if err != nil {
		return nil, err
	}
	topology := c.unmarshalTopology(resp)
	c.Lock()
	c.Cluster = topology
	c.Unlock()
	return topology, nil
}

func (c *Client) manageTopology() {
	for {
		select {
		case <-time.After(TopologyRefreshInterval * time.Second):
			if time.Since(c.Cluster.UpdatedAt) > TopologyRefreshInterval*time.Second {
				c.Topology()
			}

			break
		}
	}
}

func (c *Client) CreateTopic(name string, partitionNum int) (*zbmsgpack.Topic, error) {
	execCommandRequest := &zbsbe.ExecuteCommandRequest{
		PartitionId: 0,
		Position:    0,
		EventType:   zbsbe.EventType.TOPIC_EVENT,
		TopicName:   []byte(SystemTopic),
	}
	execCommandRequest.Key = execCommandRequest.KeyNullValue()

	topic := zbmsgpack.NewTopic(name, TopicCreate, partitionNum)
	resp, err := MessageRetry(func() (*Message, error) {
		return c.responder(c.createTopicRequest(execCommandRequest, topic))
	})

	if err != nil {
		return nil, err
	}

	return c.unmarshalTopic(resp), nil
}

func (c *Client) UnmarshalFromFile(path string) (*Message, error) {
	data, fsErr := ioutil.ReadFile(path)

	if fsErr != nil {
		return nil, fsErr
	}

	messageReader := MessageReader{nil}
	return messageReader.readMessage(data)

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
		&sync.Mutex{},
		requestHandler{},
		responseHandler{},
		conn,
		nil,
		make(chan bool),
		New(),
		New(),
	}
	go c.receiver()
	//go c.manageTopology()

	_, err = c.Topology()

	if err != nil {
		log.Printf("TopologyRequest err: %+v\n", err)
		return nil, err
	}

	return c, nil
}

// NewTask is constructor for Task object. Function signature denotes mandatory fields.
func NewTask(typeName, lockOwner string) *zbmsgpack.Task {
	return &zbmsgpack.Task{
		State:        TaskCreate,
		Headers:      make(map[string]interface{}),
		CustomHeader: make(map[string]interface{}),

		Type:      typeName,
		LockOwner: lockOwner,
		Retries:   3,
	}
}

func NewWorkflowInstance(bpmnProcessId string, version int, payload map[string]interface{}) *zbmsgpack.WorkflowInstance {
	b, err := msgpack.Marshal(payload)
	if err != nil {
		return nil
	}
	return &zbmsgpack.WorkflowInstance{
		State:         CreateWorkflowInstance,
		BPMNProcessID: bpmnProcessId,
		Version:       version,
		Payload:       b,
	}
}
