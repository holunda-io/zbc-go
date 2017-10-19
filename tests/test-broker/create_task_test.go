package testbroker

import (
	"testing"
	"github.com/zeebe-io/zbc-go/zbc"
)

func TestCreateTask(t *testing.T) {
	zbClient, err := zbc.NewClient(brokerAddr)
	assert(t, nil, err, true)
	assert(t, nil, zbClient, false)

	task := zbc.NewTask("testType", "test-owner")
	responseTask, err := zbClient.CreateTask(topicName, task)
	assert(t, nil, err, true)
	assert(t, nil, responseTask, false)

	assert(t, nil, responseTask, false)
	assert(t, nil, responseTask.State, false)
	assert(t, responseTask.State, zbc.TaskCreated, true)
}
