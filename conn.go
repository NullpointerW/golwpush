package golwpush

import (
	"context"
	"encoding/json"
	"github.com/NullpointerW/golwpush/errs"
	"github.com/NullpointerW/golwpush/logger"
	"github.com/NullpointerW/golwpush/netrw"
	"github.com/NullpointerW/golwpush/persist"
	"github.com/NullpointerW/golwpush/pkg"
	"github.com/NullpointerW/golwpush/protocol"
	"github.com/NullpointerW/golwpush/utils"
	"github.com/go-redis/redis"
	"math"
	"net"
	"strconv"
	"sync/atomic"
	"time"
)

type Conn struct {
	Uid     uint64
	tcpConn net.Conn
	/*readBuf []byte
	wBufPos int*/
	tcpReader netrw.Reader
	wch       chan<- pkg.SendMarshal
	Addr      Addr
	errMsg    chan<- error
	sendSeq   atomic.Value
}

type ackTicker struct {
	pkg                      pkg.SendMarshal
	pending                  context.Context
	deadline, actualSendTime time.Time
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
	msg.Id = utils.GenerateId(conn.Uid) + strconv.FormatUint(conn.incrSeq(), 10) //生成消息id
	marshaled, err := msg.MarshalToSend()                                        //避免在写goroutine中编码
	if err != nil {
		conn.errMsg <- err
		return
	}
	conn.wch <- marshaled
}

func (conn *Conn) read() (msg string, err error) {
	return conn.tcpReader.Read()
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
	reset := make(chan struct{}, 100) //心跳重置
	ctx, cancel := context.WithCancel(context.Background())
	rHandler := make(chan string, 2048)
	ackCh := make(chan ackTicker, 2048)

	//避免在读goroutine解码，通过一个goroutine处理所有读到的包
	go readHandle(ctx, rHandler, errCh, pingCh, reset, ackBuf)
	//读 goroutine
	go readLoop(ctx, conn, errCh, rHandler)
	//心跳检测
	go heartBeatCheck(ctx, conn, errCh, pingCh, reset, wch)
	//丢失消息重发
	go msgRetransmission(conn, ctx)
	//ack消息确认/持久化
	go ackPipelineV2(ackCh, conn.Uid, ctx, ackBuf.Del)

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
				curr := time.Now()
				deadline := curr.Add(30 * time.Second)
				if ok := ackBuf.Put(msg.MsgId, ackPeek{
					ack:            ack,
					actualSendTime: curr,
				}); ok {
					ackCh <- ackTicker{msg, pending, deadline, curr}
				} else {
					err = errs.AckBuffCapLimit
					goto fatal
				}
			}
			//var n int
			_, err = tcpConn.Write(protocol.Pack(msg.Marshaled))
			logger.Infof("send msg:%s len:%d", string(protocol.Pack(msg.Marshaled)), len(protocol.Pack(msg.Marshaled)))
			//logger.Warn(msg.Marshaled)
			//logger.Debugf("write %d", n)
			if err != nil {
				goto fatal
			}

		case id := <-ackBuf.Del:
			if peek, exist := ackBuf0[id]; exist {
				logger.Infof("ack from cli %s,msg id:%s", conn.Addr.String(), id)
				peek.ack()
				delete(ackBuf0, id)
			}

		case err = <-errCh:
			goto fatal
		}
	}
fatal:
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
func readHandle(ctx context.Context, rHandler chan string, errCh chan error, pingCh, reset chan struct{},
	ackBuf utils.ChanMap[string, ackPeek]) {
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
				reset <- struct{}{}
				ackBuf.Del <- p.Id

			case pkg.PING:
				pingCh <- struct{}{}
			}
		}
	}
}

func readLoop(ctx context.Context, conn *Conn, errCh chan error, rHandler chan string) {
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
}

func heartBeatCheck(ctx context.Context, conn *Conn, errCh chan error, pingCh, reset chan struct{}, wch chan pkg.SendMarshal) {
	t := time.NewTimer(time.Minute * 1)
	defer t.Stop()
	logger.PlnWAddr(logger.L_Debug|logger.HeartBeat|logger.Srv, conn.Addr.Uid(), conn.Addr.String(),
		"start heartbeat check") //debug
	for {
		select {
		case <-ctx.Done():
			logger.PlnWAddr(logger.L_Debug|logger.PingOutput, conn.Addr.Uid(), conn.Addr.String(),
				"heartbeat check end") //debug
			return
		case <-t.C:
			logger.PlnWAddr(logger.L_Warn|logger.HeartBeat|logger.Srv, conn.Addr.Uid(), conn.Addr.String(),
				"Heartbeat timeout 60s...") //debug
			errCh <- errs.HeartbeatTimeout
			return
		case <-pingCh:
			wch <- pkg.PongMarshaled
			t.Reset(time.Minute * 1)
		case <-reset:
			t.Reset(time.Minute * 1)
		}
	}
}

func ackPipeline[K comparable, V any](ctx context.Context, pds utils.ChanMap[K, V], id K, p pkg.SendMarshal,
	s time.Time, uid uint64) {
	t := time.NewTimer(time.Second * 30)
	defer t.Stop()
	select {
	case <-t.C:
		pds.Del <- id
		//TODO 消息持久化
		mem, _ := json.Marshal(p)
		persist.Redis.ZAdd(persist.KeyCache.Key(strconv.FormatUint(uid, 10)), redis.Z{Score: float64(s.UnixMilli()),
			Member: mem})
		logger.Warnf("ack time out,msg id:%s", p.MsgId)
	case <-ctx.Done():
		logger.Warnf("ack from id:%s", p.MsgId)
		return
	}
}
func ackPipelineV2(ackReceiver chan ackTicker, uid uint64, ctx context.Context, pdsDel chan string) {
	dbWriteTick := time.NewTimer(time.Second * 30)
	key := persist.RedisKeyPrefix + strconv.FormatUint(uid, 10)
	var zmem []redis.Z
	defer dbWriteTick.Stop()
	for {
		select {
		case <-ctx.Done():
			persist.Redis.ZAdd(key, zmem...)
			return
		case <-dbWriteTick.C:
			if len(zmem) > 0 {
				persist.Redis.ZAdd(key, zmem...)
				zmem = nil
			}
		case ack := <-ackReceiver:
			cmp, d := utils.TimeCmp(ack.deadline, time.Now())
			if cmp < 1 {
				goto check
			}
			time.Sleep(d)
		check:
			select {
			case <-ack.pending.Done():
			default:
				//db
				pdsDel <- ack.pkg.MsgId
				v, _ := json.Marshal(ack.pkg)
				mem := redis.Z{Score: float64(ack.actualSendTime.UnixMilli()),
					Member: v}
				zmem = append(zmem, mem)
				if len(zmem) >= 500 {
					persist.Redis.ZAdd(key, zmem...)
					zmem = nil
				}
			}
		}
	}
}

func msgRetransmission(conn *Conn, ctx context.Context) {
	k := persist.KeyCache.Key(strconv.FormatUint(conn.Uid, 10))
	c := persist.Redis.ZCard(k).Val()
	size := 5000
	partition := int(math.Ceil(float64(c) / float64(size)))
	if partition > 0 {
		logger.Warnf("uid[%d] msgRetransmission msg num:%d", conn.Uid, c)
	}
	for i := 0; i < partition; i++ {
		var del []interface{}
		select {
		case <-ctx.Done():
			return
		default:
			r, err := persist.Redis.ZRange(k, int64(i*size), int64(i*size+size)-1).Result()
			//r, err := persist.Redis.ZPopMin(k, int64(size)).Result()
			if err != nil {
				//todo
				return
			}
			for _, jsonRaw := range r {
				//jsonRaw, ok := z.Member.(string)
				//if !ok {
				//	return
				//}
				var p pkg.SendMarshal
				err = json.Unmarshal(utils.Scb(jsonRaw), &p)
				if err != nil {
					//todo handle
					return
				}
				conn.wch <- p
				del = append(del, jsonRaw)
			}
		}
		persist.Redis.ZRem(k, del...)
	}
}

func newClient(tcpConn net.Conn, id uint64) {
	wch := make(chan pkg.SendMarshal, 2048)
	errCh := make(chan error, 5)
	seq := atomic.Value{}
	seq.Store(uint64(0))
	conn := &Conn{
		Uid:       id,
		tcpConn:   tcpConn,
		tcpReader: &netrw.TcpReader{Buf: make([]byte, protocol.MaxLen), Conn: tcpConn},
		wch:       wch,
		errMsg:    errCh,
		Addr:      &ConnAddr{tcpConn.RemoteAddr(), id},
		sendSeq:   seq,
	}
	go connHandle(wch, errCh, id, tcpConn, conn)
	ConnAddCh <- conn
}
