package zbmsgpack

// TaskSubscription is structure which we use to handle a subscription on a task.
type TaskSubscription struct {
	SubscriberKey uint64 `zbmsgpack:"subscriberKey"`
	TopicName     string `zbmsgpack:"topicName"`
	PartitionID   uint16  `zbmsgpack:"partitionId"`
	TaskType      string `zbmsgpack:"taskType"`
	LockDuration  uint64 `zbmsgpack:"lockDuration"`
	LockOwner     string `zbmsgpack:"lockOwner"`
	Credits       int32  `zbmsgpack:"credits"`
}