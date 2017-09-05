package zbc

import "time"

// RequestTimeout specifies default timeout for Responder.
const RequestTimeout = 5

// TopologyRefreshInterval defines time to live of topology object.
const TopologyRefreshInterval = 30

// Retry constants
const (
	BackoffMin      = 1 * time.Millisecond
	BackoffMax      = 100 * time.Millisecond
	BackoffDeadline = 10 * time.Second
)


// Sbe template ID constants
const (
	templateIDExecuteCommandRequest  = 20
	templateIDExecuteCommandResponse = 21
	templateIDControlMessageResponse = 11
	templateIDSubscriptionEvent      = 30
)

// Zeebe protocol constants
const (
	FrameHeaderSize           = 12
	TransportHeaderSize       = 2
	RequestResponseHeaderSize = 8
	SBEMessageHeaderSize      = 8

	TotalHeaderSizeNoFrame = 18
	LengthFieldSize        = 2
)

// Message pack states
const (
	CreateDeployment = "CREATE_DEPLOYMENT"
	WorkflowInstanceRejected = "WORKFLOW_INSTANCE_REJECTED"
)

const (
	TopicSubscriptionSubscribeState  = "SUBSCRIBE"
	TopicSubscriptionSubscribedState = "SUBSCRIBED"
)

const (
	TopicSubscriptionAckState          = "ACKNOWLEDGE"
	TopicSubscriptionAcknowledgedState = "ACKNOWLEDGED"
)

