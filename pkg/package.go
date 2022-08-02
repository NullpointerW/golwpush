package pkg

type Type uint8

const (
	PING Type = iota
	MSG
	ERR
	OFFLINE
	ONLINE
	KICK
)

type Package struct {
	Mode Type   `json:"mode"`
	Data string `json:"data"`
}