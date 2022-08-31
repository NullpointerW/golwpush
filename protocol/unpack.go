package protocol

import (
	"encoding/binary"
	"github.com/NullpointerW/golwpush/errs"
	"github.com/NullpointerW/golwpush/logger"
	"github.com/NullpointerW/golwpush/pkg"
	"io"
	"net"
)

const (
	EndFlag byte = '|'
)

func Unpack(b []byte, wIdx *int, jmp bool) (msg string, retry bool, err error) {
	//r:=b[readIdx:]
	logger.Debug("before:" + string(b))
	for i, v := range b {
		if v == EndFlag {
			msg = string(b[:i])
			if len(b) == i+1 {
				*wIdx = 0
				logger.Debug("after:" + string(b))
				return
			}
			*wIdx = copy(b, b[i+1:])
			logger.Debug("after:" + string(b))
			return
		}
	}

	if len(b) == pkg.MaxLen {
		err = errs.UnpackOutOfSize
		return
	}
	if !jmp {
		*wIdx += len(b)
	}
	retry = true
	logger.Debug("after:" + string(b))
	return
}

func UnPackByteStream(conn net.Conn) (data []byte, err error) {
	h := make([]byte, heartLen)
	_, err = io.ReadFull(conn, h)
	if err != nil {
		return
	}
	dataLen := binary.BigEndian.Uint16(h)
	data = make([]byte, dataLen)
	_, err = io.ReadFull(conn, data)
	return
}
