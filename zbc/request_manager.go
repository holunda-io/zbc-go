package zbc

import (
	"github.com/zeebe-io/zbc-go/zbc/zbmsgpack"
	"github.com/zeebe-io/zbc-go/zbc/zbsbe"
)

type requestManager struct {
	*requestFactory
	*responseHandler

	*topologyManager
}

func (rm *requestManager) partitionRequest() (*zbmsgpack.PartitionCollection, error) {
	message := rm.createPartitionRequest()
	request := newRequestWrapper(message)
	resp, err := rm.executeRequest(request)
	if err != nil {
		return nil, err
	}
	return rm.unmarshalPartition(resp), nil
}

func (rm *requestManager) createTask(topic string, task *zbmsgpack.Task) (*zbmsgpack.Task, error) {
	partitionID, err := rm.partitionID(topic)

	if err != nil {
		return nil, err
	}

	message := rm.createTaskRequest(partitionID, 0, task)
	request := newRequestWrapper(message)
	resp, err := rm.executeRequest(request)
	if err != nil {
		return nil, err
	}
	return rm.unmarshalTask(resp), nil
}

func (rm *requestManager) createWorkflow(topic string, resource []*zbmsgpack.Resource) (*zbmsgpack.Workflow, error) {
	message := rm.deployWorkflowRequest(topic, resource)
	request := newRequestWrapper(message)
	resp, err := rm.executeRequest(request)
	if err != nil {
		return nil, err
	}
	return rm.unmarshalWorkflow(resp), nil
}

func (rm *requestManager) createWorkflowInstance(topic string, wfi *zbmsgpack.WorkflowInstance) (*zbmsgpack.WorkflowInstance, error) {
	partitionID, err := rm.partitionID(topic)

	if err != nil {
		return nil, err
	}

	message := rm.createWorkflowInstanceRequest(partitionID, 0, topic, wfi)
	request := newRequestWrapper(message)
	resp, err := rm.executeRequest(request)
	if err != nil {
		return nil, err
	}
	return rm.unmarshalWorkflowInstance(resp), nil
}

func (rm *requestManager) completeTask(task *SubscriptionEvent) (*zbmsgpack.Task, error) {
	message := rm.completeTaskRequest(task)
	request := newRequestWrapper(message)
	resp, err := rm.executeRequest(request)
	if err != nil {
		return nil, err
	}
	return rm.unmarshalTask(resp), nil
}

func (rm *requestManager) taskConsumer(topic, lockOwner, taskType string, credits int32) (chan *SubscriptionEvent, *zbmsgpack.TaskSubscriptionInfo, error) {
	partitions, err := rm.topicPartitionsAddrs(topic)
	if err != nil {
		return nil, nil, err
	}

	send := func(endSubscriptionCh chan *SubscriptionEvent, subscriptionCh <-chan *SubscriptionEvent) {
		for {
			msg := <-subscriptionCh
			endSubscriptionCh <- msg
		}
	}

	tsi := zbmsgpack.NewTaskSubscriptionInfo()
	endSubscriptionCh := make(chan *SubscriptionEvent)

	for partitionID := range *partitions {
		subscriptionCh := make(chan *SubscriptionEvent, credits)
		message := rm.openTaskSubscriptionRequest(partitionID, lockOwner, taskType, credits)
		request := newRequestWrapper(message)
		resp, err := rm.executeRequest(request)
		if err != nil {
			return nil, nil, err
		}

		taskSubInfo := rm.unmarshalTaskSubscription(resp)
		if taskSubInfo != nil {
			tsi.AddSubInfo(*taskSubInfo)
			request.sock.addTaskSubscription(taskSubInfo.SubscriberKey, subscriptionCh)
			go send(endSubscriptionCh, subscriptionCh)
		}

	}

	return endSubscriptionCh, tsi, nil
}

func (rm *requestManager) increaseTaskSubscriptionCredits(task *zbmsgpack.TaskSubscription) (*zbmsgpack.TaskSubscription, error) {
	message := rm.increaseTaskSubscriptionCreditsRequest(task)

	request := newRequestWrapper(message)

	resp, err := rm.executeRequest(request)
	if err != nil {
		return nil, err
	}

	return rm.unmarshalTaskSubscription(resp), nil
}

func (rm *requestManager) closeTaskSubscription(sub *zbmsgpack.TaskSubscriptionInfo) []error {
	var errs []error
	for _, taskSub := range sub.Subs {
		_, err := rm.closeTaskSubscriptionPartition(&taskSub)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (rm *requestManager) closeTaskSubscriptionPartition(task *zbmsgpack.TaskSubscription) (*Message, error) {
	message := rm.closeTaskSubscriptionRequest(task)
	request := newRequestWrapper(message)
	resp, err := rm.executeRequest(request)
	request.sock.removeTaskSubscription(task.SubscriberKey)
	return resp, err
}
func (rm *requestManager) closeTopicSubscription(sub *zbmsgpack.TopicSubscriptionInfo) []error {
	var errs []error
	for _, taskSub := range sub.Subs {
		_, err := rm.closeTopicSubscriptionPartition(&taskSub)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (rm *requestManager) closeTopicSubscriptionPartition(topicPartition *zbmsgpack.TopicSubscription) (*Message, error) {
	message := rm.closeTopicSubscriptionRequest(topicPartition)
	request := newRequestWrapper(message)
	resp, err := rm.executeRequest(request)
	return resp, err
}

func (rm *requestManager) topicSubscriptionAck(ts *zbmsgpack.TopicSubscription, s *SubscriptionEvent) (*zbmsgpack.TopicSubscriptionAck, error) {
	message := rm.topicSubscriptionAckRequest(ts, s)
	request := newRequestWrapper(message)
	resp, err := rm.executeRequest(request)
	return rm.unmarshalTopicSubAck(resp), err
}

func (rm *requestManager) createTopic(name string, partitionNum int) (*zbmsgpack.Topic, error) {
	topic := zbmsgpack.NewTopic(name, TopicCreate, partitionNum)
	message := rm.createTopicRequest(topic)
	request := newRequestWrapper(message)

	resp, err := rm.executeRequest(request)
	if err != nil {
		return nil, err
	}
	return rm.unmarshalTopic(resp), nil
}

func (rm *requestManager) topicConsumer(topic, subName string, startPosition int64) (chan *SubscriptionEvent, *zbmsgpack.TopicSubscriptionInfo, error) {
	partitions, err := rm.topicPartitionsAddrs(topic)
	if err != nil {
		return nil, nil, err
	}

	send := func(subInfo *zbmsgpack.TopicSubscription, endSubscriptionCh chan *SubscriptionEvent, subscriptionCh <-chan *SubscriptionEvent) {
		for {
			msg := <-subscriptionCh
			endSubscriptionCh <- msg
			rm.topicSubscriptionAck(subInfo, msg)
		}
	}

	tsi := zbmsgpack.NewTopicSubscriptionInfo()
	endSubscriptionCh := make(chan *SubscriptionEvent, 1000)

	for partitionID := range *partitions {
		subscriptionCh := make(chan *SubscriptionEvent, 1000)
		message := rm.openTopicSubscriptionRequest(partitionID, topic, subName, startPosition)
		request := newRequestWrapper(message)
		resp, err := rm.executeRequest(request)
		if err != nil {
			return nil, nil, err
		}

		cmdResponse := (*resp.SbeMessage).(*zbsbe.ExecuteCommandResponse)
		subscriberKey := cmdResponse.Key

		request.sock.addTopicSubscription(cmdResponse.Key, subscriptionCh)
		subscriptionInfo := zbmsgpack.TopicSubscription{
			TopicName:        topic,
			PartitionID:      partitionID,
			SubscriberKey:    subscriberKey,
			SubscriptionName: subName,
		}

		tsi.AddSubInfo(subscriptionInfo)
		go send(&subscriptionInfo, endSubscriptionCh, subscriptionCh)
	}

	return endSubscriptionCh, tsi, nil
}

func newRequestManager(bootstrapAddr string) *requestManager {
	return &requestManager{
		newRequestFactory(),
		newResponseHandler(),
		newTopologyManager(bootstrapAddr),
	}
}
