package httphandler

import (
	"GoPush"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

var PushHandler = Handler{
	Adapter: gopush.Default,
}

type Handler struct {
	gopush.Adapter
}

func (httpPush Handler) Push(w http.ResponseWriter, req *http.Request) {
	_msg, _ := ioutil.ReadAll(req.Body)
	idStr := req.URL.Query().Get("id")
	idInt, _ := strconv.ParseUint(idStr, 10, 64)
	err := httpPush.Adapter.Push(idInt, string(_msg))
	if err != nil {
		RespSrvErr(w, err)
	}
	RespOk(w, "ok")
}

func (httpPush Handler) Broadcast(w http.ResponseWriter, req *http.Request) {
	_msg, _ := ioutil.ReadAll(req.Body)
	err := httpPush.Adapter.Broadcast(string(_msg))
	if err != nil {
		fmt.Fprintf(w, "%s\n", err)
	}
	fmt.Fprintf(w, "ok")
}

func (httpPush Handler) MultiPush(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		w.WriteHeader(405)
		fmt.Fprintf(w, "method not allowed")
		return
	}
	jsBody, _ := ioutil.ReadAll(req.Body)
	cts := &gopush.Contents{}
	err := json.Unmarshal(jsBody, cts)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "json unmarshal error")
		return
	}
	cts.Res = make(chan uint64, 1)
	err, ok := httpPush.Adapter.MultiPush(cts)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "sever internal error %s", err)
		return
	}
	fmt.Fprintf(w, "%d", ok)
}
