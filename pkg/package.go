package pkg

import "encoding/json"

type Type uint8

const (
	PING Type = iota
	PONG
	MSG
	ERR
	OFFLINE
	ONLINE
	KICK
)

var (
	MaxLen = 1024
)

type Package struct {
	Mode Type   `json:"mode"`
	Data string `json:"data"`
}

var (
	Ping = &Package{Mode: PING}
	Pong = &Package{Mode: PONG}
)

func New(msg string) (*Package, error) {
	pkg := &Package{}
	err := json.Unmarshal([]byte(msg), pkg)
	if err != nil {
		return nil, err
	}
	return pkg, nil
}

func (p *Package) ConvStr() (string, error) {
	b, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
