package zbmsgpack

type PartitionDetails struct {
	ID    uint   `msgpack:"id"`
	Topic string `msgpack:"topic"`
}

type PartitionCollection struct {
	Partitions []PartitionDetails `msgpack:"partitions"`
}
