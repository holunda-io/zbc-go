package testzbdump

import (
	"os"
	"testing"

	"github.com/zeebe-io/zbc-go/zbc"
)

const (
	AppendRequest = "dumps/append-request.bin"

	CreateDeploymentRequest        = "dumps/create-deployment-request.bin"
	CreateDeploymentResponse       = "dumps/create-deployment-response.bin"
	CreateTaskRequest              = "dumps/create-task-request.bin"
	CreateTaskResponse             = "dumps/create-task-response.bin"
	CreateWorkflowInstanceRequest  = "dumps/create-workflow-instance-request.bin"
	CreateWorkflowInstanceResponse = "dumps/create-workflow-instance-response.bin"

	CloseTaskSubscriptionRequest   = "dumps/close-task-subscription-request.bin"
	CloseTaskSubscriptionResponse  = "dumps/close-task-subscription-response.bin"
	CloseTopicSubscriptionRequest  = "dumps/close-topic-subscription-request.bin"
	CloseTopicSubscriptionResponse = "dumps/close-topic-subscription-response.bin"

	ErrorTopicNotFound = "dumps/error-topic-not-found.bin"
	KeepAlive          = "dumps/keep-alive.bin"

	OpenTaskSubscriptionRequest   = "dumps/open-task-subscription-request.bin"
	OpenTaskSubscriptionResponse  = "dumps/open-task-subscription-response.bin"
	OpenTopicSubscriptionRequest  = "dumps/open-topic-subscription-request.bin"
	OpenTopicSubscriptionResponse = "dumps/open-topic-subscription-response.bin"

	TaskSubscriptionLockedTask = "dumps/task-subscription-locked-task.bin"
	TopologyRequest            = "dumps/topology-request.bin"
	TopologyResponse           = "dumps/topology-response.bin"
)

func TestAppendRequest(t *testing.T) {
	reader, fsErr := os.Open(AppendRequest)
	defer reader.Close()

	if fsErr != nil {
		t.Fatalf("Cannot read file %s", AppendRequest)
	}

	client := zbc.Client{}
	message, err := client.UnmarshalFromFile(AppendRequest)
	if err != nil {
		t.Fatalf("Error reading file.")
	}
	if message == nil {
		t.Fatalf("Message is nil.")
	}
}

//func TestCreateDeploymentRequest(t *testing.T) {
//	reader, fsErr := os.Open(CreateDeploymentRequest)
//	defer reader.Close()
//
//	if fsErr != nil {
//		t.Fatalf("Cannot read file %s", CreateDeploymentRequest)
//	}
//
//	client := zbc.Client{}
//	message, err := client.UnmarshalFromFile(CreateDeploymentRequest)
//	if err != nil {
//		t.Fatalf("Error reading file.", err)
//	}
//	if message == nil {
//		t.Fatalf("Message is nil.")
//	}
//}

func TestCreateDeploymentResponse(t *testing.T) {
	reader, fsErr := os.Open(CreateDeploymentResponse)
	defer reader.Close()

	if fsErr != nil {
		t.Fatalf("Cannot read file %s", CreateDeploymentResponse)
	}

	client := zbc.Client{}
	message, err := client.UnmarshalFromFile(CreateDeploymentResponse)
	if err != nil {
		t.Fatalf("Error reading file.")
	}
	if message == nil {
		t.Fatalf("Message is nil.")
	}
}

func TestCreateTaskRequest(t *testing.T) {
	reader, fsErr := os.Open(CreateTaskRequest)
	defer reader.Close()

	if fsErr != nil {
		t.Fatalf("Cannot read file %s", CreateTaskRequest)
	}

	client := zbc.Client{}
	message, err := client.UnmarshalFromFile(CreateTaskRequest)
	if err != nil {
		t.Fatalf("Error reading file.")
	}
	if message == nil {
		t.Fatalf("Message is nil.")
	}
}

func TestCreateTaskResponse(t *testing.T) {
	reader, fsErr := os.Open(CreateTaskResponse)
	defer reader.Close()

	if fsErr != nil {
		t.Fatalf("Cannot read file %s", CreateTaskResponse)
	}

	client := zbc.Client{}
	message, err := client.UnmarshalFromFile(CreateTaskResponse)
	if err != nil {
		t.Fatalf("Error reading file.")
	}
	if message == nil {
		t.Fatalf("Message is nil.")
	}
}

//func TestWorkflowInstanceRequest(t *testing.T) {
//	reader, fsErr := os.Open(CreateWorkflowInstanceRequest)
//	defer reader.Close()
//
//	if fsErr != nil {
//		t.Fatalf("Cannot read file %s", CreateWorkflowInstanceRequest)
//	}
//
//	client := zbc.Client{}
//	message, err := client.UnmarshalFromFile(CreateWorkflowInstanceRequest)
//	if err != nil {
//		t.Fatalf("Error reading file.")
//	}
//	if message == nil {
//		t.Fatalf("Message is nil.")
//	}
//}

func TestWorkflowInstanceResponse(t *testing.T) {
	reader, fsErr := os.Open(CreateWorkflowInstanceResponse)
	defer reader.Close()

	if fsErr != nil {
		t.Fatalf("Cannot read file %s", CreateWorkflowInstanceResponse)
	}

	client := zbc.Client{}
	message, err := client.UnmarshalFromFile(CreateWorkflowInstanceResponse)
	if err != nil {
		t.Fatalf("Error reading file.")
	}
	if message == nil {
		t.Fatalf("Message is nil.")
	}
}

func TestCloseTaskSubscriptionRequest(t *testing.T) {
	reader, fsErr := os.Open(CloseTaskSubscriptionRequest)
	defer reader.Close()

	if fsErr != nil {
		t.Fatalf("Cannot read file %s", CloseTaskSubscriptionRequest)
	}

	client := zbc.Client{}
	message, err := client.UnmarshalFromFile(CloseTaskSubscriptionRequest)
	if err != nil {
		t.Fatalf("Error reading file.")
	}
	if message == nil {
		t.Fatalf("Message is nil.")
	}
}

func TestCloseTaskSubscriptionResponse(t *testing.T) {
	reader, fsErr := os.Open(CloseTaskSubscriptionResponse)
	defer reader.Close()

	if fsErr != nil {
		t.Fatalf("Cannot read file %s", CloseTaskSubscriptionResponse)
	}

	client := zbc.Client{}
	message, err := client.UnmarshalFromFile(CloseTaskSubscriptionResponse)
	if err != nil {
		t.Fatalf("Error reading file.")
	}
	if message == nil {
		t.Fatalf("Message is nil.")
	}
}

func TestCloseTopicSubscriptionRequest(t *testing.T) {
	reader, fsErr := os.Open(CloseTopicSubscriptionRequest)
	defer reader.Close()

	if fsErr != nil {
		t.Fatalf("Cannot read file %s", CloseTopicSubscriptionRequest)
	}

	client := zbc.Client{}
	message, err := client.UnmarshalFromFile(CloseTopicSubscriptionRequest)
	if err != nil {
		t.Fatalf("Error reading file.")
	}
	if message == nil {
		t.Fatalf("Message is nil.")
	}
}

func TestCloseTopicSubscriptionResponse(t *testing.T) {
	reader, fsErr := os.Open(CloseTopicSubscriptionResponse)
	defer reader.Close()

	if fsErr != nil {
		t.Fatalf("Cannot read file %s", CloseTopicSubscriptionResponse)
	}

	client := zbc.Client{}
	message, err := client.UnmarshalFromFile(CloseTopicSubscriptionResponse)
	if err != nil {
		t.Fatalf("Error reading file.")
	}
	if message == nil {
		t.Fatalf("Message is nil.")
	}
}

func TestErrorTopicNotFound(t *testing.T) {
	reader, fsErr := os.Open(ErrorTopicNotFound)
	defer reader.Close()

	if fsErr != nil {
		t.Fatalf("Cannot read file %s", ErrorTopicNotFound)
	}

	client := zbc.Client{}
	message, err := client.UnmarshalFromFile(ErrorTopicNotFound)
	if err != nil {
		t.Fatalf("Error reading file.")
	}
	if message == nil {
		t.Fatalf("Message is nil.")
	}
}

//func TestKeepAlive(t *testing.T) {
//	reader, fsErr := os.Open(KeepAlive)
//	defer reader.Close()
//
//	if fsErr != nil {
//		t.Fatalf("Cannot read file %s", KeepAlive)
//	}
//
//	client := zbc.Client{}
//	message, err := client.UnmarshalFromFile(KeepAlive)
//	if err != nil {
//		t.Fatalf("Error reading file.")
//	}
//	if message == nil {
//		t.Fatalf("Message is nil.")
//	}
//}

//func TestOpenTaskSubscriptionRequest(t *testing.T) {
//	reader, fsErr := os.Open(OpenTaskSubscriptionRequest)
//	defer reader.Close()
//
//	if fsErr != nil {
//		t.Fatalf("Cannot read file %s", OpenTaskSubscriptionRequest)
//	}
//
//	client := zbc.Client{}
//	message, err := client.UnmarshalFromFile(OpenTaskSubscriptionRequest)
//	if err != nil {
//		t.Fatalf("Error reading file.")
//	}
//	if message == nil {
//		t.Fatalf("Message is nil.")
//	}
//}

//func TestOpenTaskSubscriptionResponse(t *testing.T) {
//	reader, fsErr := os.Open(OpenTaskSubscriptionResponse)
//	defer reader.Close()
//
//	if fsErr != nil {
//		t.Fatalf("Cannot read file %s", OpenTaskSubscriptionResponse)
//	}
//
//	client := zbc.Client{}
//	message, err := client.UnmarshalFromFile(OpenTaskSubscriptionResponse)
//	if err != nil {
//		t.Fatalf("Error reading file.")
//	}
//	if message == nil {
//		t.Fatalf("Message is nil.")
//	}
//}

//func TestOpenTopicSubscriptionRequest(t *testing.T) {
//	reader, fsErr := os.Open(OpenTopicSubscriptionRequest)
//	defer reader.Close()
//
//	if fsErr != nil {
//		t.Fatalf("Cannot read file %s", OpenTopicSubscriptionRequest)
//	}
//
//	client := zbc.Client{}
//	message, err := client.UnmarshalFromFile(OpenTopicSubscriptionRequest)
//	if err != nil {
//		t.Fatalf("Error reading file.")
//	}
//	if message == nil {
//		t.Fatalf("Message is nil.")
//	}
//}
//
//func TestOpenTopicSubscriptionResponse(t *testing.T) {
//	reader, fsErr := os.Open(OpenTopicSubscriptionResponse)
//	defer reader.Close()
//
//	if fsErr != nil {
//		t.Fatalf("Cannot read file %s", OpenTopicSubscriptionResponse)
//	}
//
//	client := zbc.Client{}
//	message, err := client.UnmarshalFromFile(OpenTopicSubscriptionResponse)
//	if err != nil {
//		t.Fatalf("Error reading file.")
//	}
//	if message == nil {
//		t.Fatalf("Message is nil.")
//	}
//}

func TestTaskSubscriptionLockedTask(t *testing.T) {
	reader, fsErr := os.Open(TaskSubscriptionLockedTask)
	defer reader.Close()

	if fsErr != nil {
		t.Fatalf("Cannot read file %s", TaskSubscriptionLockedTask)
	}

	client := zbc.Client{}
	message, err := client.UnmarshalFromFile(TaskSubscriptionLockedTask)
	if err != nil {
		t.Fatalf("Error reading file.")
	}
	if message == nil {
		t.Fatalf("Message is nil.")
	}
}

func TestTopologyRequest(t *testing.T) {
	reader, fsErr := os.Open(TopologyRequest)
	defer reader.Close()

	if fsErr != nil {
		t.Fatalf("Cannot read file %s", TopologyRequest)
	}

	client := zbc.Client{}
	message, err := client.UnmarshalFromFile(TopologyRequest)
	if err != nil {
		t.Fatalf("Error reading file.")
	}
	if message == nil {
		t.Fatalf("Message is nil.")
	}
}

func TestTopologyResponse(t *testing.T) {
	reader, fsErr := os.Open(TopologyResponse)
	defer reader.Close()

	if fsErr != nil {
		t.Fatalf("Cannot read file %s", TopologyResponse)
	}

	client := zbc.Client{}
	message, err := client.UnmarshalFromFile(TopologyResponse)
	if err != nil {
		t.Fatalf("Error reading file.")
	}
	if message == nil {
		t.Fatalf("Message is nil.")
	}
}
