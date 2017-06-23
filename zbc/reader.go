package zbc

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"io"

	"github.com/jsam/zbc-go/zbc/protocol"
	"github.com/jsam/zbc-go/zbc/sbe"
	"gopkg.in/vmihailenco/msgpack.v2"
)

var (
	errFrameHeaderRead    = errors.New("Cannot read bytes for frame header")
	errFrameHeaderDecode  = errors.New("Cannot decode bytes into frame header")
	errProtocolIDNotFound = errors.New("ProtocolId not found")
)

// MessageReader is builder which will read byte array and construct Message with all their parts.
type MessageReader struct {
	io.Reader
}

func (mr *MessageReader) readNext(n uint32) ([]byte, error) {
	buffer := make([]byte, n)

	numBytes, err := mr.Read(buffer)

	if uint32(numBytes) != n || err != nil {
		return nil, err
	}

	return buffer, nil
}

func (mr *MessageReader) readFrameHeader(data []byte) (*protocol.FrameHeader, error) {
	var frameHeader protocol.FrameHeader
	if frameHeader.Decode(bytes.NewReader(data[:12]), binary.LittleEndian, 0) != nil {
		return nil, errFrameHeaderDecode
	}
	return &frameHeader, nil
}

func (mr *MessageReader) readTransportHeader(data []byte) (*protocol.TransportHeader, error) {
	var transport protocol.TransportHeader
	err := transport.Decode(bytes.NewReader(data[:2]), binary.LittleEndian, 0)
	if err != nil {
		return nil, err
	}
	if transport.ProtocolID == protocol.RequestResponse || transport.ProtocolID == protocol.FullDuplexSingleMessage {
		return &transport, nil
	}

	return nil, errProtocolIDNotFound
}

func (mr *MessageReader) readRequestResponseHeader(data []byte) (*protocol.RequestResponseHeader, error) {
	var requestResponse protocol.RequestResponseHeader
	err := requestResponse.Decode(bytes.NewReader(data[:16]), binary.LittleEndian, 0)
	if err != nil {
		return nil, err
	}

	return &requestResponse, nil
}

func (mr *MessageReader) readSbeMessageHeader(data []byte) (*sbe.MessageHeader, error) {
	var sbeMessageHeader sbe.MessageHeader
	err := sbeMessageHeader.Decode(bytes.NewReader(data[:8]), binary.LittleEndian, 0)
	if err != nil {
		return nil, err
	}
	return &sbeMessageHeader, nil
}

// ReadHeaders will read entire message and interpret all headers. It will return pointer to headers object and tail of the message as byte array.
func (mr *MessageReader) ReadHeaders() (*Headers, *[]byte, error) {
	var header Headers

	headerByte, err := mr.readNext(12)
	if err != nil {
		return nil, nil, errFrameHeaderRead
	}

	frameHeader, err := mr.readFrameHeader(headerByte)
	if err != nil {
		return nil, nil, err
	}
	header.SetFrameHeader(frameHeader)

	message, err := mr.readNext(frameHeader.Length)
	if err != nil {
		return nil, nil, err
	}

	transport, err := mr.readTransportHeader(message[:2])
	if err != nil {
		return nil, nil, err
	}
	header.SetTransportHeader(transport)

	sbeIndex := 2
	switch transport.ProtocolID {
	case protocol.RequestResponse:
		requestResponse, errHeader := mr.readRequestResponseHeader(message[2:18])
		if errHeader != nil {
			return nil, nil, err
		}
		header.SetRequestResponseHeader(requestResponse)
		sbeIndex = 18
		break

	case protocol.FullDuplexSingleMessage:
		header.SetRequestResponseHeader(nil)
		break
	}

	sbeMessageHeader, err := mr.readSbeMessageHeader(message[sbeIndex : sbeIndex+8])
	if err != nil {
		return nil, nil, err
	}
	header.SetSbeMessageHeader(sbeMessageHeader)

	// TODO: make client less reliable on garbage -> this should align the reader for the next message
	mr.align()

	body := message[sbeIndex+8:]
	return &header, &body, nil
}

func (mr *MessageReader) align() {
	// TODO:
}

func (mr *MessageReader) decodeCmdRequest(reader *bytes.Reader, header *sbe.MessageHeader) (*sbe.ExecuteCommandRequest, error) {
	var commandRequest sbe.ExecuteCommandRequest
	err := commandRequest.Decode(reader,
		binary.LittleEndian,
		header.Version,
		header.BlockLength,
		true)
	if err != nil {
		return nil, err
	}
	return &commandRequest, nil
}

func (mr *MessageReader) decodeCmdResponse(reader *bytes.Reader, header *sbe.MessageHeader) (*sbe.ExecuteCommandResponse, error) {
	var commandResponse sbe.ExecuteCommandResponse
	err := commandResponse.Decode(reader,
		binary.LittleEndian,
		header.Version,
		header.BlockLength,
		true)
	if err != nil {
		return nil, err
	}
	return &commandResponse, nil
}

func (mr *MessageReader) decodeCtlResponse(reader *bytes.Reader, header *sbe.MessageHeader) (*sbe.ControlMessageResponse, error) {
	var controlResponse sbe.ControlMessageResponse
	err := controlResponse.Decode(reader, binary.LittleEndian, header.Version, header.BlockLength, true)
	if err != nil {
		return nil, err
	}
	return &controlResponse, nil
}

func (mr *MessageReader) decodeSubEvent(reader *bytes.Reader, header *sbe.MessageHeader) (*sbe.SubscribedEvent, error) {
	var subEvent sbe.SubscribedEvent
	err := subEvent.Decode(reader, binary.LittleEndian, header.Version, header.BlockLength, true)
	if err != nil {
		return nil, err
	}
	return &subEvent, nil
}

func (mr *MessageReader) parseMessagePack(data *[]byte) (*map[string]interface{}, error) {
	var item map[string]interface{}
	err := msgpack.Unmarshal(*data, &item)

	if err != nil {
		return nil, err
	}
	return &item, nil
}

// ParseMessage will take the headers and tail and construct Message.
func (mr *MessageReader) ParseMessage(headers *Headers, message *[]byte) (*Message, error) {
	var msg Message
	msg.SetHeaders(headers)
	reader := bytes.NewReader(*message)

	switch headers.SbeMessageHeader.TemplateId {

	case templateIDExecuteCommandRequest: // Testing purposes.
		commandRequest, err := mr.decodeCmdRequest(reader, headers.SbeMessageHeader)
		if err != nil {
			return nil, err
		}
		msg.SetSbeMessage(commandRequest)

		msgPackData, err := mr.parseMessagePack(&commandRequest.Command)
		if err != nil {
			return nil, err
		}

		msg.SetData(msgPackData)
		break

	case templateIDExecuteCommandResponse: // Read response from the socket.
		commandResponse, err := mr.decodeCmdResponse(reader, headers.SbeMessageHeader)
		if err != nil {
			return nil, err
		}
		msg.SetSbeMessage(commandResponse)

		msgPackData, err := mr.parseMessagePack(&commandResponse.Event)
		if err != nil {
			return nil, err
		}
		msg.SetData(msgPackData)
		break

	case templateIDControlMessageResponse:
		ctlResponse, err := mr.decodeCtlResponse(reader, headers.SbeMessageHeader)
		if err != nil {
			return nil, err
		}
		msg.SetSbeMessage(ctlResponse)
		msgPackData, err := mr.parseMessagePack(&ctlResponse.Data)
		if err != nil {
			return nil, err
		}
		msg.SetData(msgPackData)
		break

	case templateIDSubscriptionEvent:
		subscribedEvent, err := mr.decodeSubEvent(reader, headers.SbeMessageHeader)
		if err != nil {
			return nil, err
		}

		msg.SetSbeMessage(subscribedEvent)
		msgPackData, err := mr.parseMessagePack(&subscribedEvent.Event)
		if err != nil {
			return nil, err
		}
		msg.SetData(msgPackData)
		break
	}
	return &msg, nil
}

// NewMessageReader is constructor for MessageReader builder.
func NewMessageReader(rd *bufio.Reader) *MessageReader {
	return &MessageReader{
		rd,
	}
}
