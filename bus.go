package golwpush

import (
	"context"
	"encoding/binary"
	"github.com/NullpointerW/golwpush/errs"
	"github.com/NullpointerW/golwpush/logger"
	"github.com/NullpointerW/golwpush/pkg"
	"github.com/NullpointerW/golwpush/protocol"
	"github.com/NullpointerW/golwpush/utils"
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
	broadcast0   chan string      = make(chan string, 2048)
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
	lingerMs := time.NewTicker(time.Millisecond * 100)
	defer lingerMs.Stop()
	for {
		select {
		case <-lingerMs.C:
			lingerSend()

		case content := <-pushCh0:
			select {
			case <-lingerMs.C:
				lingerSend()
			default:
			}
			if _, exist := conns[content.Id]; exist {
				conns[content.Id].write(content.pkg())
			}

		case conn := <-connAddCh0:
			select {
			case <-lingerMs.C:
				lingerSend()
			default:
			}
			if _, exist := conns[conn.Uid]; exist {
				conn.errMsg <- errs.NewDuplicateConnIdErr(conn.Uid)
				continue
			}
			conns[conn.Uid] = conn
			now := time.Now().Format(utils.TimeParseLayout)
			storeConnNum(uint64(len(conns)))
			connInfos[conn.Uid] = ConnInfo{conn.Uid, conn.tcpConn.RemoteAddr().String(), now}

		case conn := <-connRmCh0:
			select {
			case <-lingerMs.C:
				lingerSend()
			default:
			}
			delete(conns, conn.Uid)
			storeConnNum(uint64(len(conns)))
			delete(connInfos, conn.Uid)

		case msg := <-broadcast0:
			select {
			case <-lingerMs.C:
				lingerSend()
			default:
			}
			mergeMsg(msg)
			//broadcaster(&pkg.Package{Mode: pkg.MSG,
			//	Data: msg})
		case contents := <-multiPushCh0:
			select {
			case <-lingerMs.C:
				lingerSend()
			default:
			}
			multiSend(contents.pkg(), contents.Ids, contents.Res)

		case req := <-bizCh0:
			select {
			case <-lingerMs.C:
				lingerSend()
			default:
			}
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
	uid, err := AuthCli(tcpConn)
	if err != nil {
		logger.Error(err)
		return
	}
	newClient(tcpConn, uid)
}

func AuthCli(conn net.Conn) (uid uint64, err error) {
	ctx, cancel := context.WithCancel(context.Background())

	go func(ctx context.Context) {
		t := time.NewTimer(time.Minute * 1)
		defer t.Stop()
		select {
		case <-ctx.Done():
		case <-t.C:
			conn.Close()
			logger.PlnNUid(logger.L_Err|logger.Host, conn.RemoteAddr().String(), errs.SendUidTimeOut)
		}
	}(ctx)
	//接收客户端uid
	data, err := protocol.UnPackByteStream(conn)
	if err != nil {
		logger.PfNUid(logger.CliErr, conn.RemoteAddr().String(), "read error:%v", err)
		cancel()
		return
	}
	cancel()
	uid = binary.BigEndian.Uint64(data)
	logger.PfNUid(logger.Cli|logger.Login|logger.Host, conn.RemoteAddr().String(), "recv uid:%d", uid)
	return
}
