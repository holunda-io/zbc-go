package protocol

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"testing"
	"unsafe"
)

func TestNewRequestResponseHeader(t *testing.T) {
	obj := NewRequestResponseHeader()
	if int(unsafe.Sizeof(*obj)) == 0 {
		t.Fatal("Object not allocated.")
	}
}

func TestRequestResponseHeader_Encode(t *testing.T) {
	var rrHeader RequestResponseHeader
	rrHeader.ConnectionID = 0
	rrHeader.RequestID = 1

	buff := bytes.Buffer{}
	err := rrHeader.Encode(&buff)

	if err != nil || len(buff.Bytes()) != 16 {
		t.Fatal("RequestResponseHeader encoding failed.")
	}

	expected := []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	if !reflect.DeepEqual(buff.Bytes(), expected) {
		t.Fatalf("Encoding failed. Expected %+v, recevied %+v", expected, buff.Bytes())
	}
}

func TestRequestResponseHeader_Decode(t *testing.T) {
	payload := []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	var rrHeader RequestResponseHeader
	err := rrHeader.Decode(bytes.NewReader(payload), binary.LittleEndian, 0)

	if err != nil {
		t.Fatalf("Decoding went wrong. %s", err)
	}

	if rrHeader.RequestID != 1 {
		t.Fatalf("Wrong RequestId. Expected 1, received %d.", rrHeader.RequestID)
	}

	if rrHeader.ConnectionID != 0 {
		t.Fatalf("Wrong ConnectionId. Expected 0, received %d.", rrHeader.ConnectionID)
	}
}
