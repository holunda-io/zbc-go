package zbc

import (
	"github.com/jsam/zbc-go/zbc/protocol"
	"github.com/jsam/zbc-go/zbc/sbe"
	"gopkg.in/vmihailenco/msgpack.v2"
)

type CreateTask struct {
	EventType string                 `yaml:"eventType" msgpack:"eventType"`
	Headers   map[string]interface{} `yaml:"headers" msgpack:"headers"`
	Payload   []uint8                `yaml:"payload" msgpack:"payload"`
	Retries   int                    `yaml:"retries" msgpack:"retries"`
	Type      string                 `yaml:"type" msgpack:"type"`
}

func NewCreateTaskMessage(commandRequest *sbe.ExecuteCommandRequest, createTask *CreateTask) *Message {
	var msg Message

	b, err := msgpack.Marshal(createTask)
	if err != nil {
		return nil
	}
	commandRequest.Command = b
	msg.SetSbeMessage(commandRequest)

	// We add +2 to every variable length attribute since all variable length attributes will have 2 bytes in front
	// which will denote their size. Then we add 11 bytes which is size of non-variable length attributes of
	// ExecuteCommandRequest and 26 bytes which is for SbeMessageHeader, RequestResponse and Transport.
	length := uint32(2+len(commandRequest.TopicName)) + uint32(2+len(commandRequest.Command)) + 11 + 26

	var headers Headers
	headers.SetSbeMessageHeader(&sbe.MessageHeader{
		BlockLength: commandRequest.SbeBlockLength(),
		TemplateId:  commandRequest.SbeTemplateId(),
		SchemaId:    commandRequest.SbeSchemaId(),
		Version:     commandRequest.SbeSchemaVersion(),
	})

	headers.SetRequestResponseHeader(protocol.NewRequestResponseHeader())
	headers.SetTransportHeader(protocol.NewTransportHeader(protocol.RequestResponse))

	// Writer will set FrameHeader after serialization to byte array.
	headers.SetFrameHeader(protocol.NewFrameHeader(uint32(length), 0, 0, 0, 2))

	msg.SetHeaders(&headers)
	return &msg
}

func NewCommandRequestMessage(commandRequest *sbe.ExecuteCommandRequest, payload *map[string]interface{}) *Message {
	var msg Message
	msg.SetData(payload)

	b, err := msgpack.Marshal(payload)
	if err != nil {
		return nil
	}
	commandRequest.Command = b
	msg.SetSbeMessage(commandRequest)

	// We add +2 to every variable length attribute since all variable length attributes will have 2 bytes in front
	// which will denote their size. Then we add 11 bytes which is size of non-variable length attributes of
	// ExecuteCommandRequest and 26 bytes which is for SbeMessageHeader, RequestResponse and Transport.
	length := uint32(2+len(commandRequest.TopicName)) + uint32(2+len(commandRequest.Command)) + 11 + 26

	var headers Headers
	headers.SetSbeMessageHeader(&sbe.MessageHeader{
		BlockLength: commandRequest.SbeBlockLength(),
		TemplateId:  commandRequest.SbeTemplateId(),
		SchemaId:    commandRequest.SbeSchemaId(),
		Version:     commandRequest.SbeSchemaVersion(),
	})

	headers.SetRequestResponseHeader(protocol.NewRequestResponseHeader())
	headers.SetTransportHeader(protocol.NewTransportHeader(protocol.RequestResponse))

	// Writer will set FrameHeader after serialization to byte array.
	headers.SetFrameHeader(protocol.NewFrameHeader(uint32(length), 0, 0, 0, 2))

	msg.SetHeaders(&headers)
	return &msg
}
