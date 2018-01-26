package zbc

import "fmt"

type subscriptionsManager struct {
	taskSubscriptions  SafeMap
	topicSubscriptions SafeMap
}

func (sm *subscriptionsManager) addTaskSubscription(key uint64, value interface{}) {
	sm.taskSubscriptions.Set(fmt.Sprintf("%d", key), value)
}

func (sm *subscriptionsManager) addTopicSubscription(key uint64, value interface{}) {
	sm.topicSubscriptions.Set(fmt.Sprintf("%d", key), value)
}

func (sm *subscriptionsManager) removeTaskSubscription(key uint64) {
	sm.taskSubscriptions.Remove(fmt.Sprintf("%d", key))
}

func (sm *subscriptionsManager) removeTopicSubscription(key uint64) {
	sm.topicSubscriptions.Remove(fmt.Sprintf("%d", key))
}

func (sm *subscriptionsManager) getTaskChannel(key uint64) chan *SubscriptionEvent {
	if ch, ok := sm.taskSubscriptions.Get(fmt.Sprintf("%d", key)); ok {
		return ch.(chan *SubscriptionEvent)
	}
	return nil
}

func (sm *subscriptionsManager) getTopicChannel(key uint64) chan *SubscriptionEvent {
	if ch, ok := sm.topicSubscriptions.Get(fmt.Sprintf("%d", key)); ok {
		return ch.(chan *SubscriptionEvent)
	}
	return nil
}

func newSubscriptionsManager() *subscriptionsManager {
	return &subscriptionsManager{
		NewSafeMap(),
		NewSafeMap(),
	}
}
