package cli

import (
	"context"
	"github.com/NullpointerW/golwpush/errs"
	"github.com/NullpointerW/golwpush/logger"
	"github.com/NullpointerW/golwpush/netrw"
	"github.com/NullpointerW/golwpush/pkg"
	"github.com/NullpointerW/golwpush/protocol"
	"net"
	"sync"
	"time"
)

type PushCli interface {
	Read() (msg string, err error)
	Write(p string) (length int, err error)
	Close()
	PongRecv()
}

type client struct {
	ctx   context.Context
	cFunc context.CancelFunc
	/*buffer  []byte
	wBufPos int*/
	netrw.ReaderBuff
	id      uint64
	tcpConn net.Conn
	pongCh  chan struct{}
	//WMu     sync.Mutex
}

var pongRecv = &sync.Once{}

func (cli *client) Read() (msg string, err error) {
	return netrw.ReadTcp(cli.tcpConn, &cli.ReaderBuff)
}

func (cli *client) Write(p string) (length int, err error) {
	//cli.WMu.Lock()
	//defer cli.WMu.Unlock()
	var b []byte
	if err != nil {
		goto fatal
	}
	b = protocol.Pack(p)
	length, err = cli.tcpConn.Write(b)
	if err != nil {
		goto fatal
	}
	return
fatal:
	cli.Close()
	return
}

func (cli *client) Close() {
	cli.cFunc()
	cli.tcpConn.Close()
}

func (cli *client) PongRecv() {
	cli.pongCh <- struct{}{}
}

func HeartbeatCheck(pushCli PushCli) {
	var (
		cli *client
	)
	if conv, ok := pushCli.(*client); !ok {
		pushCli.Close()
		return
	} else {
		cli = conv
	}
	t := time.NewTimer(time.Second * 60)
	defer t.Stop()
	for {
		select {
		case <-cli.ctx.Done():
			return
		case <-t.C:
			logger.Error(errs.HeartbeatTimeout.Error())
			pushCli.Close()
			return
		case <-cli.pongCh:
			pongRecv.Do(
				func() {
					logger.PrintfNonUid(logger.PingOutput|logger.Host, cli.tcpConn.RemoteAddr().String(), "recved pong")
				})

			t.Reset(time.Second * 60)
		}
	}
}

func SendHeartbeat(pushCli PushCli) {
	var (
		cli *client
	)
	if conv, ok := pushCli.(*client); !ok {
		pushCli.Close()
		return
	} else {
		cli = conv
	}
	t := time.NewTimer(time.Second * 30)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			_, err := cli.Write(pkg.PingMarshaled.Marshaled)
			if err != nil {
				logger.Error(err)
				pushCli.Close()
				return
			}
			t.Reset(time.Second * 30)
		case <-cli.ctx.Done():
			return
		}
	}
}

func NewClient(conn net.Conn, id uint64) (cli PushCli, cancelFunc context.CancelFunc) {
	var ctx context.Context
	ctx, cancelFunc = context.WithCancel(context.Background())
	cli = &client{
		ctx:        ctx,
		ReaderBuff: netrw.ReaderBuff{Buffer: make([]byte, pkg.MaxLen)},
		id:         id,
		tcpConn:    conn,
		cFunc:      cancelFunc,
		pongCh:     make(chan struct{}, 1000),
	}
	return cli, cancelFunc
}
