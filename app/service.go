package main

import (
	"github.com/NullpointerW/gopush"
	"github.com/NullpointerW/gopush/httphandler"
	"github.com/NullpointerW/gopush/logger"
	"log"
	"net"
	"net/http"
)

func init() {
	logger.ModifyLv(logger.Dev)
}

func main() {
	logger.Info("start pushServer")
	go func() {
		logger.Infof("staring http server...")
		mux := http.NewServeMux()
		h := httphandler.PushHandler
		mux.Handle("/push", http.HandlerFunc(h.Push))
		mux.Handle("/broadcast", http.HandlerFunc(h.Broadcast))
		mux.Handle("/multiPush", http.HandlerFunc(h.MultiPush))
		mux.Handle("/count", http.HandlerFunc(h.Count))
		mux.Handle("/info", http.HandlerFunc(h.Info))
		log.Fatal(http.ListenAndServe("localhost:8000", mux))
	}()
	logger.Infof("staring tcp server...")
	listener, err := net.Listen("tcp", "localhost:9000")
	if err != nil {
		logger.Fatal(err)
	}
	go gopush.Handle()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}

		go gopush.InitConn(conn)

	}
}
