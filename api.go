package golwpush

import (
	"github.com/NullpointerW/golwpush/errs"
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
	select {
	case Broadcast <- broadMsg: //限流

	default:
		// TODO log
	}

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

func (p defaultPush) Info(biz BizReq) (*ConnInfo, error) {
	t := time.NewTicker(time.Second * 10)
	defer t.Stop()
	BizCh <- biz
	select {
	case r := <-biz.Res:
		if info, exist := r.(ConnInfo); exist {
			return &info, nil
		}
		return nil, nil //offline
	case <-t.C:
		err := errs.ServiceCallTimedOut
		return nil, err
	}

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
	Info(biz BizReq) (*ConnInfo, error)
}

type Adapter interface {
	SinglePush
	AllPush
	ConnManger
}
