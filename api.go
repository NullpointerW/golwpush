package gopush

import (
	"gopush/errs"
	"time"
)

type BizTyp int

const (
	Info BizTyp = iota
	Kick
)

type BizReq struct {
	Res chan any
	Uid uint64
	Typ BizTyp
}

type defaultPush struct {
	banned map[uint64]struct{}
}

func (p defaultPush) Push(id uint64, msg string) (err error) {
	PushCh <- Content{Id: id, Msg: msg}
	return
}

func (p defaultPush) Broadcast(broadMsg string) (err error) {
	Broadcast <- broadMsg
	return
}

func (p defaultPush) MultiPush(cts *Contents) (err error, success uint64) {
	MultiPushCh <- cts
	t := time.NewTicker(time.Second * 10)
	defer t.Stop()
	select {
	case success = <-cts.Res:
	case <-t.C:
		err = errs.ServiceCallTimedOut
	}
	return
}
func (p defaultPush) Count() uint64 {
	return LoadConnNum()
}

func (p defaultPush) Info(uid uint64) (ConnInfo, error) {
	return ConnInfo{}, nil
}

var (
	Default Adapter = defaultPush{banned: nil}
)

type SinglePush interface {
	Push(uint64, string) error
}
type AllPush interface {
	Broadcast(string) error
	MultiPush(cts *Contents) (error, uint64)
}

type ConnManger interface {
	Count() uint64
	Info(uid uint64) (ConnInfo, error)
}

type Adapter interface {
	SinglePush
	AllPush
	ConnManger
}
