package main

import (
	"GoPush/cli"
	"GoPush/logger"
	"net"
	"strconv"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		logger.Fatal(err)
	}
	logger.Infof("connect to server %s\n", conn.RemoteAddr().String())

	var (
		id   = 114514
		buff = make([]byte, 1024)
	)
	pCli, _ := cli.NewClient(conn, int64(id))
	defer pCli.Close()
	msg := []byte(strconv.Itoa(id))
	_, wErr := pCli.Write(msg)
	if wErr != nil {
		logger.Fatalf("write error: %v", wErr)
	}
	logger.Infof("send id:%d  succeed \n", id)
	go cli.SendHeartbeat(pCli)
	go cli.HeartbeatCheck(pCli)

	for {
		l, err := pCli.Read(buff)
		if err != nil {
			logger.Fatal(err)
			pCli.Close()
			return
		}
		msg := string(buff[:l])
		if strings.EqualFold(msg, "pong") {
			pCli.PongRecv()
		} else {
			logger.Infof(msg)
		}
	}
}
