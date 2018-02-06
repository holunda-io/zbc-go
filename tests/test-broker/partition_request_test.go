package testbroker

import (
	"testing"

	"github.com/zeebe-io/zbc-go/zbc"
)

func TestGetPartitions(t *testing.T) {
	zbClient, err := zbc.NewClient(brokerAddr)
	assert(t, nil, err, true)
	assert(t, nil, zbClient, false)
	partitions, err := zbClient.GetPartitions()

	assert(t, nil, err, true)
	assert(t, nil, partitions, false)
}
