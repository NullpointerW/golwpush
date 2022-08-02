package errs

import (
	"errors"
	"fmt"
)

var (
	HeartbeatTimeout = errors.New("heartbeat timeout")

	UnpackOutOfSize = errors.New("package out of  read buffer size")

	SendUidTimeOut = errors.New("send uid timeout")
)

type duplicateConnIdErr struct {
	s string
}

func (e *duplicateConnIdErr) Error() string {
	return e.s
}

func NewDuplicateConnIdErr(id int64) error {
	err := new(duplicateConnIdErr)
	err.s = fmt.Sprintf("duplicate connection id %d:already exisit", id)
	return err
}
