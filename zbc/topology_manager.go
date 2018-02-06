package zbc

import (
	"errors"
	"math/rand"

	"github.com/zeebe-io/zbc-go/zbc/zbmsgpack"

	"time"
)

var (
	errPartitionNotFound = errors.New("partition not found")
	errNoBrokersFound    = errors.New("no brokers found")
)

type topologyManager struct {
	*transportManager

	topologyWorkload chan *requestWrapper

	lastIndexes map[string]uint16

	bootstrapAddr string
	cluster       *zbmsgpack.ClusterTopology
}

func (tm *topologyManager) topicPartitionsAddrs(topic string) (*map[uint16]string, error) {
	addrs := make(map[uint16]string)
	if partitions, ok := tm.cluster.PartitionIDByTopicName[topic]; ok {
		for _, partitionID := range partitions {
			if addr, ok := tm.cluster.AddrByPartitionID[partitionID]; ok {
				addrs[partitionID] = addr
			} else {
				return nil, errPartitionNotFound
			}

		}
	}

	return &addrs, nil
}

func (tm *topologyManager) partitionID(topic string) (uint16, error) {
	lastPartitionUsedIndex, ok := tm.lastIndexes[topic]
	if !ok {
		lastPartitionUsedIndex = 0
	}
	brokers, ok := tm.cluster.PartitionIDByTopicName[topic]
	if !ok {
		tm.refreshTopology()

		brokers, ok = tm.cluster.PartitionIDByTopicName[topic]

		if !ok {
			return 0, errTopicLeaderNotFound
		}
	}

	var partitionID uint16
	if tm.lastIndexes[topic] >= uint16(len(brokers)) {
		tm.lastIndexes[topic] = 0
	} else {
		partitionID = brokers[lastPartitionUsedIndex]
		tm.lastIndexes[topic] = lastPartitionUsedIndex + 1
	}

	// TODO: zbc-go/issues#40 + zbc-go/issues#48
	return partitionID, nil
}

func (tm *topologyManager) initTopology() (*zbmsgpack.ClusterTopology, error) {
	factory := newRequestFactory()
	responseHandler := newResponseHandler()

	resp, err := tm.executeRequest(newRequestWrapper(factory.topologyRequest()))
	topology := responseHandler.unmarshalTopology(resp)
	if err != nil {
		return nil, err
	}

	tm.cluster = &topology
	return tm.cluster, nil
}

func (tm *topologyManager) refreshTopology() (*zbmsgpack.ClusterTopology, error) {
	rand.Seed(time.Now().Unix())

	factory := newRequestFactory()
	responseHandler := newResponseHandler()

	request := newRequestWrapper(factory.topologyRequest())
	if len(tm.cluster.Brokers) == 0 {
		return nil, errNoBrokersFound
	}
	broker := tm.cluster.GetRandomBroker()
	request.addr = broker.Addr() //.Brokers[rand.Int()%len(tm.cluster.Brokers)].Addr() // Get random broker
	resp, err := tm.executeRequest(request)
	topology := responseHandler.unmarshalTopology(resp)
	if err != nil {
		return nil, err
	}

	tm.cluster = &topology
	return tm.cluster, err
}

func (tm *topologyManager) getDestinationAddr(msg *Message) (string, error) {
	if msg.isTopologyMessage() && tm.cluster == nil {
		return tm.bootstrapAddr, nil
	}

	partitionID := msg.forPartitionId()
	if partitionID == nil {
		return "", brokerNotFound
	}

	if addr, ok := tm.cluster.AddrByPartitionID[*partitionID]; ok {
		return addr, nil
	}

	return "", brokerNotFound
}

func (tm *topologyManager) executeRequest(request *requestWrapper) (*Message, error) {
	addr, err := tm.getDestinationAddr(request.payload)

	if err == brokerNotFound {
		return nil, brokerNotFound
	}
	request.addr = addr
	tm.topologyWorkload <- request

	select {

	case resp := <-request.responseCh:
		return resp, nil

	case err := <-request.errorCh:
		if tm.cluster != nil {
			tm.refreshTopology()
		}
		return nil, err

	}
}

func (tm *topologyManager) topologyTicker() {
	for {
		select {
		case <-time.After(TopologyRefreshInterval * time.Second):
			if time.Since(tm.cluster.UpdatedAt) > TopologyRefreshInterval*time.Second {
				tm.topologyWorkload <- newRequestWrapper(newRequestFactory().topologyRequest())
			}
			break
		}
	}
}

func (tm *topologyManager) topologyWorker() {
	for {
		select {
		case request := <-tm.topologyWorkload:
			tm.execTransport(request)
		}
	}
}

func newTopologyManager(bootstrapAddr string) *topologyManager {
	tm := &topologyManager{
		newTransportManager(),
		make(chan *requestWrapper, requestQueueSize),
		make(map[string]uint16),
		bootstrapAddr,
		nil,
	}

	go tm.topologyWorker()
	go tm.topologyTicker()

	tm.initTopology()
	return tm
}
