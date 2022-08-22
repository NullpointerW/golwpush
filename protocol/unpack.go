package protocol

import (
	"encoding/binary"
	"gopush/errs"
	"gopush/pkg"
	"io"
	"net"
)

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

	if len(b) == pkg.MaxLen {
		err = errs.UnpackOutOfSize
		return
	}
	readIdx += len(b)
	retry = true
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
