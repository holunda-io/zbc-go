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
	"math/rand"
	"sync"
)

var (
	errTimeout             = errors.New("Request timeout")
	errSocketWrite         = errors.New("Tried to write more bytes to socket")
	errTopicLeaderNotFound = errors.New("Topic leader not found")
)

// Client for Zeebe broker with support for clustered deployment.
type Client struct {
	*sync.Mutex

	requestHandler
	responseHandler

	Connection         net.Conn
	Cluster            *zbmsgpack.ClusterTopology
	closeCh            chan bool
	transactions       map[uint64]chan *Message
	//topicSubscriptions map[uint64]chan *TopicEvent
	//taskSubscriptions  map[uint64]chan *TaskEvent
	subscriptions map[uint64]chan *SubscriptionEvent
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
				continue
			}

			if headers.IsSingleMessage() && message != nil {
				event := (*message.SbeMessage).(*zbsbe.SubscribedEvent)
				c.subscriptions[event.SubscriberKey] <- &SubscriptionEvent{Task: message.Task(), Event: event}
			}
		}
	}
}

// responder implements synchronous way of sending ExecuteCommandRequest and waiting for ExecuteCommandResponse.
func (c *Client) responder(message *Message) (*Message, error) {
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

	return MessageRetry(func() (*Message, error) {
		return c.responder(message)
	})
}

// CreateWorkflow will deploy process to the broker.
func (c *Client) CreateWorkflow(topic string, resourceType string, resource []byte) (*Message, error) {
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

	return MessageRetry(func() (*Message, error) {
		return c.responder(c.newWorkflowRequest(commandRequest, &deployment))
	})
}

// CreateWorkflowInstance will create new workflow instance on the broker.
func (c *Client) CreateWorkflowInstance(topic string, m *zbmsgpack.WorkflowInstance) (*Message, error) {
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

	return MessageRetry(func() (*Message, error) {
		return c.responder(c.createWorkflowInstanceRequest(commandRequest, m))
	})
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

	d := response.TaskSubscription()
	subscriberKey := d.SubscriberKey
	c.subscriptions[subscriberKey] = subscriptionCh

	return subscriptionCh, d, nil
}

// CompleteTask will notify broker about finished task.
func (c *Client) CompleteTask(task *SubscriptionEvent) (*zbmsgpack.Task, error) {
	msg, err := MessageRetry(func() (*Message, error) { return c.responder(c.completeTaskRequest(task)) })
	return msg.Task(), err
}

//  IncreaseTaskSubscriptionCredits will increase the current credits of the task subscription.
func (c *Client) IncreaseTaskSubscriptionCredits(task *zbmsgpack.TaskSubscription) (*zbmsgpack.TaskSubscription, error) {
	msg := c.increaseTaskSubscriptionCreditsRequest(task)
	response, err := MessageRetry(func() (*Message, error) { return c.responder(msg) })
	if err != nil {
		return nil, err
	}

	d := response.TaskSubscription()

	return d, nil
}

// CloseTaskSubscription will close currently active task subscription.
func (c *Client) CloseTaskSubscription(task *zbmsgpack.TaskSubscription) (*Message, error) {
	return MessageRetry(func() (*Message, error) { return c.responder(c.closeTaskSubscriptionRequest(task)) })
}

// TopicConsumer opens a subscription on topic and returns a channel where all the SubscribedEvents will arrive.
func (c *Client) TopicConsumer(topic, subName string, startPosition int64) (chan *SubscriptionEvent, *zbmsgpack.TopicSubscription, error) {
	partitionID, err := c.partitionID(topic)
	if err != nil {
		return nil, nil, err
	}

	topicSub := &zbmsgpack.TopicSubscription{
		StartPosition:    startPosition,
		Name:             subName,
		PrefetchCapacity: 0,
		ForceStart:       false,
		State:            TopicSubscriptionSubscribeState,
	}
	execCommandRequest := &zbsbe.ExecuteCommandRequest{
		PartitionId: partitionID,
		Position:    0,
		EventType:   zbsbe.EventType.SUBSCRIBER_EVENT,
		TopicName:   []byte(topic),
	}
	execCommandRequest.Key = execCommandRequest.KeyNullValue()

	subscriptionCh := make(chan *SubscriptionEvent)
	msg := c.openTopicSubscriptionRequest(execCommandRequest, topicSub)

	response, err := MessageRetry(func() (*Message, error) { return c.responder(msg) })
	if err != nil {
		return nil, nil, err
	}

	subscriberKey := (*response.SbeMessage).(*zbsbe.ExecuteCommandResponse).Key
	c.subscriptions[subscriberKey] = subscriptionCh

	return subscriptionCh, topicSub, nil
}

// TopologyRequest will retrieve latest cluster topology information.
func (c *Client) Topology() (*zbmsgpack.ClusterTopology, error) {
	resp, err := MessageRetry(func() (*Message, error) { return c.responder(c.topologyRequest()) })
	if err != nil {
		return nil, err
	}
	topology := c.newClusterTopologyResponse(resp)
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

	topic := &zbmsgpack.Topic{
		name,
		TopicCreate,
		partitionNum,
	}

	resp, err := MessageRetry(func() (*Message, error) { return c.responder(c.createTopicRequest(execCommandRequest, topic)) })
	if err != nil {
		return nil, err
	}

	return c.newCreateTopicResponse(resp), nil
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
		make(map[uint64]chan *Message),
		make(map[uint64]chan *SubscriptionEvent),
	}
	go c.receiver()
	go c.manageTopology()

	_, err = c.Topology()

	if err != nil {
		log.Printf("TopologyRequest err: %+v\n", err)
		return nil, err
	}

	return c, nil
}

// NewTask is constructor for Task object. Function signature denotes mandatory fields.
func NewTask(typeName string) *zbmsgpack.Task {
	return &zbmsgpack.Task{
		State:   TaskCreate,
		Type:    typeName,
		Retries: 3,
	}
}

func NewWorkflowInstance(bpmnProcessId string, version int, payload []uint8) *zbmsgpack.WorkflowInstance {
	return &zbmsgpack.WorkflowInstance{
		State:         CreateWorkflowInstance,
		BPMNProcessID: bpmnProcessId,
		Version:       version,
		Payload:       payload,
	}
}
