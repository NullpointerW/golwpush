package main

import (
	"GoPush/src/cli"
	"fmt"
	"log"
	"net"
	"strconv"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:9003")
	if err != nil {
		log.Fatal(err)
	}

	var (
		id = 114514
	)
	pCli, cancelFunc := cli.NewClient(conn, int64(id))
	defer pCli.Close()
	msg := []byte(strconv.Itoa(id))
	_, wErr := pCli.Write(msg)
	if wErr != nil {
		fmt.Errorf("")
		return
	}

	for {

	}

}
