package zbc

import "sync/atomic"
import "runtime"

// Masking is faster than division
const indexMask uint64 = requestQueueSize - 1

type ringBuffer struct {
	// The padding members 1 to 5 below are here to ensure each item is on a separate cache line.
	// This prevents false sharing and hence improves performance.
	padding1 [8]uint64
	lastCommittedIndex uint64
	padding2 [8]uint64
	nextFreeIndex uint64
	padding3 [8]uint64
	readerIndex uint64
	padding4 [8]uint64
	contents [requestQueueSize]*requestWrapper
	padding5 [8]uint64
}

func newRingBuffer() *ringBuffer {
	return &ringBuffer{lastCommittedIndex : 0, nextFreeIndex : 1, readerIndex : 1}
}

func (rb *ringBuffer) Write(value *requestWrapper) {
	var myIndex = atomic.AddUint64(&rb.nextFreeIndex, 1) - 1
	//Wait for reader to catch up, so we don't clobber a slot which it is (or will be) reading
	for myIndex > (rb.readerIndex + requestQueueSize - 2) {
		runtime.Gosched()
	}
	//Write the item into it's slot
	rb.contents[myIndex & indexMask] = value
	//Increment the lastCommittedIndex so the item is available for reading
	for !atomic.CompareAndSwapUint64(&rb.lastCommittedIndex, myIndex - 1, myIndex) {
		runtime.Gosched()
	}
}

func (rb *ringBuffer) Read() *requestWrapper {
	var myIndex = atomic.AddUint64(&rb.readerIndex, 1) - 1
	//If reader has out-run writer, wait for a value to be committed
	for myIndex > rb.lastCommittedIndex {
		runtime.Gosched()
	}
	return rb.contents[myIndex & indexMask]
}

