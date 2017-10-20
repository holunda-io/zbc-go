package testbroker

import (
	"github.com/zeebe-io/zbc-go/zbc"
	"testing"
)

func TestCreateInstance(t *testing.T) {
	zbClient, err := zbc.NewClient(brokerAddr)
	assert(t, nil, err, true)
	assert(t, nil, zbClient, false)

	workflow, err := zbClient.CreateWorkflowFromFile(topicName, zbc.BpmnXml, "../../examples/demoProcess.bpmn")
	assert(t, nil, err, true)
	assert(t, nil, workflow, false)
	assert(t, zbc.DeployementCreated, workflow.State, true)

	instance := zbc.NewWorkflowInstance("demoProcess", -1, nil)
	createdInstance, err := zbClient.CreateWorkflowInstance(topicName, instance)
	assert(t, nil, err, true)
	assert(t, nil, createdInstance, false)
	assert(t, zbc.WorkflowInstanceCreated, createdInstance.State, true)
}
