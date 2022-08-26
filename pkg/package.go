package pkg

import (
	"encoding/json"
	"gopush/utils"
)

type Type uint8

const (
	PING Type = iota
	PONG
	MSG
	ERR
	ONLINE
	KICK
	ACK
)

var (
	MaxLen = 1024
)

type Package struct {
	Id   string `json:"id"`
	Mode Type   `json:"mode"`
	Data string `json:"data"`
}

type SendMarshal struct {
	MsgId     string
	Marshaled string
	Mode      Type
}

var (
	Ping             = &Package{Mode: PING}
	PingMarshaled, _ = Ping.MarshalToSend()
	Pong             = &Package{Mode: PONG}
	PongMarshaled, _ = Pong.MarshalToSend()
)

func New(msg string) (*Package, error) {
	pkg := &Package{}
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
	conv.Marshaled = s
	return conv, nil
}
