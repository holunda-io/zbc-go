package zbmsgpack

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

// BrokerPartition contains information about partition contained in topology request.
type BrokerPartition struct {
	State       string `msgpack:"state"`
	TopicName   string `msgpack:"topicName"`
	PartitionID uint16 `msgpack:"partitionId"`
}

// Broker is used to hold broker contact information.
type Broker struct {
	Host       string            `msgpack:"host"`
	Port       uint64            `msgpack:"port"`
	Partitions []BrokerPartition `msgpack:"partitions"`
}

func (t *Broker) String() string {
	b, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return fmt.Sprintf("json marshaling failed\n")
	}
	return fmt.Sprintf("%+v", string(b))
}

func (t *Broker) Addr() string {
	return fmt.Sprintf("%s:%d", t.Host, t.Port)
}

// TopicLeader is used to hold information about the relation of broker and topic/partitions.
// // type TopicLeader struct {
// // 	Broker
// // }

// func (t *TopicLeader) String() string {
// 	b, err := json.MarshalIndent(t, "", "  ")
// 	if err != nil {
// 		return fmt.Sprintf("json marshaling failed\n")
// 	}
// 	return fmt.Sprintf("%+v", string(b))
// }

// TopologyRequest is used to make a topology request.
type TopologyRequest struct{}

func (t *TopologyRequest) String() string {
	b, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return fmt.Sprintf("json marshaling failed\n")
	}
	return fmt.Sprintf("%+v", string(b))
}

// ClusterTopologyResponse is used to parse the topology response from the broker.
type ClusterTopologyResponse struct {
	Brokers []Broker `msgpack:"brokers"`
}

func (t *ClusterTopologyResponse) String() string {
	b, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return fmt.Sprintf("json marshaling failed\n")
	}
	return fmt.Sprintf("%+v", string(b))
}

// ClusterTopology is structure used by the client object to hold information about cluster.
type ClusterTopology struct {
	AddrByPartitionID      map[uint16]string
	PartitionIDByTopicName map[string][]uint16

	Partitions *PartitionCollection
	Brokers    []Broker
	UpdatedAt  time.Time
}

// GetRandomBroker will select one broker on random using time.Now().UTC().UnixNano() as the seed for randomness.
func (t *ClusterTopology) GetRandomBroker() Broker {
	rand.Seed(time.Now().UTC().UnixNano())
	index := rand.Intn(len(t.Brokers))
	broker := t.Brokers[index] //Brokers[rand.Int()%len(tm.cluster.Brokers)].Addr()
	return broker
}

// String will marshal the data structure as series of characters.
func (t *ClusterTopology) String() string {
	b, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return fmt.Sprintf("json marshaling failed\n")
	}
	return fmt.Sprintf("%+v", string(b))
}
