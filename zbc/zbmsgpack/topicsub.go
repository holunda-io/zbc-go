package zbmsgpack

// OpenTopicSubscription is used to open a topic subscription.
type OpenTopicSubscription struct {
	StartPosition    int64  `msgpack:"startPosition"`
	PrefetchCapacity int32  `msgpack:"prefetchCapacity"`
	Name             string `msgpack:"name"`

	ForceStart bool   `msgpack:"forceStart"`
	State      string `msgpack:"state"`
}

// TopicSubscriptionAck is used to acknowledge receiving of an event.
type TopicSubscriptionAck struct {
	Name        string `msgpack:"name"`
	AckPosition uint64 `msgpack:"ackPosition"`
	State       string `msgpack:"state"`
}

type TopicSubscription struct {
	TopicName     string `msgpack:"topicName"`
	PartitionID   uint16 `msgpack:"partitionId"`
	SubscriberKey uint64 `msgpack:"subscriberKey"`
}
