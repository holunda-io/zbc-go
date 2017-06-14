package zbc

import (
	"bytes"
	"encoding/binary"
)

type MessageWriter struct {
	message *Message
}

func (mw *MessageWriter) writeFrameHeader(writer *bytes.Buffer) error {
	err := mw.message.Headers.FrameHeader.Encode(writer)
	if err != nil {
		return err
	}
	return nil
}

func (mw *MessageWriter) writeTransportHeader(writer *bytes.Buffer) error {
	err := mw.message.Headers.TransportHeader.Encode(writer)
	if err != nil {
		return err
	}
	return nil
}

func (mw *MessageWriter) writeRequestResponseHeader(writer *bytes.Buffer) error {
	if mw.message.Headers.IsSingleMessage() {
		return nil
	}
	err := mw.message.Headers.RequestResponseHeader.Encode(writer)
	if err != nil {
		return err
	}
	return nil
}

func (mw *MessageWriter) writeSbeMessageHeader(writer *bytes.Buffer) error {
	if err := mw.message.Headers.SbeMessageHeader.Encode(writer, binary.LittleEndian); err != nil {
		return err
	}
	return nil
}

func (mw *MessageWriter) writeHeaders(writer *bytes.Buffer) error {
	if err := mw.writeFrameHeader(writer); err != nil {
		return err
	}
	if err := mw.writeTransportHeader(writer); err != nil {
		return err
	}
	if err := mw.writeRequestResponseHeader(writer); err != nil {
		return err
	}
	if err := mw.writeSbeMessageHeader(writer); err != nil {
		return err
	}
	return nil
}

func (mw *MessageWriter) writeMessage(writer *bytes.Buffer) error {
	switch mw.message.Headers.SbeMessageHeader.TemplateId {

	case SBE_ExecuteCommandRequest_TemplateId:
		if err := (*mw.message.SbeMessage).Encode(writer, binary.LittleEndian, false); err != nil {
			return err
		}
		return nil
	}

	return nil
}

func (mw *MessageWriter) align(writer *bytes.Buffer) {
	currentSize := len(writer.Bytes())
	expectedSize := (currentSize + 7) & ^7

	for currentSize < expectedSize {
		writer.Write([]byte{0x00})
		currentSize += 1
	}
}

func (mw *MessageWriter) Write(writer *bytes.Buffer) {
	mw.writeHeaders(writer)
	mw.writeMessage(writer)
	mw.align(writer)
}

func NewMessageWriter(message *Message) *MessageWriter {
	return &MessageWriter{
		message,
	}
}
