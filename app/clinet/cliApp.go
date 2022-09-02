package main

import (
	"encoding/json"
	"github.com/NullpointerW/golwpush/cli"
	"github.com/NullpointerW/golwpush/logger"
	"github.com/NullpointerW/golwpush/pkg"
	"github.com/NullpointerW/golwpush/utils"
	"math/rand"
	"net"
	"time"
)

func init() {
	logger.ModifyLv(logger.Dev)
}

func main() {
	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		logger.Fatal(err)
	}
	logger.PlnNUid(logger.Login|logger.Cli|logger.Host, conn.RemoteAddr().String(), "connected to server")
	rand.Seed(time.Now().UnixNano())
	var (
		uid = rand.Intn(10000) //随机生成id
	)
	pCli, _ := cli.NewClient(conn, uint64(uid))
	defer pCli.Close()
	pCli.Auth()
	go cli.SendHeartbeat(pCli)
	go cli.HeartbeatCheck(pCli)

	for {
		msg, err := pCli.Read()
		if err != nil {
			fatal(err, pCli)
		}
		tPkg, pkgErr := pkg.New(msg)
		if pkgErr != nil {
			logger.Error("json err")
			logger.Error(msg)
			fatal(pkgErr, pCli)
		}
		switch tPkg.Mode {
		case pkg.PONG:
			pCli.PongRecv()
		case pkg.MSG:
			logger.PlnNUid(logger.MsgOutput|logger.Host, conn.RemoteAddr().String(), tPkg.Data)
			recvTime := time.Now().Format(utils.TimeParseLayout)
			raw, _ := json.Marshal(&pkg.Package{Mode: pkg.ACK, Id: tPkg.Id, Data: recvTime}) //ack 确认
			pCli.Write(utils.Bcs(raw))
		}
	}

}
func fatal(err error, pCli cli.PushCli) {
	pCli.Close()
	logger.Fatal(err)
}
