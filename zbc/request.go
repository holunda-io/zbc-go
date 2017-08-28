package zbc

import (
	"github.com/zeebe-io/zbc-go/zbc/zbmsgpack"
	"github.com/zeebe-io/zbc-go/zbc/zbprotocol"
	"github.com/zeebe-io/zbc-go/zbc/zbsbe"
	"gopkg.in/vmihailenco/msgpack.v2"
)

type RequestHandler struct{}

func (rf *RequestHandler) newCommandMessage(commandRequest *zbsbe.ExecuteCommandRequest, command interface{}) *Message {
	var msg Message

	b, err := msgpack.Marshal(command)
	if err != nil {
		return nil
	}
	commandRequest.Command = b
	msg.SetSbeMessage(commandRequest)

	// We add +2 to every variable length attribute since all variable length attributes will have 2 bytes in front
	// which will denote their size. Then we add 19 bytes which is size of non-variable length attributes of
	// ExecuteCommandRequest and 26 bytes which is for SbeMessageHeader, RequestResponse and Transport.
	length := uint32(LengthFieldSize+len(commandRequest.TopicName)) + uint32(LengthFieldSize+len(commandRequest.Command))
	length += uint32(commandRequest.SbeBlockLength()) + TotalHeaderSizeNoFrame

	var headers Headers
	headers.SetSbeMessageHeader(&zbsbe.MessageHeader{
		BlockLength: commandRequest.SbeBlockLength(),
		TemplateId:  commandRequest.SbeTemplateId(),
		SchemaId:    commandRequest.SbeSchemaId(),
		Version:     commandRequest.SbeSchemaVersion(),
	})

	headers.SetRequestResponseHeader(zbprotocol.NewRequestResponseHeader())
	headers.SetTransportHeader(zbprotocol.NewTransportHeader(zbprotocol.RequestResponse))

	// Writer will set FrameHeader after serialization to byte array.
	headers.SetFrameHeader(zbprotocol.NewFrameHeader(uint32(length), 0, 0, 0, 0))

	msg.SetHeaders(&headers)
	return &msg
}

func (rf *RequestHandler) newControlMessage(req *zbsbe.ControlMessageRequest, payload interface{}) *Message {
	var msg Message

	b, err := msgpack.Marshal(payload)
	if err != nil {
		return nil
	}
	req.Data = b
	msg.SetSbeMessage(req)

	length := uint32(LengthFieldSize + len(req.Data))
	length += uint32(req.SbeBlockLength()) + TotalHeaderSizeNoFrame

	var headers Headers
	headers.SetSbeMessageHeader(&zbsbe.MessageHeader{
		BlockLength: req.SbeBlockLength(),
		TemplateId:  req.SbeTemplateId(),
		SchemaId:    req.SbeSchemaId(),
		Version:     req.SbeSchemaVersion(),
	})
	headers.SetRequestResponseHeader(zbprotocol.NewRequestResponseHeader())
	headers.SetTransportHeader(zbprotocol.NewTransportHeader(zbprotocol.RequestResponse))

	// Writer will set FrameHeader after serialization to byte array.
	headers.SetFrameHeader(zbprotocol.NewFrameHeader(uint32(length), 0, 0, 0, 0))
	msg.SetHeaders(&headers)

	return &msg
}

func (rf *RequestHandler) createTaskRequest(commandRequest *zbsbe.ExecuteCommandRequest, task *zbmsgpack.Task) *Message {
	commandRequest.EventType = zbsbe.EventTypeEnum(0)

	if task.Payload == nil {
		b, err := msgpack.Marshal(task.PayloadJSON)
		if err != nil {
			return nil
		}
		task.Payload = b
	}

	return rf.newCommandMessage(commandRequest, task)

}

func (rf *RequestHandler) completeTaskRequest(taskMessage *Message) *Message {
	m, _ := taskMessage.ParseToMap()
	payload := *m
	payload["state"] = "COMPLETE"

	cmdReq := &zbsbe.ExecuteCommandRequest{
		PartitionId: (*taskMessage.SbeMessage).(*zbsbe.SubscribedEvent).PartitionId,
		Position:    (*taskMessage.SbeMessage).(*zbsbe.SubscribedEvent).Position,
		Key:         (*taskMessage.SbeMessage).(*zbsbe.SubscribedEvent).Key,
		TopicName:   (*taskMessage.SbeMessage).(*zbsbe.SubscribedEvent).TopicName,
	}
	return rf.newCommandMessage(cmdReq, payload)
}

func (rf *RequestHandler) createWorkflowInstanceRequest(commandRequest *zbsbe.ExecuteCommandRequest, wf *zbmsgpack.WorkflowInstance) *Message {
	commandRequest.EventType = zbsbe.EventTypeEnum(5)

	if wf.Payload == nil {
		b, err := msgpack.Marshal(wf.PayloadJSON)
		if err != nil {
			return nil
		}
		wf.Payload = b
	}

	return rf.newCommandMessage(commandRequest, wf)
}

func (rf *RequestHandler) topologyRequest() *Message {
	t := &zbmsgpack.TopologyRequest{}
	cmr := &zbsbe.ControlMessageRequest{
		MessageType: zbsbe.ControlMessageType.REQUEST_TOPOLOGY,
		Data:        nil,
	}

	return rf.newControlMessage(cmr, t)
}

func (rf *RequestHandler) newDeploymentRequest(commandRequest *zbsbe.ExecuteCommandRequest, d *zbmsgpack.Deployment) *Message {
	commandRequest.EventType = zbsbe.EventTypeEnum(4)
	return rf.newCommandMessage(commandRequest, d)
}

func (rf *RequestHandler) openTaskSubscriptionRequest(ts *zbmsgpack.TaskSubscription) *Message {
	var msg Message

	b, err := msgpack.Marshal(ts)
	if err != nil {
		return nil
	}
	controlRequest := &zbsbe.ControlMessageRequest{
		MessageType: zbsbe.ControlMessageType.ADD_TASK_SUBSCRIPTION,
		Data:        b,
	}
	msg.SetSbeMessage(controlRequest)

	length := 1 + uint32(LengthFieldSize+len(controlRequest.Data)) + TotalHeaderSizeNoFrame

	var headers Headers
	headers.SetSbeMessageHeader(&zbsbe.MessageHeader{
		BlockLength: controlRequest.SbeBlockLength(),
		TemplateId:  controlRequest.SbeTemplateId(),
		SchemaId:    controlRequest.SbeSchemaId(),
		Version:     controlRequest.SbeSchemaVersion(),
	})

	headers.SetRequestResponseHeader(zbprotocol.NewRequestResponseHeader())
	headers.SetTransportHeader(zbprotocol.NewTransportHeader(zbprotocol.RequestResponse))
	headers.SetFrameHeader(zbprotocol.NewFrameHeader(uint32(length), 0, 0, 0, 0))

	msg.SetHeaders(&headers)
	return &msg
}

func (rf *RequestHandler) closeTaskSubscriptionRequest(ts *zbmsgpack.TaskSubscription) *Message {
	var msg Message

	b, err := msgpack.Marshal(ts)
	if err != nil {
		return nil
	}
	controlRequest := &zbsbe.ControlMessageRequest{
		MessageType: zbsbe.ControlMessageType.REMOVE_TASK_SUBSCRIPTION,
		Data:        b,
	}
	msg.SetSbeMessage(controlRequest)

	length := 1 + uint32(LengthFieldSize+len(controlRequest.Data)) + TotalHeaderSizeNoFrame

	var headers Headers
	headers.SetSbeMessageHeader(&zbsbe.MessageHeader{
		BlockLength: controlRequest.SbeBlockLength(),
		TemplateId:  controlRequest.SbeTemplateId(),
		SchemaId:    controlRequest.SbeSchemaId(),
		Version:     controlRequest.SbeSchemaVersion(),
	})

	headers.SetRequestResponseHeader(zbprotocol.NewRequestResponseHeader())
	headers.SetTransportHeader(zbprotocol.NewTransportHeader(zbprotocol.RequestResponse))
	headers.SetFrameHeader(zbprotocol.NewFrameHeader(uint32(length), 0, 0, 0, 0))

	msg.SetHeaders(&headers)
	return &msg
}

func (rf *RequestHandler) openTopicSubscriptionRequest(cmdReq *zbsbe.ExecuteCommandRequest, ts *zbmsgpack.TopicSubscription) *Message {
	var msg Message

	b, err := msgpack.Marshal(ts)
	if err != nil {
		return nil
	}
	cmdReq.Command = b
	msg.SetSbeMessage(cmdReq)

	length := uint32(LengthFieldSize+len(cmdReq.TopicName)) + uint32(LengthFieldSize+len(cmdReq.Command))
	length += uint32(cmdReq.SbeBlockLength()) + TotalHeaderSizeNoFrame

	var headers Headers
	headers.SetSbeMessageHeader(&zbsbe.MessageHeader{
		BlockLength: cmdReq.SbeBlockLength(),
		TemplateId:  cmdReq.SbeTemplateId(),
		SchemaId:    cmdReq.SbeSchemaId(),
		Version:     cmdReq.SbeSchemaVersion(),
	})

	headers.SetRequestResponseHeader(zbprotocol.NewRequestResponseHeader())
	headers.SetTransportHeader(zbprotocol.NewTransportHeader(zbprotocol.RequestResponse))
	headers.SetFrameHeader(zbprotocol.NewFrameHeader(uint32(length), 0, 0, 0, 0))

	msg.SetHeaders(&headers)
	return &msg
}

func (rf *RequestHandler) topicSubscriptionAckRequest(cmdReq *zbsbe.ExecuteCommandRequest, tsa *zbmsgpack.TopicSubscriptionAck) *Message {
	var msg Message

	b, err := msgpack.Marshal(tsa)
	if err != nil {
		return nil
	}
	cmdReq.Command = b
	msg.SetSbeMessage(cmdReq)

	length := uint32(LengthFieldSize+len(cmdReq.TopicName)) + uint32(LengthFieldSize+len(cmdReq.Command))
	length += uint32(cmdReq.SbeBlockLength()) + TotalHeaderSizeNoFrame

	var headers Headers
	headers.SetSbeMessageHeader(&zbsbe.MessageHeader{
		BlockLength: cmdReq.SbeBlockLength(),
		TemplateId:  cmdReq.SbeTemplateId(),
		SchemaId:    cmdReq.SbeSchemaId(),
		Version:     cmdReq.SbeSchemaVersion(),
	})

	headers.SetRequestResponseHeader(zbprotocol.NewRequestResponseHeader())
	headers.SetTransportHeader(zbprotocol.NewTransportHeader(zbprotocol.RequestResponse))
	headers.SetFrameHeader(zbprotocol.NewFrameHeader(uint32(length), 0, 0, 0, 0))

	msg.SetHeaders(&headers)
	return &msg
}
