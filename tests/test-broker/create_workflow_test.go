package testbroker

import (
	"github.com/zeebe-io/zbc-go/zbc"
	"testing"
)

func TestCreateWorkflow(t *testing.T) {
	zbClient, err := zbc.NewClient(brokerAddr)
	assert(t, nil, err, true)
	assert(t, nil, zbClient, false)

	workflow, err := zbClient.CreateWorkflowFromFile(topicName, zbc.BpmnXml, "../../examples/demoProcess.bpmn")
	assert(t, nil, err, true)
	assert(t, nil, workflow, false)
	assert(t, zbc.DeployementCreated, workflow.State, true)
}
