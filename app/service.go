package main

import (
	"GoPush"
	"GoPush/httpHandler"
	"GoPush/logger"
	"log"
	"net"
	"net/http"
)

func init() {
	logger.ModifyLv(logger.Prod)
}

func main() {
	logger.Info("start pushServer")
	go func() {
		logger.Infof("staring http server...")
		mux := http.NewServeMux()
		h := httpHandler.PushHandler
		mux.Handle("/push", http.HandlerFunc(h.Push))
		mux.Handle("/broadcast", http.HandlerFunc(h.Broadcast))
		log.Fatal(http.ListenAndServe("localhost:8000", mux))
	}()
	logger.Infof("staring tcp server...")
	listener, err := net.Listen("tcp", "localhost:9000")
	if err != nil {
		logger.Fatal(err)
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
