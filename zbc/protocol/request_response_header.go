package protocol

import (
	"encoding/binary"
	"io"
	"math/rand"
	"time"
)

type RequestResponseHeader struct {
	ConnectionId uint64 //
	RequestId    uint64 // tid
}

func (fh RequestResponseHeader) Encode(writer io.Writer) error {
	return binary.Write(writer, binary.LittleEndian, fh)
}

func (fh *RequestResponseHeader) Decode(reader io.Reader, order binary.ByteOrder, _ uint16) error {
	return binary.Read(reader, order, fh)
}

func NewRequestResponseHeader() *RequestResponseHeader {
	max := ^uint64(0)
	var s float64 = 10
	var v float64 = 1000
	zipf := rand.NewZipf(rand.New(rand.NewSource(time.Now().UnixNano())), s, v, max)
	return &RequestResponseHeader{
		zipf.Uint64(),
		zipf.Uint64(),
	}
}
