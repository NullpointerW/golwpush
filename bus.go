package GoPush

import (
	"GoPush/errs"
	"GoPush/logger"
	"GoPush/protocol"
	"context"
	"encoding/binary"
	"net"
	"time"
)

type Content struct {
	Id  uint64
	Msg string
}

var (
	connAddCh0 chan *Conn       = make(chan *Conn)
	ConnAddCh  chan<- *Conn     = connAddCh0
	connRmCh0  chan *Conn       = make(chan *Conn)
	ConnRmCh   chan<- *Conn     = connRmCh0
	broadcast0 chan string      = make(chan string, 1024)
	Broadcast  chan<- string    = broadcast0
	conns      map[uint64]*Conn = make(map[uint64]*Conn)
	pushCh0    chan Content     = make(chan Content, 1024)
	PushCh     chan<- Content   = pushCh0
)

func Handle() {
	for {
		select {
		case content := <-pushCh0:
			if _, exist := conns[content.Id]; exist {
				conns[content.Id].write(content.Msg)
			}
		case conn := <-connAddCh0:
			if _, exist := conns[conn.Id]; exist {
				conn.errMsg <- errs.NewDuplicateConnIdErr(conn.Id)
			}
			conns[conn.Id] = conn
		case conn := <-connRmCh0:
			delete(conns, conn.Id)
		case msg := <-broadcast0:
			broadcaster(msg)
		}
	}
}

func InitConn(tcpConn net.Conn) {
	//buf := make([]byte, 128)

	ctx, cancel := context.WithCancel(context.Background())

	go func(ctx context.Context) {
		t := time.NewTimer(time.Minute * 1)
		defer t.Stop()
		select {
		case <-ctx.Done():
		case <-t.C:
			logger.Error(errs.SendUidTimeOut)
			tcpConn.Close()
		}
	}(ctx)
	data, err := protocol.UnPackByteStream(tcpConn)
	if err != nil {
		logger.Errorf("read error:%v", err)
		cancel()
		return
	}
	cancel()
	id := binary.BigEndian.Uint64(data)
	//	var (
	//		length = 0
	//	)
	//Loop:
	//	for {
	//		var (
	//			err error
	//		)
	//		for {
	//			var l int
	//			l, err = tcpConn.Read(buf[length:])
	//			length += l
	//			if err != nil {
	//				goto Fatal
	//			}
	//			if buf[length-1] == protocol.EndFlag {
	//				break Loop
	//			}
	//			if length >= len(buf) {
	//				goto Fatal
	//			}
	//		}
	//	Fatal:
	//		logger.Errorf("read error:%v", err)
	//		cancel()
	//		tcpConn.Close()
	//		return
	//	}
	logger.Debugf("recv id:%d", id)
	//cancel()
	//id, convErr := strconv.ParseInt(string(buf[:length-1]), 10, 64)
	//if convErr != nil {
	//	logger.Errorf("parse error:%v", convErr)
	//	return
	//}
	newClient(tcpConn, id)
}
