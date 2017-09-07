package zbc

import (
	"github.com/zeebe-io/zbc-go/zbc/zbsbe"
	"github.com/zeebe-io/zbc-go/zbc/zbmsgpack"
)

// TaskEvent structure is used on task subscription channel.
type TaskEvent struct {
	*zbmsgpack.Task
	Event *zbsbe.SubscribedEvent
}
