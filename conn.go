package gopush

import (
	"context"
	"gopush/errs"
	"gopush/logger"
	"gopush/pkg"
	"gopush/protocol"
	"gopush/utils"
	"net"
	"time"
)

type Conn struct {
	Id         uint64
	tcpConn    net.Conn
	readBuf    []byte
	readBufPtr int
	wch        chan<- pkg.SendMarshal
	Addr       string
	errMsg     chan<- error
}

type acklog struct {
	ack context.CancelFunc
	actualSendTime time.Time
}

func (conn *Conn) write(msg *pkg.Package) {
	msg.Id = utils.GenerateId(conn.Id)    //生成消息id
	marshaled, err := msg.MarshalToSend() //避免在写goroutine中编码
	if err != nil {
		conn.errMsg <- err
		return
	}
	conn.wch <- marshaled
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
	logger.Debug("close conn" + conn.Addr)
	err := conn.tcpConn.Close()
	if err != nil {
		logger.Error(err)
	}
	ConnRmCh <- conn
}

func connHandle(wch chan pkg.SendMarshal, errCh chan error, id uint64, tcpConn net.Conn, conn *Conn) {
	pendAck, pends := utils.NewChMap[string, context.CancelFunc](10000) //max 1000
	pingCh := make(chan struct{}, 100)
	ctx, cancel := context.WithCancel(context.Background())
	readHandler := make(chan string, 1024)
	//避免在读goroutine解码，通过一个goroutine处理所有读到的包
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case unDecode := <-readHandler:
				p, err := pkg.New(unDecode)
				if err != nil {
					errCh <- err
				}
				switch p.Mode {
				case pkg.ACK:
					pendAck.Del <- p.Id
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
				readHandler <- msg
			}
		}
	}(ctx)

	go func(ctx context.Context) {
		t := time.NewTimer(time.Minute * 1)
		defer t.Stop()
		cliAddr := "[cli][" + tcpConn.RemoteAddr().String() + "] "
		logger.Debug(cliAddr + "start heartbeat check") //debug
		for {
			select {
			case <-ctx.Done():
				logger.Debug(cliAddr + "heartbeat check end") //debug
				return
			case <-t.C:
				logger.Warn(cliAddr + "Heartbeat timeout 60s...")
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

	for {
		select {
		case msg := <-wch:
			//var strMsg string
			//strMsg, err = msg.ConvStr()
			//if err != nil {
			//	goto Fatal
			//}
			if msg.Mode == pkg.MSG {
				//TODO ACK
				pending, ack := context.WithCancel(context.Background())
				pends[]

				go ackPipeline(pending, pends)
			}
			_, err = tcpConn.Write(protocol.Pack(msg.Marshaled))
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
func ackPipeline[K comparable, V any](ctx context.Context, a utils.ChanMap[K, V]) {

}

func newClient(tcpConn net.Conn, id uint64) {
	wch := make(chan pkg.SendMarshal, 1024)
	errCh := make(chan error, 5)
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
