package pkg

type Type uint8

const (
	PING Type = iota
	MSG
)

type Package struct {
	Mode Type
	Data string
}
