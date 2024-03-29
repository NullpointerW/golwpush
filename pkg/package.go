package pkg

import (
	"encoding/json"
	"github.com/NullpointerW/golwpush/protocol"
	"github.com/NullpointerW/golwpush/utils"
)

type Typ uint8

const (
	PING Typ = iota
	PONG
	MSG
	ERR
	ONLINE
	KICK
	ACK

	MsgModeExtraLen = 90
)

type Package struct {
	Uid  string          `json:"uid,omitempty"`
	Id   string          `json:"id,omitempty"`
	Mode Typ             `json:"mode"`
	Data json.RawMessage `json:"data,omitempty"` // ACK
}

type SendMarshal struct {
	//Uid       uint64 `json:"-"`
	MsgId     string `json:"id"`
	Marshaled string `json:"marshaled"`
	Mode      Typ    `json:"mode"`
}

var (
	Ping             = &Package{Mode: PING}
	PingMarshaled, _ = Ping.MarshalToSend()
	Pong             = &Package{Mode: PONG}
	PongMarshaled, _ = Pong.MarshalToSend()
)

func New(msg string) (*Package, error) {
	pkg := &Package{}
	// fmt.Println("json_raw is"+msg)
	err := json.Unmarshal(utils.Scb(msg), pkg)
	if err != nil {
		return nil, err
	}
	return pkg, nil
}

func (p *Package) Marshal() (string, error) {
	b, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return utils.Bcs(b), nil
}

func (p *Package) MarshalToSend() (SendMarshal, error) {
	conv := SendMarshal{Mode: p.Mode}
	s, err := p.Marshal()
	if err != nil {
		return conv, err
	}
	conv.Marshaled = protocol.CatEndFlag(s)
	conv.MsgId = p.Id
	return conv, nil
}
