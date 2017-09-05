package zbc

import (
	"github.com/zeebe-io/zbc-go/zbc/zbmsgpack"
	"github.com/vmihailenco/msgpack"
	"time"
)

type ResponseHandler struct {}

func (rf *ResponseHandler) newClusterTopologyResponse(msg *Message) *zbmsgpack.ClusterTopology {
	var resp zbmsgpack.ClusterTopologyResponse
	msgpack.Unmarshal(msg.Data, &resp)

	ct := &zbmsgpack.ClusterTopology{
		Brokers: resp.Brokers,
		TopicLeaders: make(map[string][]zbmsgpack.TopicLeader),
		UpdatedAt: time.Now(),
	}

	for _, leader := range resp.TopicLeaders {
		ct.TopicLeaders[leader.TopicName] = append(ct.TopicLeaders[leader.TopicName], leader)
	}
	return ct
}
