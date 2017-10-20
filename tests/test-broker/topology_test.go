package testbroker

import (
	"github.com/zeebe-io/zbc-go/zbc"
	"testing"
)

func TestTopology(t *testing.T) {
	zbClient, err := zbc.NewClient(brokerAddr)
	assert(t, nil, err, true)
	assert(t, nil, zbClient, false)
	assert(t, 0, len(zbClient.Cluster.TopicLeaders), false)
}
