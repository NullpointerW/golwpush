package cli

import (
	"context"
	"net"
	"time"
)

type PushCli interface {
	Read() (length int, err error)
	Write(b []byte) (length int, err error)
	Close()
}

type client struct {
	ctx     context.Context
	buffer  []byte
	id      int64
	tcpConn net.Conn
}

func (cli *client) Read() (length int, err error) {
	length, err = cli.tcpConn.Read(cli.buffer)
	return
}

func (cli *client) Write(b []byte) (length int, err error) {
	length, err = cli.tcpConn.Write(b)
	return
}

func (cli *client) Close() {
	cli.tcpConn.Close()
}

func Heartbeat(pushCli PushCli) {
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
	}
	return cli, cancelFunc
}
