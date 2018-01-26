package zbc

import (
	"github.com/zeebe-io/zbc-go/zbc/zbmsgpack"
	"github.com/zeebe-io/zbc-go/zbc/zbsbe"
)

type dispatcher struct {
	lastTransactionSeed uint64
	activeTransactions  []*requestWrapper

	lastSubscriptionSeed uint64
	subscriptions        *subscriptionsManager
}

func (d *dispatcher) addTransaction(request *requestWrapper) {
	d.lastTransactionSeed += requestQueueSize + 1
	request.payload.Headers.RequestResponseHeader.RequestID = d.lastTransactionSeed

	index := d.lastTransactionSeed & (requestQueueSize - 1)
	d.activeTransactions[index] = request
}

func (d *dispatcher) dispatchTransaction(transactionID uint64, response *Message) {
	index := transactionID & (requestQueueSize - 1)
	request := d.activeTransactions[index]
	request.responseCh <- response
}

func (d *dispatcher) addTaskSubscription(key uint64, value interface{}) {
	d.subscriptions.addTaskSubscription(key, value)
}

func (d *dispatcher) addTopicSubscription(key uint64, value interface{}) {
	d.subscriptions.addTopicSubscription(key, value)
}

func (d *dispatcher) dispatchTaskEvent(key uint64, message *zbsbe.SubscribedEvent, task *zbmsgpack.Task) {
	if ch := d.subscriptions.getTaskChannel(key); ch != nil {
		ch <- &SubscriptionEvent{
			Task: task, Event: message,
		}
	}
}

func (d *dispatcher) dispatchTopicEvent(key uint64, message *zbsbe.SubscribedEvent) {
	if ch := d.subscriptions.getTopicChannel(key); ch != nil {
		ch <- &SubscriptionEvent{
			Task:  nil,
			Event: message,
		}
	}
}

func (d *dispatcher) removeTaskSubscription(key uint64) {
	d.subscriptions.removeTaskSubscription(key)
}

func (d *dispatcher) removeTopicSubscription(key uint64) {
	d.subscriptions.removeTopicSubscription(key)
}

func newDispatcher() dispatcher {
	return dispatcher{
		0,
		make([]*requestWrapper, requestQueueSize),
		0,
		newSubscriptionsManager(),
	}
}
