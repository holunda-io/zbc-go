package zbc

import "time"

// RequestTimeout specifies default timeout for Responder.
const RequestTimeout = 5

// TopologyRefreshInterval defines time to live of topology object.
const TopologyRefreshInterval = 30

// MaxRetries specifies total number of retries for failed requests.
const MaxRetries = 5

// RetryDeadline specified total number of seconds before we give up on retrying on requests.
const RetryDeadline = 15

const (
	templateIDExecuteCommandRequest  = 20
	templateIDExecuteCommandResponse = 21
	templateIDControlMessageResponse = 11
	templateIDSubscriptionEvent      = 30
)

const (
	FrameHeaderSize           = 12
	TransportHeaderSize       = 2
	RequestResponseHeaderSize = 8
	SBEMessageHeaderSize      = 8

	TotalHeaderSizeNoFrame = 18
	LengthFieldSize        = 2
)

const CreateDeployment = "CREATE_DEPLOYMENT"
const WorkflowInstanceRejected = "WORKFLOW_INSTANCE_REJECTED"

const (
	TopicSubscriptionSubscribeState  = "SUBSCRIBE"
	TopicSubscriptionSubscribedState = "SUBSCRIBED"
)

const (
	TopicSubscriptionAckState          = "ACKNOWLEDGE"
	TopicSubscriptionAcknowledgedState = "ACKNOWLEDGED"
)

const (
	BackoffMin      = 1 * time.Millisecond
	BackoffMax      = 100 * time.Millisecond
	BackoffDeadline = 10 * time.Second
)
