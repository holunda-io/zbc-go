package zbc

import (
	"github.com/vmihailenco/msgpack"
	"github.com/zeebe-io/zbc-go/zbc/zbmsgpack"
	"time"
)

type responseHandler struct{}

func (rf *responseHandler) newClusterTopologyResponse(msg *Message) *zbmsgpack.ClusterTopology {
	var resp zbmsgpack.ClusterTopologyResponse
	msgpack.Unmarshal(msg.Data, &resp)

	ct := &zbmsgpack.ClusterTopology{
		Brokers:      resp.Brokers,
		TopicLeaders: make(map[string][]zbmsgpack.TopicLeader),
		UpdatedAt:    time.Now(),
	}

	for _, leader := range resp.TopicLeaders {
		ct.TopicLeaders[leader.TopicName] = append(ct.TopicLeaders[leader.TopicName], leader)
	}
	return ct
}

func (rf *responseHandler) newCreateTopicResponse(msg *Message) *zbmsgpack.Topic {
	var topic zbmsgpack.Topic
	msgpack.Unmarshal(msg.Data, &topic)

	return &topic
}