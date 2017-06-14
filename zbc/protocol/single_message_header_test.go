package protocol

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"testing"
	"unsafe"
)

func TestNewSingleMessageHeader(t *testing.T) {
	obj := NewSingleMessageHeader()
	if int(unsafe.Sizeof(*obj)) != 0 {
		t.Fatal("Expected size is wrong.")
	}
}

func TestSingleMessageHeader_Encode(t *testing.T) {
	var single SingleMessageHeader

	buff := bytes.Buffer{}
	err := single.Encode(&buff)

	if err != nil || buff.Len() != 0 {
		t.Fatalf("SingleMessageHeader encoding failed: %s", err)
	}

	expected := []byte(nil)
	if !reflect.DeepEqual(buff.Bytes(), expected) {
		t.Fatalf("Encoding failed. Expected %+#v, received: %+#v", expected, buff.Bytes())
	}
}

func TestSingleMessageHeader_Decode(t *testing.T) {
	payload := []byte(nil)

	var single SingleMessageHeader
	err := single.Decode(bytes.NewReader(payload), binary.LittleEndian, 0)

	if err != nil {
		t.Fatalf("Decoding went wrong. %s", err)
	}

	expected := SingleMessageHeader{}
	if !reflect.DeepEqual(single, expected) {
		t.Fatalf("Decoding failed. Expected %+#v, received %+#v", expected, single)
	}
}
