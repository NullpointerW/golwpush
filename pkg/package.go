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

var (
	Ping             = &Package{Mode: PING}
	PingMarshaled, _ = Ping.Marshal()
	Pong             = &Package{Mode: PONG}
	PongMarshaled, _ = Pong.Marshal()
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
