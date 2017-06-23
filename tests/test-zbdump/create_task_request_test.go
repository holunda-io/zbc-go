package test_zbdump

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/jsam/zbc-go/zbc"
	"github.com/jsam/zbc-go/zbc/sbe"
)

const (
	CloseChannelFile          = "dumps/close-channel"               // TODO:
	CloseSubscriptionRequest  = "dumps/close-subscription-request"  // TODO:
	CloseSubscriptionResponse = "dumps/close-subscription-response" // TODO:
	CreateTaskRequest         = "dumps/create-task-request"
	CreateTaskResponse        = "dumps/create-task-response" // TODO:
	EndOfStream               = "dumps/end-of-stream"        // TODO:
	Gossip                    = "dumps/gossip"               // TODO:
)

func TestCreateTaskRequest(t *testing.T) {
	reader, fsErr := os.Open(CreateTaskRequest)
	defer reader.Close()

	if fsErr != nil {
		t.Fatalf("Cannot read file %s", CreateTaskRequest)
	}

	buffer := bufio.NewReader(reader)
	msgReader := zbc.NewMessageReader(buffer)

	headers, message, err := msgReader.ReadHeaders()
	if err != nil {
		t.Fatalf("Something went wrong during reading of headers. %+#v", err)
	}

	if int(headers.FrameHeader.Length) != len(*message)+26 {
		t.Fatalf("FrameHeader contains wrong length. Expected %d, received %d", len(*message)+26, headers.FrameHeader.Length)
	}

	msg, err := msgReader.ParseMessage(headers, message)
	if err != nil {
		t.Fatalf("Nothing got decoded. Something is wrong! %+#v", err)
	}

	msgpackItem := *msg.Data
	if !reflect.DeepEqual([]byte(msgpackItem["eventType"].(string)), []byte("CREATE")) {
		t.Fatalf("Wrong eventType. Expected CREATE, received %s", msgpackItem["eventType"].(string))
	}

	if headers.FrameHeader.Length != 129 {
		t.Fatalf("Wrong FrameHeader Length. Expected %+#v, received: %d", 147, headers.FrameHeader.Length)
	}

	if err != nil {
		t.Fatalf("init failed: %+#v", err)
	}

}

func TestCreateTaskRequest_CommandRequest(t *testing.T) {
	reader, fsErr := os.Open(CreateTaskRequest)
	defer reader.Close()

	if fsErr != nil {
		t.Fatalf("Cannot read file %s", CreateTaskRequest)
	}

	buffer := bufio.NewReader(reader)
	msgReader := zbc.NewMessageReader(buffer)

	headers, message, err := msgReader.ReadHeaders()
	if int(headers.FrameHeader.Length) != len(*message)+26 {
		t.Fatalf("FrameHeader contains wrong length. Expected %d, received %d",
			len(*message)+26, headers.FrameHeader.Length)
	}

	if err != nil {
		t.Fatal("Cannot read headers.")
	}

	msg, err := msgReader.ParseMessage(headers, message)
	if err != nil {
		t.Fatal(err)
	}
	executeCmdRequestI := *msg.SbeMessage
	cmdReq := executeCmdRequestI.(*sbe.ExecuteCommandRequest)
	size := int(cmdReq.SbeBlockLength()) + 2 + len(cmdReq.TopicName) + 2 + len(cmdReq.Command) + 26

	if uint32(size) != headers.FrameHeader.Length {
		t.Fatalf("size of sbe.ExecuteCommandRequest is %d and FrameHeader.Length is %d", size, headers.FrameHeader.Length)
	}
}

func TestCreateTaskRequest_CommandRequest2(t *testing.T) {
	reader, fsErr := ioutil.ReadFile(CreateTaskRequest)
	if fsErr != nil {
		t.Fatalf("fserr %+#v", fsErr)
	}

	buffer := bytes.NewReader(reader)
	msgReader := zbc.MessageReader{Reader: buffer}
	headers, message, err := msgReader.ReadHeaders()

	if err != nil {
		t.Fatal("Cannot read headers.")
	}

	msg, err := msgReader.ParseMessage(headers, message)
	if err != nil {
		t.Fatal(err)
	}

	executeCmdRequestI := *msg.SbeMessage
	cmdReq := executeCmdRequestI.(*sbe.ExecuteCommandRequest)
	size := int(cmdReq.SbeBlockLength()) + 2 + len(cmdReq.TopicName) + 2 + len(cmdReq.Command) + 26

	if uint32(size) != headers.FrameHeader.Length {
		t.Fatalf("size of sbe.ExecuteCommandRequest is %d and FrameHeader.Length is %d", size, headers.FrameHeader.Length)
	}

	writer := zbc.NewMessageWriter(msg)
	byteBuffer := &bytes.Buffer{}
	writer.Write(byteBuffer)

	if !reflect.DeepEqual(reader, byteBuffer.Bytes()) {
		t.Fatalf("Expected %+#v, received %+#v", reader, byteBuffer.Bytes())
	}
}
