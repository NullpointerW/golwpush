package GoPush

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

type Content struct {
	Id  int64
	Msg string
}

var (
	connAddCh0 chan *Conn      = make(chan *Conn)
	ConnAddCh  chan<- *Conn    = connAddCh0
	connRmCh0  chan *Conn      = make(chan *Conn)
	ConnRmCh   chan<- *Conn    = connRmCh0
	broadcast0 chan string     = make(chan string, 1024)
	Broadcast  chan<- string   = broadcast0
	conns      map[int64]*Conn = make(map[int64]*Conn)
	pushCh0    chan Content    = make(chan Content, 1024)
	PushCh     chan<- Content  = pushCh0
)

func Handle() {
	for {
		select {
		case content := <-pushCh0:
			if _, exist := conns[content.Id]; exist {
				conns[content.Id].write(content.Msg)
			}
		case conn := <-connAddCh0:
			conns[conn.Id] = conn
		case conn := <-connRmCh0:
			delete(conns, conn.Id)
		case msg := <-broadcast0:
			broadcaster(msg)
		}
	}
}

func InitConn(tcpConn net.Conn) {
	buf := make([]byte, 128)

	ctx, cancel := context.WithCancel(context.Background())

	go func(ctx context.Context) {
		t := time.NewTimer(time.Minute * 1)
		defer t.Stop()
		select {
		case <-ctx.Done():
		case <-t.C:
			tcpConn.Close()
		}
	}(ctx)
	length, err := tcpConn.Read(buf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read error:%v", err)
		cancel()
		return
	}
	cancel()
	id, convErr := strconv.ParseInt(string(buf[:length]), 10, 64)
	if convErr != nil {
		fmt.Fprintf(os.Stderr, "parse error:%v", convErr)
		return
	}
	newClient(tcpConn, id)
}
