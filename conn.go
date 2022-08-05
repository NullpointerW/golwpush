package GoPush

import (
	"GoPush/errs"
	"GoPush/logger"
	"GoPush/pkg"
	"GoPush/protocol"
	"context"
	"net"
	"time"
)

type Conn struct {
	Id         uint64
	tcpConn    net.Conn
	readBuf    []byte
	readBufPtr int
	wch        chan<- string
	Addr       string
	errMsg     chan<- error
}

func (conn *Conn) write(msg string) {
	conn.wch <- msg
}

func (conn *Conn) read() (msg string, err error) {
	length, TCPErr := conn.tcpConn.Read(conn.readBuf[conn.readBufPtr:])
	if TCPErr != nil {
		return msg, TCPErr
	}
	var retry bool
unpack:
	msg, conn.readBufPtr, retry, err = protocol.Unpack(conn.readBuf[:length+conn.readBufPtr], conn.readBufPtr)
	if retry {
		goto unpack
	}
	if err != nil {
		return msg, err
	}
	return
}

func (conn *Conn) close() {
	//close(conn.wch)
	//close(conn.errMsg)
	conn.tcpConn.Close()
	ConnRmCh <- conn
}

func connHandle(wch chan string, errCh chan error, id uint64, tcpConn net.Conn, conn *Conn) {

	pingCh := make(chan string, 100)
	ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				msg, err := conn.read()
				if err != nil {
					//可能导致写入已关闭channel panic
					//解决方法有两种：
					//1.不关闭channel 由gc回收 连接过多时可能会导致效率下降
					//2.使用mutex维护一个关闭状态
					errCh <- err
					return
				}
				pingCh <- msg
			}
		}
	}(ctx)

	go func(ctx context.Context) {
		t := time.NewTimer(time.Minute * 1)
		defer t.Stop()
		logger.Debug("start heartbeat check") //debug

		for {
			select {
			case <-ctx.Done():
				logger.Debug("heartbeat check end") //debug
				return
			case <-t.C:
				logger.Warn("Heartbeat timeout 60s...")
				errCh <- errs.HeartbeatTimeout
				return
			case <-pingCh:
				wch <- "pong"
				t.Reset(time.Minute * 1)
			}
		}
	}(ctx)

	var (
		err error
	)

	for {
		select {
		case msg := <-wch:
			_, err = tcpConn.Write(protocol.Pack(msg))
			if err != nil {
				goto Fatal
			}
		case err = <-errCh:
			goto Fatal
		}
	}
Fatal:
	connFatal(err, conn, cancel)

}
func connFatal(err error, conn *Conn, cancelFunc context.CancelFunc) {
	logger.Error(err)
	conn.close()
	cancelFunc()
}

func newClient(tcpConn net.Conn, id uint64) {
	wch := make(chan string, 100)
	errCh := make(chan error, 3)
	conn := &Conn{
		Id:         id,
		tcpConn:    tcpConn,
		readBuf:    make([]byte, pkg.MaxLen),
		readBufPtr: 0,
		wch:        wch,
		errMsg:     errCh,
		Addr:       tcpConn.RemoteAddr().String(),
	}
	go connHandle(wch, errCh, id, tcpConn, conn)
	ConnAddCh <- conn
}
