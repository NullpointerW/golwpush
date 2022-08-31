package netrw

import (
	"github.com/NullpointerW/golwpush/protocol"
	"net"
)

type ReaderBuff struct {
	Buffer  []byte
	WBufPos int
}

func ReadTcp(conn net.Conn, rb *ReaderBuff) (msg string, err error) {
	var (
		retry        bool
		rPos, length int
		tcpErr       error
		jmp          = false
	)
	if rb.WBufPos != 0 {
		rPos = rb.WBufPos
		jmp = true
		goto readBuf
	}
netPull:
	length, tcpErr = conn.Read(rb.Buffer[rb.WBufPos:])
	if tcpErr != nil {
		return msg, tcpErr
	}
	rPos = length + rb.WBufPos
readBuf:
	msg, retry, err = protocol.Unpack(rb.Buffer[:rPos], &rb.WBufPos, jmp)
	if err != nil {
		return msg, err
	}
	if retry {
		goto netPull
	}
	return
}
