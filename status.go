package GoPush

import (
	"sync/atomic"
)

var connNum uint64

type ConnInfo struct {
	Id          uint64 `json:"id"`
	Addr        string `json:"addr"`
	ConnectTime string `json:"connectTime"`
}

func LoadConnNum() uint64 {
	return atomic.LoadUint64(&connNum)
}

func storeConnNum(num uint64) {
	atomic.StoreUint64(&connNum, num)
}
