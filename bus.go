package gopush

import (
	"context"
	"encoding/binary"
	"gopush/errs"
	"gopush/logger"
	"gopush/pkg"
	"gopush/protocol"
	"gopush/utils"
	"net"
	"time"
)

type Content struct {
	Id  uint64
	Msg string
}

func (c Content) pkg() *pkg.Package {
	return &pkg.Package{
		Data: c.Msg,
		Mode: pkg.MSG,
	}
}

var (
	connAddCh0   chan *Conn       = make(chan *Conn)
	ConnAddCh    chan<- *Conn     = connAddCh0
	connRmCh0    chan *Conn       = make(chan *Conn)
	ConnRmCh     chan<- *Conn     = connRmCh0
	broadcast0   chan string      = make(chan string, 1024)
	Broadcast    chan<- string    = broadcast0
	multiPushCh0 chan *Contents   = make(chan *Contents, 1024)
	MultiPushCh  chan<- *Contents = multiPushCh0
	conns        map[uint64]*Conn = make(map[uint64]*Conn)
	pushCh0      chan Content     = make(chan Content, 1024)
	PushCh       chan<- Content   = pushCh0
	bizCh0       chan BizReq      = make(chan BizReq, 1024)
	BizCh        chan<- BizReq    = bizCh0
)

func Handle() {
	for {
		select {
		case content := <-pushCh0:
			if _, exist := conns[content.Id]; exist {
				conns[content.Id].write(content.pkg())
			}
		case conn := <-connAddCh0:
			if _, exist := conns[conn.Id]; exist {
				conn.errMsg <- errs.NewDuplicateConnIdErr(conn.Id)
				continue
			}
			conns[conn.Id] = conn
			now := time.Now().Format(utils.TimeParseLayout)
			storeConnNum(uint64(len(conns)))
			connInfos[conn.Id] = ConnInfo{conn.Id, conn.tcpConn.RemoteAddr().String(), now}
		case conn := <-connRmCh0:
			delete(conns, conn.Id)
			storeConnNum(uint64(len(conns)))
			delete(connInfos, conn.Id)
		case msg := <-broadcast0:
			//logger.Debug(len(conns))
			broadcaster(&pkg.Package{Mode: pkg.MSG,
				Data: msg})
		case contents := <-multiPushCh0:
			multiSend(contents.pkg(), contents.Ids, contents.Res)
		case req := <-bizCh0:
			switch req.Typ {
			case Info:
				if info, exist := connInfos[req.Uid]; exist {
					req.Res <- info
				} else {
					req.Res <- nil
				}
			case Kick:
				//TODO
			}
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
			tcpConn.Close()
			logger.Fatal(errs.SendUidTimeOut)
		}
	}(ctx)
	//接收客户端uid
	data, err := protocol.UnPackByteStream(tcpConn)
	if err != nil {
		logger.PrintfNonUid(logger.CliErr, tcpConn.RemoteAddr().String(), "read error:%v", err)
		cancel()
		return
	}

	cancel()

	uid := binary.BigEndian.Uint64(data)

	logger.PrintfNonUid(logger.Cli|logger.Login|logger.Host, tcpConn.RemoteAddr().String(), "recv uid:%d", uid)

	newClient(tcpConn, uid)
}
