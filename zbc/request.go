package zbc

import (
	"github.com/zeebe-io/zbc-go/zbc/zbmsgpack"
	"github.com/zeebe-io/zbc-go/zbc/zbprotocol"
	"github.com/zeebe-io/zbc-go/zbc/zbsbe"
	"github.com/vmihailenco/msgpack"
)

type requestHandler struct{}

func (rf *requestHandler) newCommandMessage(commandRequest *zbsbe.ExecuteCommandRequest, command interface{}) *Message {
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

func (rf *requestHandler) newControlMessage(req *zbsbe.ControlMessageRequest, payload interface{}) *Message {
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

func (rf *requestHandler) createTaskRequest(commandRequest *zbsbe.ExecuteCommandRequest, task *zbmsgpack.Task) *Message {
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

func (rf *requestHandler) completeTaskRequest(taskMessage *TaskEvent) *Message {
	taskMessage.State = TaskComplete
	cmdReq := &zbsbe.ExecuteCommandRequest{
		PartitionId: taskMessage.Event.PartitionId,
		Position:    taskMessage.Event.Position,
		Key:         taskMessage.Event.Key,
		TopicName:   taskMessage.Event.TopicName,
	}
	return rf.newCommandMessage(cmdReq, taskMessage.Task)
}

func (rf *requestHandler) createWorkflowInstanceRequest(commandRequest *zbsbe.ExecuteCommandRequest, wf *zbmsgpack.WorkflowInstance) *Message {
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

func (rf *requestHandler) topologyRequest() *Message {
	t := &zbmsgpack.TopologyRequest{}
	cmr := &zbsbe.ControlMessageRequest{
		MessageType: zbsbe.ControlMessageType.REQUEST_TOPOLOGY,
		Data:        nil,
	}

	return rf.newControlMessage(cmr, t)
}

func (rf *requestHandler) newWorkflowRequest(commandRequest *zbsbe.ExecuteCommandRequest, d *zbmsgpack.Workflow) *Message {
	commandRequest.EventType = zbsbe.EventTypeEnum(4)
	return rf.newCommandMessage(commandRequest, d)
}

func (rf *requestHandler) openTaskSubscriptionRequest(ts *zbmsgpack.TaskSubscription) *Message {
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

func (rf *requestHandler) increaseTaskSubscriptionCreditsRequest(ts *zbmsgpack.TaskSubscription) *Message {
	var msg Message

	b, err := msgpack.Marshal(ts)
	if err != nil {
		return nil
	}
	controlRequest := &zbsbe.ControlMessageRequest{
		MessageType: zbsbe.ControlMessageType.INCREASE_TASK_SUBSCRIPTION_CREDITS,
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

func (rf *requestHandler) closeTaskSubscriptionRequest(ts *zbmsgpack.TaskSubscription) *Message {
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

func (rf *requestHandler) openTopicSubscriptionRequest(cmdReq *zbsbe.ExecuteCommandRequest, ts *zbmsgpack.TopicSubscription) *Message {
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

func (rf *requestHandler) topicSubscriptionAckRequest(cmdReq *zbsbe.ExecuteCommandRequest, tsa *zbmsgpack.TopicSubscriptionAck) *Message {
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

func (rf *requestHandler) createTopicRequest(cmdReq *zbsbe.ExecuteCommandRequest, t *zbmsgpack.Topic) *Message {
	var msg Message

	b, err := msgpack.Marshal(t)
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