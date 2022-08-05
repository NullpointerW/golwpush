package GoPush

type defaultPush struct {
	banned map[uint64]struct{}
}

func (p defaultPush) Push(id uint64, msg string) (err error) {
	PushCh <- Content{Id: id, Msg: msg}
	return
}

func (p defaultPush) Broadcast(msg string) (err error) {
	Broadcast <- msg
	return
}

var (
	Default Adapter = defaultPush{banned: nil}
)

type SinglePush interface {
	Push(uint64, string) error
}
type AllPush interface {
	Broadcast(string) error
}

type Adapter interface {
	SinglePush
	AllPush
}
