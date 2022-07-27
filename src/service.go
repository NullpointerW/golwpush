package main

import (
	"GoPush/src/push"
	"GoPush/src/push/httpHandler"
	"log"
	"net"
	"net/http"
)

type a interface {
	a()
}

type st struct {
	a
}

func main() {
	log.Println("start pushServer")

	go func() {
		log.Println("正在启动http服务...")
		mux := http.NewServeMux()
		h := httpHandler.HttpPushHandler
		mux.Handle("/push", http.HandlerFunc(h.ReqPush))
		mux.Handle("/broadcast", http.HandlerFunc(h.ReqBroadcast))
		log.Fatal(http.ListenAndServe("localhost:8000", mux))
	}()
	log.Println("正在启动tcp服务...")
	listener, err := net.Listen("tcp", "localhost:9000")
	if err != nil {
		log.Fatal(err)
	}
	go push.Broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}

		go push.InitConn(conn)

	}
}
