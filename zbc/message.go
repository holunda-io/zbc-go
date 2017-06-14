package zbc

import (
	"encoding/binary"
	"github.com/jsam/zbc-go/zbc/protocol"
	"github.com/jsam/zbc-go/zbc/sbe"
	"io"
)

const (
	SBE_ExecuteCommandRequest_TemplateId  = 20
	SBE_ExecuteCommandResponse_TemplateId = 21
)

type Headers struct {
	FrameHeader           *protocol.FrameHeader
	TransportHeader       *protocol.TransportHeader
	RequestResponseHeader *protocol.RequestResponseHeader // If this is nil then struct is equal to SingleMessage
	SbeMessageHeader      *sbe.MessageHeader
}

func (h *Headers) SetFrameHeader(header *protocol.FrameHeader) {
	h.FrameHeader = header
}

func (h *Headers) SetTransportHeader(header *protocol.TransportHeader) {
	h.TransportHeader = header
}

func (h *Headers) SetRequestResponseHeader(header *protocol.RequestResponseHeader) {
	h.RequestResponseHeader = header
}

func (h *Headers) IsSingleMessage() bool {
	return h.RequestResponseHeader == nil
}

func (h *Headers) SetSbeMessageHeader(header *sbe.MessageHeader) {
	h.SbeMessageHeader = header
}

type SBE interface {
	Encode(writer io.Writer, order binary.ByteOrder, doRangeCheck bool) error
	Decode(reader io.Reader, order binary.ByteOrder, actingVersion uint16, blockLength uint16, doRangeCheck bool) error
}

type Message struct {
	Headers    *Headers
	SbeMessage *SBE
	Data       *map[string]interface{}
}

func (m *Message) SetHeaders(headers *Headers) {
	m.Headers = headers
}

func (m *Message) SetSbeMessage(data SBE) {
	m.SbeMessage = &data
}

func (m *Message) SetData(data *map[string]interface{}) {
	m.Data = data
}
