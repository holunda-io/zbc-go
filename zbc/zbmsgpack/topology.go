package zbmsgpack

import (
	"fmt"
	//"time"
	"time"
)

// Broker is used to hold broker contact information.
type Broker struct {
	Host string `msgpack:"host"`
	Port uint64 `msgpack:"port"`
}

func (t *Broker) Addr() string {
	return fmt.Sprintf("%s:%d", t.Host, t.Port)
}

// TopicLeader is used to hold information about the relation of broker and topic/partitions.
type TopicLeader struct {
	Broker
	TopicName   string `msgpack:"topicName"`
	PartitionID uint16 `msgpack:"partitionId"`
}

// TopologyRequest is used to make a topology request.
type TopologyRequest struct{}

// ClusterTopologyResponse is used to parse the topology response from the broker.
type ClusterTopologyResponse struct {
	Brokers      []Broker      `msgpack:"brokers"`
	TopicLeaders []TopicLeader `msgpack:"topicLeaders"`
}

// ClusterTopology is structure used by the client object to hold information about cluster.
type ClusterTopology struct {
	TopicLeaders map[string][]TopicLeader `msgpack:"topicLeaders"`
	Brokers      []Broker                 `msgpack:"brokers"`
	UpdatedAt    time.Time                `msgpack:"-"`
}
