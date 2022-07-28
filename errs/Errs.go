package errs

import "errors"

var (
	HeartbeatTimeout = errors.New("heartbeat timeout")
)
