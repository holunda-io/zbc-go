package protocol

import (
	"testing"
	"unsafe"
)

func TestNewSingleMessageHeader(t *testing.T) {
	obj := NewSingleMessageHeader()
	if int(unsafe.Sizeof(*obj)) != 0 {
		t.Fatal("Expected size is wrong.")
	}
}
