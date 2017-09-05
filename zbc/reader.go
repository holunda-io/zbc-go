package zbc

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"io"

	"github.com/zeebe-io/zbc-go/zbc/zbprotocol"
	"github.com/zeebe-io/zbc-go/zbc/zbsbe"
)

var (
	errFrameHeaderRead    = errors.New("Cannot read bytes for frame header")
	errFrameHeaderDecode  = errors.New("Cannot decode bytes into frame header")
	errProtocolIDNotFound = errors.New("ProtocolID not found")
)

// MessageReader is builder which will read byte array and construct Message with all their parts.
type MessageReader struct {
	io.Reader
}

func (mr *MessageReader) readNext(n uint32) ([]byte, error) {
	data := make([]byte, n)
	bytesRead, err := mr.Read(data)

	switch {
	case bytesRead > 0:
		return data, nil

	case err != nil:
		//log.Println("Found error .... backing off.")
		//log.Println(err)
		return nil, err
	}

	return nil, nil
}

func (mr *MessageReader) readFrameHeader(data io.Reader) (*zbprotocol.FrameHeader, error) {
	var frameHeader zbprotocol.FrameHeader
	if frameHeader.Decode(data, binary.LittleEndian, 0) != nil {
		return nil, errFrameHeaderDecode
	}
	return &frameHeader, nil
}

func (mr *MessageReader) readTransportHeader(data io.Reader) (*zbprotocol.TransportHeader, error) {
	var transport zbprotocol.TransportHeader
	err := transport.Decode(data, binary.LittleEndian, 0)
	if err != nil {
		return nil, err
	}
	if transport.ProtocolID == zbprotocol.RequestResponse || transport.ProtocolID == zbprotocol.FullDuplexSingleMessage {
		return &transport, nil
	}

	return nil, errProtocolIDNotFound
}

func (mr *MessageReader) readRequestResponseHeader(data io.Reader) (*zbprotocol.RequestResponseHeader, error) {
	var requestResponse zbprotocol.RequestResponseHeader
	err := requestResponse.Decode(data, binary.LittleEndian, 0)
	if err != nil {
		return nil, err
	}

	return &requestResponse, nil
}

func (mr *MessageReader) readSbeMessageHeader(data io.Reader) (*zbsbe.MessageHeader, error) {
	var sbeMessageHeader zbsbe.MessageHeader
	err := sbeMessageHeader.Decode(data, binary.LittleEndian, 0)
	if err != nil {
		return nil, err
	}
	return &sbeMessageHeader, nil
}

// ReadHeaders will read entire message and interpret all headers. It will return pointer to headers object and tail of the message as byte array.
func (mr *MessageReader) ReadHeaders() (*Headers, *[]byte, error) {
	var header Headers

	headerByte, err := mr.readNext(FrameHeaderSize)
	switch {
	case err == errTimeout:
		return nil, nil, errTimeout

	case err != nil:
		return nil, nil, errFrameHeaderRead
	}

	frameHeaderReader := bytes.NewReader(headerByte)
	frameHeader, err := mr.readFrameHeader(frameHeaderReader)
	if err != nil {
		return nil, nil, err
	}
	header.SetFrameHeader(frameHeader)

	message, err := mr.readNext(frameHeader.Length)
	if err != nil {
		return nil, nil, err
	}

	transportReader := bytes.NewReader(message[:TransportHeaderSize])
	transport, err := mr.readTransportHeader(transportReader)
	if err != nil {
		return nil, nil, err
	}
	header.SetTransportHeader(transport)

	sbeIndex := TransportHeaderSize
	switch transport.ProtocolID {
	case zbprotocol.RequestResponse:
		reqRespReader := bytes.NewReader(message[TransportHeaderSize : TransportHeaderSize+RequestResponseHeaderSize])
		requestResponse, errHeader := mr.readRequestResponseHeader(reqRespReader)
		if errHeader != nil {
			return nil, nil, err
		}
		header.SetRequestResponseHeader(requestResponse)
		sbeIndex = TransportHeaderSize + RequestResponseHeaderSize
		break

	case zbprotocol.FullDuplexSingleMessage:
		header.SetRequestResponseHeader(nil)
		break
	}

	sbeHeaderReader := bytes.NewReader(message[sbeIndex : sbeIndex+SBEMessageHeaderSize])
	sbeMessageHeader, err := mr.readSbeMessageHeader(sbeHeaderReader)
	if err != nil {
		return nil, nil, err
	}
	header.SetSbeMessageHeader(sbeMessageHeader)

	body := message[sbeIndex+8:]
	return &header, &body, nil
}

func (mr *MessageReader) decodeCmdRequest(reader *bytes.Reader, header *zbsbe.MessageHeader) (*zbsbe.ExecuteCommandRequest, error) {
	var commandRequest zbsbe.ExecuteCommandRequest

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

func (mr *MessageReader) decodeCmdResponse(reader *bytes.Reader, header *zbsbe.MessageHeader) (*zbsbe.ExecuteCommandResponse, error) {
	var commandResponse zbsbe.ExecuteCommandResponse
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

func (mr *MessageReader) decodeCtlResponse(reader *bytes.Reader, header *zbsbe.MessageHeader) (*zbsbe.ControlMessageResponse, error) {
	var controlResponse zbsbe.ControlMessageResponse
	err := controlResponse.Decode(reader, binary.LittleEndian, header.Version, header.BlockLength, true)
	if err != nil {
		return nil, err
	}
	return &controlResponse, nil
}

func (mr *MessageReader) decodeSubEvent(reader *bytes.Reader, header *zbsbe.MessageHeader) (*zbsbe.SubscribedEvent, error) {
	var subEvent zbsbe.SubscribedEvent
	err := subEvent.Decode(reader, binary.LittleEndian, header.Version, header.BlockLength, true)
	if err != nil {
		return nil, err
	}
	return &subEvent, nil
}

//func (mr *MessageReader) parseMessagePack(data *[]byte) (*map[string]interface{}, error) {
//	var item map[string]interface{}
//	err := zbmsgpack.Unmarshal(*data, &item)
//
//	if err != nil {
//		return nil, err
//	}
//	return &item, nil
//}
//
//
//func (mr *MessageReader) parseMessagePackToType(data *[]byte, t interface{}) (interface{}, error) {
//	err := zbmsgpack.Unmarshal(*data, &t)
//
//	if err != nil {
//		return nil, err
//	}
//	return &t, nil
//}

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
		msg.SetData([]byte(commandRequest.Command))

		break

	case templateIDExecuteCommandResponse: // Read response from the socket.
		commandResponse, err := mr.decodeCmdResponse(reader, headers.SbeMessageHeader)
		if err != nil {
			return nil, err
		}
		msg.SetSbeMessage(commandResponse)
		msg.SetData([]byte(commandResponse.Event))
		break

	case templateIDControlMessageResponse:
		ctlResponse, err := mr.decodeCtlResponse(reader, headers.SbeMessageHeader)
		if err != nil {
			return nil, err
		}
		msg.SetSbeMessage(ctlResponse)
		msg.SetData([]byte(ctlResponse.Data))

		break

	case templateIDSubscriptionEvent:
		subscribedEvent, err := mr.decodeSubEvent(reader, headers.SbeMessageHeader)
		if err != nil {
			return nil, err
		}
		msg.SetSbeMessage(subscribedEvent)
		msg.SetData([]byte(subscribedEvent.Event))

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
