package testbroker

import (
	"testing"
	"github.com/zeebe-io/zbc-go/zbc"
)

func TestCreateTopic(t *testing.T) {
	zbClient, err := zbc.NewClient(brokerAddr)
	assert(t, nil, err, true)
	assert(t, nil, zbClient, false)

	hash := RandStringBytes(25)
	topic, err := zbClient.CreateTopic(hash, 3)
	assert(t, nil, err, true)
	assert(t, nil, topic, false)

	assert(t, zbc.TopicCreated, topic.State, true)
}

