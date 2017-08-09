package zbc

import (
	"encoding/binary"
	"io"

	"github.com/zeebe-io/zbc-go/zbc/protocol"
	"github.com/zeebe-io/zbc-go/zbc/sbe"
)

const (
	templateIDExecuteCommandRequest  = 20
	templateIDExecuteCommandResponse = 21
	templateIDControlMessageResponse = 11
	templateIDSubscriptionEvent      = 30
)

const (
	FrameHeaderSize = 12
	TransportHeaderSize = 2
	RequestResponseHeaderSize = 8
	SBEMessageHeaderSize = 8

	TotalHeaderSizeNoFrame = 18
	LengthFieldSize = 2
)

// Headers is aggregator for all headers. It holds pointer to every layer. If RequestResponseHeader is nil, then IsSingleMessage will always return true.
type Headers struct {
	FrameHeader           *protocol.FrameHeader
	TransportHeader       *protocol.TransportHeader
	RequestResponseHeader *protocol.RequestResponseHeader // If this is nil then struct is equal to SingleMessage
	SbeMessageHeader      *sbe.MessageHeader
}

// SetFrameHeader is a setter for FrameHeader.
func (h *Headers) SetFrameHeader(header *protocol.FrameHeader) {
	h.FrameHeader = header
}

// SetTransportHeader is a setter for TransportHeader.
func (h *Headers) SetTransportHeader(header *protocol.TransportHeader) {
	h.TransportHeader = header
}

// SetRequestResponseHeader is a setting for RequestResponseHeader.
func (h *Headers) SetRequestResponseHeader(header *protocol.RequestResponseHeader) {
	h.RequestResponseHeader = header
}

// IsSingleMessage is helper to determine which model of communication we use.
func (h *Headers) IsSingleMessage() bool {
	return h.RequestResponseHeader == nil
}

// SetSbeMessageHeader is a setter for SBEMessageHeader.
func (h *Headers) SetSbeMessageHeader(header *sbe.MessageHeader) {
	h.SbeMessageHeader = header
}

// SBE interface is apstraction over all SBE Messages.
type SBE interface {
	Encode(writer io.Writer, order binary.ByteOrder, doRangeCheck bool) error
	Decode(reader io.Reader, order binary.ByteOrder, actingVersion uint16, blockLength uint16, doRangeCheck bool) error
}

// Message is zeebe message which will contain pointers to all parts of the message. Data is Message Pack layer.
type Message struct {
	Headers    *Headers
	SbeMessage *SBE
	Data       *map[string]interface{}
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
func (m *Message) SetData(data *map[string]interface{}) {
	m.Data = data
}
