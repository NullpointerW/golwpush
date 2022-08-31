package golwpush

import (
	"context"
	"github.com/NullpointerW/golwpush/errs"
	"github.com/NullpointerW/golwpush/logger"
	"github.com/NullpointerW/golwpush/netrw"
	"github.com/NullpointerW/golwpush/pkg"
	"github.com/NullpointerW/golwpush/protocol"
	"github.com/NullpointerW/golwpush/utils"
	"net"
	"strconv"
	"sync/atomic"
	"time"
)

type Conn struct {
	Id      uint64
	tcpConn net.Conn
	/*readBuf []byte
	wBufPos int*/
	netrw.ReaderBuff
	wch     chan<- pkg.SendMarshal
	Addr    Addr
	errMsg  chan<- error
	sendSeq atomic.Value
}

type ackPeek struct {
	ack            context.CancelFunc
	actualSendTime time.Time
}
type Addr interface {
	net.Addr
	Uid() uint64
}
type ConnAddr struct {
	net.Addr
	uid uint64
}

func (c *ConnAddr) Uid() uint64 {
	return c.uid
}

func (conn *Conn) write(msg *pkg.Package) {
	msg.Id = utils.GenerateId(conn.Id) + strconv.FormatUint(conn.incrSeq(), 10) //生成消息id
	marshaled, err := msg.MarshalToSend()                                       //避免在写goroutine中编码
	if err != nil {
		conn.errMsg <- err
		return
	}
	conn.wch <- marshaled
}

func (conn *Conn) read() (msg string, err error) {
	return netrw.ReadTcp(conn.tcpConn, &conn.ReaderBuff)
}

func (conn *Conn) close() {
	//close(conn.wch)
	//close(conn.errMsg)
	logger.Debug("close conn" + conn.Addr.String())
	err := conn.tcpConn.Close()
	if err != nil {
		logger.Error(err)
	}
	ConnRmCh <- conn
}

func (conn *Conn) incrSeq() uint64 { //cas 获取发送序列
	for {
		o := conn.sendSeq.Load().(uint64)
		n := o + 1
		if swa := conn.sendSeq.CompareAndSwap(o, n); swa {
			return n
		}
	}

}

func connHandle(wch chan pkg.SendMarshal, errCh chan error, uid uint64, tcpConn net.Conn, conn *Conn) {
	ackBuf, ackBuf0 := utils.NewChMap[string, ackPeek](10000) //max 10000
	pingCh := make(chan struct{}, 100)
	ctx, cancel := context.WithCancel(context.Background())
	rHandler := make(chan string, 1024)
	//避免在读goroutine解码，通过一个goroutine处理所有读到的包
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case unDec := <-rHandler:
				p, err := pkg.New(unDec)
				if err != nil {
					errCh <- err
				}
				switch p.Mode {
				case pkg.ACK:
					ackBuf.Del <- p.Id
					logger.Infof("ack from cli %s,msg id:%s", conn.Addr.String(), p.Id)
				case pkg.PING:
					pingCh <- struct{}{}
				}
			}
		}
	}(ctx)

	go func(ctx context.Context) { //readLoop
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
				rHandler <- msg
			}
		}
	}(ctx)

	go func(ctx context.Context) {
		t := time.NewTimer(time.Minute * 1)
		defer t.Stop()
		logger.PrintlnWithAddr(logger.L_Debug|logger.HeartBeat|logger.Cli, conn.Addr.Uid(), conn.Addr.String(),
			"start heartbeat check") //debug
		for {
			select {
			case <-ctx.Done():
				logger.PrintlnWithAddr(logger.L_Debug|logger.PingOutput, conn.Addr.Uid(), conn.Addr.String(),
					"heartbeat check end") //debug
				return
			case <-t.C:
				logger.PrintlnWithAddr(logger.L_Warn|logger.HeartBeat|logger.Cli, conn.Addr.Uid(), conn.Addr.String(),
					"Heartbeat timeout 60s...") //debug
				errCh <- errs.HeartbeatTimeout
				return
			case <-pingCh:
				wch <- pkg.PongMarshaled
				t.Reset(time.Minute * 1)
			}
		}
	}(ctx)

	var (
		err error
	)
	//write loop
	for {
		select {
		case msg := <-wch:
			if msg.Mode == pkg.MSG {
				//TODO ACK
				pending, ack := context.WithCancel(context.Background())
				//ackBuf0[msg.MsgId] = ackPeek{
				//	ack:            ack,
				//	actualSendTime: time.Now(),
				//}
				if ok := ackBuf.Put(msg.MsgId, ackPeek{
					ack:            ack,
					actualSendTime: time.Now(),
				}); ok {
					go ackPipeline(pending, ackBuf, msg.MsgId, msg, ackBuf0[msg.MsgId].actualSendTime)
					//continue
				} else {
					err = errs.AckBuffCapLimit
					goto Fatal
				}
			}
			var n int
			n, err = tcpConn.Write(protocol.Pack(msg.Marshaled))
			logger.Warn(msg.Marshaled)
			logger.Debugf("write %d", n)
			if err != nil {
				goto Fatal
			}

		case id := <-ackBuf.Del:
			if peek, exist := ackBuf0[id]; exist {
				peek.ack()
				delete(ackBuf0, id)
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
	if _, duplicate := err.(*errs.DuplicateConnIdErr); duplicate {
		clErr := conn.tcpConn.Close()
		cancelFunc()
		if clErr != nil {
			logger.Error(clErr)
		}
		return
	}
	conn.close()
	cancelFunc()
}
func ackPipeline[K comparable, V any](ctx context.Context, pds utils.ChanMap[K, V], id K, p pkg.SendMarshal, s time.Time) {
	t := time.NewTimer(time.Second * 30)
	defer t.Stop()
	select {
	case <-t.C:
		pds.Del <- id
		//TODO 消息持久化
		logger.Warnf("ack time out,msg id:%s", p.MsgId)
	case <-ctx.Done():
		//pds.Del <- id
	}
}

func newClient(tcpConn net.Conn, id uint64) {
	wch := make(chan pkg.SendMarshal, 1024)
	errCh := make(chan error, 5)
	seq := atomic.Value{}
	seq.Store(uint64(0))
	conn := &Conn{
		Id:         id,
		tcpConn:    tcpConn,
		ReaderBuff: netrw.ReaderBuff{Buffer: make([]byte, pkg.MaxLen)},
		wch:        wch,
		errMsg:     errCh,
		Addr:       &ConnAddr{tcpConn.RemoteAddr(), id},
		sendSeq:    seq,
	}
	go connHandle(wch, errCh, id, tcpConn, conn)
	ConnAddCh <- conn
}
