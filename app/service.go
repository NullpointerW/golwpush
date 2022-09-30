package main

import (
	"github.com/NullpointerW/golwpush"
	"github.com/NullpointerW/golwpush/httphandler"
	"github.com/NullpointerW/golwpush/logger"
	"log"
	"net"
	"net/http"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	//logger.ModifyLv(logger.Prod)
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
	go golwpush.Handle()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}

		go golwpush.InitConn(conn)

	}
}
