package cli

import (
	"context"
	"gopush/errs"
	"gopush/logger"
	"gopush/pkg"
	"gopush/protocol"
	"net"
	"time"
)

type PushCli interface {
	Read() (msg string, err error)
	Write(p *pkg.Package) (length int, err error)
	Close()
	PongRecv()
}

type client struct {
	ctx        context.Context
	cFunc      context.CancelFunc
	buffer     []byte
	readBufPtr int
	id         uint64
	tcpConn    net.Conn
	pongCh     chan struct{}
}

func (cli *client) Read() (msg string, err error) {
	length, TCPErr := cli.tcpConn.Read(cli.buffer[cli.readBufPtr:])
	if TCPErr != nil {
		return msg, TCPErr
	}
	var retry bool
unpack:
	msg, cli.readBufPtr, retry, err = protocol.Unpack(cli.buffer[:length+cli.readBufPtr], cli.readBufPtr)
	if err != nil {
		return msg, err
	}
	if retry {
		goto unpack
	}
	return
}

func (cli *client) Write(p *pkg.Package) (length int, err error) {
	var b []byte
	strMsg, err := p.ConvStr()
	if err != nil {
		goto fatal
	}
	b = protocol.Pack(strMsg)
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
			logger.Info("recv pong")
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
			_, err := cli.Write(pkg.Ping)
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
		buffer:     make([]byte, pkg.MaxLen),
		readBufPtr: 0,
		id:         id,
		tcpConn:    conn,
		cFunc:      cancelFunc,
		pongCh:     make(chan struct{}, 1000),
	}
	return cli, cancelFunc
}
