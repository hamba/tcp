package tcp

import "sync/atomic"

type atomicBool int32

func (b *atomicBool) isSet() bool {
	return atomic.LoadInt32((*int32)(b)) != 0
}

func (b *atomicBool) set() {
	atomic.StoreInt32((*int32)(b), 1)
}
