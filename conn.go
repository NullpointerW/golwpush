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

	pingCh := make(chan string)
	ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				buf := make([]byte, 128)
				length, _ := tcpConn.Read(buf)
				pingCh <- string(buf[:length])
			}
		}
	}(ctx)

	go func() {
		t := time.NewTimer(time.Minute * 1)
		defer t.Stop()
	Loop:
		for {
			select {
			case <-t.C:
				fmt.Println("Heartbeat timeout 60s...")
				errCh <- errs.HeartbeatTimeout
				break Loop
			case <-pingCh:
				wch <- "pong"
				t.Reset(time.Minute * 1)
			}
		}
	}()

Loop:
	for {
		select {
		case msg := <-wch:
			tcpConn.Write([]byte(msg))
		case err := <-errCh:
			fmt.Println(err.Error())
			conn.close()
			cancel()
			break Loop
		}
	}

}

func newClient(tcpConn net.Conn, id int64) {
	wch := make(chan string, 100)
	errCh := make(chan error, 0)
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
