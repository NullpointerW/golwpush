package main

import (
	"GoPush/cli"
	"GoPush/logger"
	"GoPush/pkg"
	"GoPush/protocol"
	"encoding/binary"
	"math/rand"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		logger.Fatal(err)
	}
	logger.Infof("connect to server %s\n", conn.RemoteAddr().String())
	rand.Seed(time.Now().UnixNano())
	var (
		id = rand.Intn(10000) //随机生成id
	)
	pCli, _ := cli.NewClient(conn, uint64(id))
	defer pCli.Close()
	//msg := strconv.Itoa(id)
	//_, wErr := pCli.Write(msg)
	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, uint64(id))
	trans := protocol.PackByteStream(8, data)
	_, wErr := conn.Write(trans)
	if wErr != nil {
		logger.Fatalf("write error: %v", wErr)
	}
	logger.Debugf("sendId:%d succeed \n", id)
	go cli.SendHeartbeat(pCli)
	go cli.HeartbeatCheck(pCli)

	for {
		msg, err := pCli.Read()
		if err != nil {
			fatal(err, pCli)
		}
		tPkg, pkgErr := pkg.New(msg)
		if pkgErr != nil {
			fatal(pkgErr, pCli)
		}
		switch tPkg.Mode {
		case pkg.PONG:
			pCli.PongRecv()
		case pkg.MSG:
			logger.Infof(tPkg.Data)
		}
	}

}
func fatal(err error, pCli cli.PushCli) {
	logger.Fatal(err)
	pCli.Close()
	return
}
