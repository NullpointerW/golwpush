package errs

import "errors"

var (
	HeartbeatTimeout = errors.New("heartbeat timeout")

	UnpackOutOfSize = errors.New("package out of  read buffer size")
)
