package protocol

import (
	"bytes"
	"encoding/binary"
	"testing"
	"unsafe"
)

func TestNewFrameHeader(t *testing.T) {
	obj := NewFrameHeader(0, 0, 0, 0, 0)
	if int(unsafe.Sizeof(*obj)) == 0 {
		t.Fatal("Object not allocated.")
	}
}

func TestFrameHeader_Encode(t *testing.T) {
	expected := []byte{
		0x93, 0x00, 0x00, 0x00,
		0x00,
		0x00,
		0x00, 0x00,
		0x02, 0x00, 0x00, 0x00,
	}

	frame := NewFrameHeader(147, 0, 0, 0, 2)

	if frame.Length != 147 {
		t.Fatal("FrameHeader construction failed. Wrong Length.")
	}
	if frame.StreamID != 2 {
		t.Fatal("FrameHeader construction failed. Wrong StreamId.")
	}

	byteBuff := &bytes.Buffer{}
	err := frame.Encode(byteBuff)

	if err != nil {
		t.Fatalf("Oh boy! %s", err)
	}

	if len(byteBuff.Bytes()) != 12 {
		t.Fatalf("Encoded buffered is not correct size. Received size: %+v", byteBuff.Bytes())
	}

	for i, b := range byteBuff.Bytes() {
		if b != expected[i] {
			t.Fatalf("Found unmatching bytes at index %d. ", i)
		}
	}
}

func TestFrameHeader_Decode(t *testing.T) {
	payload := []byte{
		0x93, 0x00, 0x00, 0x00,
		0x00,
		0x00,
		0x00, 0x00,
		0x02, 0x00, 0x00, 0x00,
	}

	if len(payload) != 12 {
		t.Fatal("Wrong header size.")
	}

	var frameHeader FrameHeader
	err := frameHeader.Decode(bytes.NewReader(payload), binary.LittleEndian, 0)

	if err != nil {
		t.Fatalf("Oh boy! %s", err)
	}

	if frameHeader.Length != 147 {
		t.Fatalf("Wrong DataFrame length. Received %d", frameHeader.Length)
	}
	if frameHeader.Version != 0 {
		t.Fatal("Wrong Version.")
	}
	if frameHeader.Flags != 0 {
		t.Fatal("Wrong Flag.")
	}
	if frameHeader.TypeID != FrameTypeMessage {
		t.Fatal("Wrong TypeId.")
	}
	if frameHeader.StreamID != 2 {
		t.Fatalf("Wrong StreamId. Received %d", frameHeader.StreamID)
	}
}
