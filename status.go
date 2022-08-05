package GoPush

import "sync/atomic"

var connNum uint64

func loadConnNum() uint64 {
	return atomic.LoadUint64(&connNum)
}

func storeConnNum(num uint64) {
	atomic.StoreUint64(&connNum, num)
}
