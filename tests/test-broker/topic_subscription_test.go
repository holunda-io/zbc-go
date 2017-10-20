package testbroker

import (
	"github.com/zeebe-io/zbc-go/zbc"
	"testing"
)

func TestTopicSubscription(t *testing.T) {
	zbClient, err := zbc.NewClient(brokerAddr)
	assert(t, nil, err, true)
	assert(t, nil, zbClient, false)

	workflow, err := zbClient.CreateWorkflowFromFile(topicName, zbc.BpmnXml, "../../examples/demoProcess.bpmn")
	assert(t, nil, err, true)
	assert(t, nil, workflow, false)
	assert(t, zbc.DeployementCreated, workflow.State, true)

	payload := make(map[string]interface{})
	payload["a"] = "b"

	instance := zbc.NewWorkflowInstance("demoProcess", -1, payload)

	for i := 0; i < 10; i++ {
		createdInstance, err := zbClient.CreateWorkflowInstance(topicName, instance)
		assert(t, nil, err, true)
		assert(t, nil, createdInstance, false)
		assert(t, zbc.WorkflowInstanceCreated, createdInstance.State, true)
	}

	subscriptionCh, subscription, err := zbClient.TopicConsumer(topicName, "default-name", 0)
	assert(t, nil, err, true)
	assert(t, nil, subscription, false)
	assert(t, nil, subscriptionCh, false)

	var message *zbc.SubscriptionEvent
	for i := 0; i < 10; i++ {
		message = <-subscriptionCh
		assert(t, nil, message, false)
	}

	msg, err := zbClient.TopicSubscriptionAck(subscription, message)
	assert(t, nil, err, true)
	assert(t, nil, msg, false)
	assert(t, "default-name", msg.Name, true)
	assert(t, message.Event.Position, msg.AckPosition, true)

	closeMsg, err := zbClient.CloseTopicSubscription(subscription)
	assert(t, nil, err, true)
	assert(t, nil, closeMsg, false)
}
