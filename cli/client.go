package cli

import (
	"GoPush/errs"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

type PushCli interface {
	Read(buff []byte) (length int, err error)
	Write(b []byte) (length int, err error)
	Close()
	PongRecv()
}

type client struct {
	ctx     context.Context
	cFunc   context.CancelFunc
	buffer  []byte
	id      int64
	tcpConn net.Conn
	pongCh  chan struct{}
}

func (cli *client) Read(buff []byte) (length int, err error) {
	length, err = cli.tcpConn.Read(buff)
	return
}

func (cli *client) Write(b []byte) (length int, err error) {
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
			fmt.Fprintf(os.Stderr, errs.HeartbeatTimeout.Error())
			pushCli.Close()
			return
		case <-cli.pongCh:
			log.Println("recv pong")
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
			cli.Write([]byte("ping"))
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
		ctx:     ctx,
		buffer:  make([]byte, 128),
		id:      id,
		tcpConn: conn,
		cFunc:   cancelFunc,
		pongCh:  make(chan struct{}, 1000),
	}
	return cli, cancelFunc
}
