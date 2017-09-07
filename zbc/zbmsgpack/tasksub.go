package zbmsgpack

// TaskSubscription is structure which we use to handle a subscription on a task.
type TaskSubscription struct {
	SubscriberKey uint64 `msgpack:"subscriberKey" json:"subscriberKey"`
	TopicName     string `msgpack:"topicName" json:"topicName"`
	PartitionID   uint16 `msgpack:"partitionId" json:"partitionId"`
	TaskType      string `msgpack:"taskType" json:"taskType"`
	LockDuration  uint64 `msgpack:"lockDuration" json:"lockDuration"`
	LockOwner     string `msgpack:"lockOwner" json:"lockOwner"`
	Credits       int32  `msgpack:"credits" json:"credits"`
}
