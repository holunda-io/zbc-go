package protocol

import (
	"encoding/binary"
	"io"
)

const (
	RequestResponse = iota
	FullDuplexSingleMessage
)

type TransportHeader struct {
	ProtocolId uint16
}

func (fh TransportHeader) Encode(writer io.Writer) error {
	return binary.Write(writer, binary.LittleEndian, fh)
}

func (fh *TransportHeader) Decode(reader io.Reader, order binary.ByteOrder, _ uint16) error {
	return binary.Read(reader, order, fh)
}

func NewTransportHeader(pid uint16) *TransportHeader {
	return &TransportHeader{
		pid,
	}
}
