package zbc

type requestWrapper struct {
	addr string
	sock       *socket
	responseCh chan *Message
	errorCh    chan error
	payload    *Message
}

func newRequestWrapper(payload *Message) *requestWrapper {
	return &requestWrapper{
		"",
		nil,
		make(chan *Message),
		make(chan error),
		payload,
	}
}
