package protocol

import (
	"encoding/binary"
	"io"
)

const (
	FrameType_Message = iota
	Padding
	Unknown
)

const (
	_            = iota
	ControlClose = 100 + iota
	ControlEndOfStream
	ControlKeepAlive
	ProtocolControlFrame
)

type FrameHeader struct {
	Length   uint32 // 4 bytes
	Version  uint8  // 1 byte
	Flags    uint8  // 1 byte
	TypeId   uint16 // 2 byte One of the above defined constants.
	StreamId uint32 // 4 bytes

}

func (fh FrameHeader) Encode(writer io.Writer) error {
	return binary.Write(writer, binary.LittleEndian, fh)
}

func (fh *FrameHeader) Decode(reader io.Reader, order binary.ByteOrder, _ uint16) error {
	return binary.Read(reader, order, fh)
}

func NewFrameHeader(length uint32, version uint8, flags uint8, typeId uint16, streamId uint32) *FrameHeader {
	return &FrameHeader{
		length,
		version,
		flags,
		typeId,
		streamId,
	}
}
