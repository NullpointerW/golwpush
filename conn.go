package GoPush

import (
	"GoPush/errs"
	"context"
	"fmt"
	"net"
	"time"
)

type Conn struct {
	Id      int64
	tcpConn net.Conn
	wch     chan<- string
	Addr    string
	errMsg  chan<- error
}

func (conn *Conn) write(msg string) {
	conn.wch <- msg
}

func (conn *Conn) close() {
	close(conn.wch)
	close(conn.errMsg)
	conn.tcpConn.Close()
	ConnRmCh <- conn
}

func connHandle(wch chan string, errCh chan error, id int64, tcpConn net.Conn, conn *Conn) {

	pingCh := make(chan string, 100)
	ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				buf := make([]byte, 128)
				length, err := tcpConn.Read(buf)
				if err != nil {
					errCh <- err
					return
				}
				pingCh <- string(buf[:length])
			}
		}
	}(ctx)

	go func(ctx context.Context) {
		t := time.NewTimer(time.Minute * 1)
		defer t.Stop()
		fmt.Println("start heartbeat check") //debug

		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				fmt.Println("Heartbeat timeout 60s...")
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
			_, err = tcpConn.Write([]byte(msg))
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
	fmt.Println(err.Error())
	conn.close()
	cancelFunc()
}

func newClient(tcpConn net.Conn, id int64) {
	wch := make(chan string, 100)
	errCh := make(chan error, 2)
	conn := &Conn{
		Id:      id,
		tcpConn: tcpConn,
		wch:     wch,
		errMsg:  errCh,
		Addr:    tcpConn.RemoteAddr().String(),
	}
	go connHandle(wch, errCh, id, tcpConn, conn)
	ConnAddCh <- conn
}
