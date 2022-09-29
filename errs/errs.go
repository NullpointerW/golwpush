package errs

import (
	"errors"
	"fmt"
)

var (
	HeartbeatTimeout = errors.New("heartbeat timeout")

	UnpackOutOfSize = errors.New("package out of read buffer size")

	SendDataOutOfSize = errors.New("package out of max size")

	SendUidTimeOut = errors.New("send uid timeout")

	ServiceCallTimedOut = errors.New("service has not responded for a long time")

	AckBuffCapLimit = errors.New("ackBuffer capacity is full")
)

type DuplicateConnIdErr struct {
	s string
}

func (e DuplicateConnIdErr) Error() string {
	return e.s
}

func NewDuplicateConnIdErr(id uint64) error {
	err := DuplicateConnIdErr{}
	err.s = fmt.Sprintf("duplicate connection id %d:already exisit", id)
	return err
}
