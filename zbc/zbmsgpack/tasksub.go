package zbmsgpack

// TaskSubscription is structure which we use to handle a subscription on a task.
type TaskSubscription struct {
	SubscriberKey uint64 `zbmsgpack:"subscriberKey" json:"subscriberKey"`
	TopicName     string `zbmsgpack:"topicName" json:"topicName"`
	PartitionID   uint16 `zbmsgpack:"partitionId" json:"partitionId"`
	TaskType      string `zbmsgpack:"taskType" json:"taskType"`
	LockDuration  uint64 `zbmsgpack:"lockDuration" json:"lockDuration"`
	LockOwner     string `zbmsgpack:"lockOwner" json:"lockOwner"`
	Credits       int32  `zbmsgpack:"credits" json:"credits"`
}
