package main

import (
	"GoPush/src/cli"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:9003")
	if err != nil {
		log.Fatal(err)
	}

	var (
		id   = 114514
		buff = make([]byte, 1024)
	)
	pCli, _ := cli.NewClient(conn, int64(id))
	defer pCli.Close()
	msg := []byte(strconv.Itoa(id))
	_, wErr := pCli.Write(msg)
	if wErr != nil {
		fmt.Errorf("write error: %v", wErr)
		return
	}
	go cli.SendHeartbeat(pCli)
	go cli.HeartbeatCheck(pCli)

	for {
		l, err := pCli.Read(buff)
		if err != nil {
			pCli.Close()
			return
		}
		msg := string(buff[:l])
		if strings.EqualFold(msg, "pong") {
			pCli.PongRecv()
		} else {
			fmt.Println(msg)
		}
	}
}
