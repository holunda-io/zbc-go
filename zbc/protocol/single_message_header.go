package protocol

import (
	"encoding/binary"
	"io"
)

type SingleMessageHeader struct{} // 0 byte

func (fh SingleMessageHeader) Encode(writer io.Writer) error {
	return binary.Write(writer, binary.LittleEndian, fh)
}

func (fh *SingleMessageHeader) Decode(reader io.Reader, order binary.ByteOrder, _ uint16) error {
	return binary.Read(reader, order, fh)
}

func NewSingleMessageHeader() *SingleMessageHeader {
	return &SingleMessageHeader{}
}
