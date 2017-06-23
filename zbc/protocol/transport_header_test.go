package protocol

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"testing"
	"unsafe"
)

func TestNewTransportHeader(t *testing.T) {
	obj := NewTransportHeader(RequestResponse)
	if int(unsafe.Sizeof(*obj)) == 0 {
		t.Fatal("Object not allocated.")
	}
}

func TestTransportHeader_Encode(t *testing.T) {
	var transport TransportHeader
	transport.ProtocolID = RequestResponse

	buff := bytes.Buffer{}
	err := transport.Encode(&buff)

	if err != nil || buff.Len() != 2 {
		t.Fatalf("Transporting encoding failed: %s", err)
	}

	if !reflect.DeepEqual(buff.Bytes(), []byte{0x00, 0x00}) {
		t.Fatalf("Encoding failed. Expected %+v, received: %+v", []byte{0x00, 0x00}, buff.Bytes())
	}
}

func TestTransportHeader_Encode2(t *testing.T) {
	var transport TransportHeader
	transport.ProtocolID = FullDuplexSingleMessage

	buff := bytes.Buffer{}
	err := transport.Encode(&buff)

	if err != nil || buff.Len() == 0 {
		t.Fatalf("Transporting encoding failed: %s", err)
	}

	if !reflect.DeepEqual(buff.Bytes(), []byte{0x01, 0x0}) {
		t.Fatalf("Encoding failed. Expected %+v, received: %+v", []byte{0x01, 0x00}, buff.Bytes())
	}
}

func TestTransportHeader_Decode(t *testing.T) {
	payload := []byte{0x00, 0x00}

	var transport TransportHeader
	err := transport.Decode(bytes.NewReader(payload), binary.LittleEndian, 0)

	if err != nil {
		t.Fatalf("Decoding went wrong. %s", err)
	}

	if transport.ProtocolID != RequestResponse {
		t.Fatal("Wrong ProtocolId.")
	}
}

func TestTransportHeader_Decode2(t *testing.T) {
	payload := []byte{0x01, 0x00}

	var transport TransportHeader
	err := transport.Decode(bytes.NewReader(payload), binary.LittleEndian, 0)

	if err != nil {
		t.Fatalf("Decoding went wrong. %s", err)
	}

	if transport.ProtocolID != FullDuplexSingleMessage {
		t.Fatal("Wrong ProtocolId.")
	}
}
