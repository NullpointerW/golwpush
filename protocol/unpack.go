package protocol

import "GoPush/errs"

const (
	EndFlag byte = '|'
)

func Unpack(b []byte, readIdx int) (msg string, readSt int, retry bool, err error) {

	for i, v := range b {
		if v == EndFlag {
			msg = string(b[:i])
			if len(b) == i+1 {
				readIdx = 0
				return
			}
			readIdx = copy(b, b[i+1:])
			return
		}
	}

	if len(b) == 1024 {
		err = errs.UnpackOutOfSize
		return
	}
	readIdx += len(b)
	retry = true
	return
}
