package testbroker

import (
	"testing"

	"github.com/zeebe-io/zbc-go/zbc"
)

func TestTaskSubscription(t *testing.T) {
	zbClient, err := zbc.NewClient(brokerAddr)
	assert(t, nil, err, true)
	assert(t, nil, zbClient, false)

	workflow, err := zbClient.CreateWorkflowFromFile(topicName, zbc.BpmnXml, "../../examples/demoProcess.bpmn")
	assert(t, nil, err, true)
	assert(t, nil, workflow, false)
	assert(t, nil, workflow.State, false)
	assert(t, zbc.DeploymentCreated, workflow.State, true)

	payload := make(map[string]interface{})
	payload["a"] = "b"

	instance := zbc.NewWorkflowInstance("demoProcess", -1, payload)
	createdInstance, err := zbClient.CreateWorkflowInstance(topicName, instance)
	assert(t, nil, err, true)
	assert(t, nil, createdInstance, false)
	assert(t, zbc.WorkflowInstanceCreated, createdInstance.State, true)

	subscriptionCh, subscription, err := zbClient.TaskConsumer(topicName, "task_subscription_test", "foo")
	assert(t, nil, err, true)
	assert(t, nil, subscription, false)
	assert(t, nil, subscriptionCh, false)

	message := <-subscriptionCh

	response, err := zbClient.CompleteTask(message)
	assert(t, nil, err, true)
	assert(t, nil, response, false)

	errs := zbClient.CloseTaskSubscription(subscription)
	assert(t, 0, len(errs), true)
}
