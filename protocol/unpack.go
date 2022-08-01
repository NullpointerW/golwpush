package protocol

import "GoPush/errs"

const (
	EndFlag byte = '|'
)

func Unpack(b []byte) (msg string, readIdx int, err error) {

	readIdx = 0

	for i, v := range b {
		if v == EndFlag {
			msg = string(b[:i])
			if len(b) == i+1 {
				return
			}
			readIdx = copy(b, b[i+1:])
			return
		}
	}

	if len(b) == 128 {
		return msg, readIdx, errs.UnpackOutOfSize
	}
	readIdx = len(b)
	return
}
