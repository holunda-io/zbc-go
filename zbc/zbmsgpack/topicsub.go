package zbmsgpack

// TopicSubscription is used to open a topic subscription.
type TopicSubscription struct {
	StartPosition    int64  `zbmsgpack:"startPosition"`
	PrefetchCapacity int32  `zbmsgpack:"prefetchCapacity"`
	Name             string `zbmsgpack:"name"`

	ForceStart bool   `zbmsgpack:"forceStart"`
	State      string `zbmsgpack:"state"`
}

// TopicSubscriptionAck is used to acknowledge receiving of an event.
type TopicSubscriptionAck struct {
	Name        string `zbmsgpack:"name"`
	AckPosition uint64 `zbmsgpack:"ackPosition"`
	State       string `zbmsgpack:"state"`
}
