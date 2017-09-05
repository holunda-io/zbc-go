package zbmsgpack

import (
	"fmt"
	//"time"
	"time"
)

type Broker struct {
	Host string `msgpack:"host"`
	Port uint64 `msgpack:"port"`
}

func (t *Broker) Addr() string {
	return fmt.Sprintf("%s:%d", t.Host, t.Port)
}

type TopicLeader struct {
	Broker
	TopicName   string `msgpack:"topicName"`
	PartitionID uint16 `msgpack:"partitionId"`
}

type TopologyRequest struct{}

type ClusterTopologyResponse struct {
	Brokers []Broker `msgpack:"brokers"`
	TopicLeaders []TopicLeader `msgpack:"topicLeaders"`
}

type ClusterTopology struct {
	TopicLeaders map[string]TopicLeader `msgpack:"topicLeaders"`
	Brokers      []Broker               `msgpack:"brokers"`
	UpdatedAt    time.Time              `msgpack:"-"`
}