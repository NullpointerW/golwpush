package netrw

import (
	"github.com/NullpointerW/golwpush/protocol"
	"net"
)

type Reader interface {
	Read() (msg string, err error)
}

type TcpReader struct {
	Buf  []byte
	WPos int
	net.Conn
}

func (r *TcpReader) Read() (msg string, err error) {
	var (
		retry        bool
		rPos, length int
		tcpErr       error
		jmp          bool
	)
	if r.WPos != 0 {
		rPos = r.WPos
		jmp = true
		goto bufPull
	}
netPull:
	length, tcpErr = r.Conn.Read(r.Buf[r.WPos:])
	if tcpErr != nil {
		return msg, tcpErr
	}
	rPos = length + r.WPos
bufPull:
	msg, retry, err = protocol.Unpack(r.Buf[:rPos], &r.WPos, jmp)
	if err != nil {
		return msg, err
	}
	if retry {
		goto netPull
	}
	return
}
