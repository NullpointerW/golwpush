package main

import (
	"GoPush"
	"GoPush/httpHandler"
	"log"
	"net"
	"net/http"
)

func main() {
	log.Println("start pushServer")

	go func() {
		log.Println("正在启动http服务...")
		mux := http.NewServeMux()
		h := httpHandler.PushHandler
		mux.Handle("/push", http.HandlerFunc(h.Push))
		mux.Handle("/broadcast", http.HandlerFunc(h.Broadcast))
		log.Fatal(http.ListenAndServe("localhost:8000", mux))
	}()
	log.Println("正在启动tcp服务...")
	listener, err := net.Listen("tcp", "localhost:9000")
	if err != nil {
		log.Fatal(err)
	}
	go GoPush.Handle()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}

		go GoPush.InitConn(conn)

	}
}
