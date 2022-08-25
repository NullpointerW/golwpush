package gopush

import (
	"sync/atomic"
)

var (
	connNum   uint64
	connInfos = make(map[uint64]ConnInfo)
)

type ConnInfo struct {
	Id       uint64 `json:"id"`
	Addr     string `json:"addr"`
	ConnTime string `json:"connectTime"`
}

func LoadConnNum() uint64 {
	return atomic.LoadUint64(&connNum)
}

func storeConnNum(num uint64) {
	atomic.StoreUint64(&connNum, num)
}
