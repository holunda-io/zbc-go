package zbc

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/jsam/zbc-go/zbc/sbe"
	"log"
	"net"
	"time"
)

const REQUEST_TIMEOUT = 5

var (
	TimeoutError     = errors.New("Request timeout.")
	SocketWriteError = errors.New("Tried to write more bytes to socket.")
)

type Client struct {
	conn          net.Conn
	transactions  map[uint64]chan *Message
	subscriptions map[uint64]chan *Message
}

func (c *Client) sender(message *Message) error {
	writer := NewMessageWriter(message)
	byteBuff := &bytes.Buffer{}
	writer.Write(byteBuff)

	n, err := c.conn.Write(byteBuff.Bytes())
	if err != nil {
		return err
	}

	if n != len(byteBuff.Bytes()) {
		return SocketWriteError
	}
	return nil
}

func (c *Client) receiver() {
	for {
		buffer := bufio.NewReader(c.conn)
		r := NewMessageReader(buffer)
		headers, tail, err := r.ReadHeaders()

		if err != nil {
			log.Printf("[R] Error %+#v\n", err)
			continue
		}
		message, err := r.ParseMessage(headers, tail)

		if err != nil && !headers.IsSingleMessage() {
			// TODO: Maybe we should panic here?
			delete(c.transactions, headers.RequestResponseHeader.RequestId)
			continue
		}

		if !headers.IsSingleMessage() && message != nil {
			c.transactions[headers.RequestResponseHeader.RequestId] <- message
			continue
		}

		if err != nil && headers.IsSingleMessage() {
			// TODO: close and delete subscription
			continue
		}

		if headers.IsSingleMessage() && message != nil {
			subscriberKey := (*message.SbeMessage).(*sbe.SubscribedEvent).SubscriberKey
			c.subscriptions[subscriberKey] <- message
			continue
		}

	}
}

// Responder implements synchronous way of sending ExecuteCommandRequest and waiting for ExecuteCommandResponse.
func (c *Client) Responder(message *Message) (*Message, error) {
	respCh := make(chan *Message)
	c.transactions[message.Headers.RequestResponseHeader.RequestId] = respCh

	if err := c.sender(message); err != nil {
		return nil, err
	}

	select {
	case resp := <-c.transactions[message.Headers.RequestResponseHeader.RequestId]:
		delete(c.transactions, message.Headers.RequestResponseHeader.RequestId)
		return resp, nil
	case <-time.After(time.Second * REQUEST_TIMEOUT):
		delete(c.transactions, message.Headers.RequestResponseHeader.RequestId)
		return nil, TimeoutError
	}
}

// TaskConsumer opens a subscription on task and returns a channel where all the SubscribedEvents will arrive.
func (c *Client) TaskConsumer(ts *TaskSubscription) (chan *Message, error) {
	subscriptionCh := make(chan *Message, ts.Credits)
	msg := NewTaskSubscriptionMessage(ts)

	response, err := c.Responder(msg)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	c.subscriptions[(*response.Data)["subscriberKey"].(uint64)] = subscriptionCh

	return subscriptionCh, nil
}

func (c *Client) Connect() {
	go c.receiver()
}

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
		conn,
		make(map[uint64]chan *Message),
		make(map[uint64]chan *Message),
	}
	c.Connect()

	return c, nil
}
