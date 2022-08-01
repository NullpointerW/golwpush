package errs

import (
	"errors"
	"fmt"
)

var (
	HeartbeatTimeout = errors.New("heartbeat timeout")

	UnpackOutOfSize = errors.New("package out of  read buffer size")
)

type duplicateConnIdErr struct {
	s string
}

func (e *duplicateConnIdErr) Error() string {
	return e.s
}

func NewDuplicateConnIdErr(id int64) error {
	err := new(duplicateConnIdErr)
	err.s = fmt.Sprintf("duplicate connection id %d:conn already exisit", id)
	return err
}
