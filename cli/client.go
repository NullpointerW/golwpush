package cli

import (
	"GoPush/errs"
	"GoPush/logger"
	"GoPush/protocol"
	"context"
	"net"
	"time"
)

type PushCli interface {
	Read() (msg string, err error)
	Write(string) (length int, err error)
	Close()
	PongRecv()
}

type client struct {
	ctx        context.Context
	cFunc      context.CancelFunc
	buffer     []byte
	readBufPtr int
	id         int64
	tcpConn    net.Conn
	pongCh     chan struct{}
}

func (cli *client) Read() (msg string, err error) {
	length, TCPErr := cli.tcpConn.Read(cli.buffer[cli.readBufPtr:])
	if TCPErr != nil {
		return msg, TCPErr
	}
	msg, cli.readBufPtr, err = protocol.Unpack(cli.buffer[:length+cli.readBufPtr])
	return
}

func (cli *client) Write(msg string) (length int, err error) {
	b := protocol.Pack(msg)
	length, err = cli.tcpConn.Write(b)
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
			_, err := cli.Write("ping")
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

func NewClient(conn net.Conn, id int64) (cli PushCli, cancelFunc context.CancelFunc) {
	var ctx context.Context
	ctx, cancelFunc = context.WithCancel(context.Background())
	cli = &client{
		ctx:        ctx,
		buffer:     make([]byte, 1024),
		readBufPtr: 0,
		id:         id,
		tcpConn:    conn,
		cFunc:      cancelFunc,
		pongCh:     make(chan struct{}, 1000),
	}
	return cli, cancelFunc
}
