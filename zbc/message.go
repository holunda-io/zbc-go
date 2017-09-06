package zbc

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/zeebe-io/zbc-go/zbc/zbmsgpack"
	"github.com/zeebe-io/zbc-go/zbc/zbprotocol"
	"github.com/zeebe-io/zbc-go/zbc/zbsbe"
	"gopkg.in/vmihailenco/msgpack.v2"
	"io"
)

// Headers is aggregator for all headers. It holds pointer to every layer. If RequestResponseHeader is nil, then IsSingleMessage will always return true.
type Headers struct {
	FrameHeader           *zbprotocol.FrameHeader
	TransportHeader       *zbprotocol.TransportHeader
	RequestResponseHeader *zbprotocol.RequestResponseHeader // If this is nil then struct is equal to SingleMessage
	SbeMessageHeader      *zbsbe.MessageHeader
}

// SetFrameHeader is a setter for FrameHeader.
func (h *Headers) SetFrameHeader(header *zbprotocol.FrameHeader) {
	h.FrameHeader = header
}

// SetTransportHeader is a setter for TransportHeader.
func (h *Headers) SetTransportHeader(header *zbprotocol.TransportHeader) {
	h.TransportHeader = header
}

// SetRequestResponseHeader is a setting for RequestResponseHeader.
func (h *Headers) SetRequestResponseHeader(header *zbprotocol.RequestResponseHeader) {
	h.RequestResponseHeader = header
}

// IsSingleMessage is helper to determine which model of communication we use.
func (h *Headers) IsSingleMessage() bool {
	return h.RequestResponseHeader == nil
}

// SetSbeMessageHeader is a setter for SBEMessageHeader.
func (h *Headers) SetSbeMessageHeader(header *zbsbe.MessageHeader) {
	h.SbeMessageHeader = header
}

// SBE interface is apstraction over all SBE Messages.
type SBE interface {
	Encode(writer io.Writer, order binary.ByteOrder, doRangeCheck bool) error
	Decode(reader io.Reader, order binary.ByteOrder, actingVersion uint16, blockLength uint16, doRangeCheck bool) error
}

// Message is Zeebe message which will contain pointers to all parts of the message. Data is Message Pack layer.
type Message struct {
	Headers    *Headers
	SbeMessage *SBE
	Data       []byte
}

// SetHeaders is a setter for Headers attribute.
func (m *Message) SetHeaders(headers *Headers) {
	m.Headers = headers
}

// SetSbeMessage is a setter for SBE attribute.
func (m *Message) SetSbeMessage(data SBE) {
	m.SbeMessage = &data
}

// SetData is a setter for unmarshaled message pack data.
func (m *Message) SetData(data []byte) {
	m.Data = data
}

func (m *Message) Task() *zbmsgpack.Task {
	var d zbmsgpack.Task
	err := msgpack.Unmarshal(m.Data, &d)
	if err != nil {
		return nil
	}
	if len(d.Type) > 0 {
		return &d
	}
	return nil
}

func (m *Message) TaskSubscription() *zbmsgpack.TaskSubscription {
	var d zbmsgpack.TaskSubscription
	err := msgpack.Unmarshal(m.Data, &d)
	if err != nil {
		return nil
	}
	if len(d.TopicName) > 0 {
		return &d
	}
	return nil
}

func (m *Message) String() string {

	if task := m.Task(); task != nil {
		b, err := json.MarshalIndent(task, "", "  ")
		if err != nil {
			return fmt.Sprintf("json marshaling failed\n")
		}
		return fmt.Sprintf("%+v", string(b))
	}

	if tasksub := m.TaskSubscription(); tasksub != nil {
		b, err := json.MarshalIndent(tasksub, "", "  ")
		if err != nil {
			return fmt.Sprintf("json marshaling failed\n")
		}
		return fmt.Sprintf("%+v", string(b))
	}

	// TODO: implement missing msgpack structs
	return "not found \n"
}

func (m *Message) ParseToMap() (*map[string]interface{}, error) {
	var item map[string]interface{}
	err := msgpack.Unmarshal(m.Data, &item)

	if err != nil {
		return nil, err
	}
	return &item, nil
}
