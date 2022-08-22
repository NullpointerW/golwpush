package httphandler

import (
	"GoPush"
	"encoding/json"
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
		respSrvErr(w, err)
		return
	}
	respOk(w, "ok")
}

func (httpPush Handler) Broadcast(w http.ResponseWriter, req *http.Request) {
	_msg, _ := ioutil.ReadAll(req.Body)
	err := httpPush.Adapter.Broadcast(string(_msg))
	if err != nil {
		respSrvErr(w, err)
		return
	}
	respOk(w, "ok")
}

func (httpPush Handler) MultiPush(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		respMethodNA(w, "non-post method not allowed")
		return
	}
	jsBody, _ := ioutil.ReadAll(req.Body)
	cts := &gopush.Contents{}
	err := json.Unmarshal(jsBody, cts)
	if err != nil {
		respBadReq(w, "json unmarshal error")
		return
	}
	cts.Res = make(chan uint64, 1)
	err, ok := httpPush.Adapter.MultiPush(cts)
	if err != nil {
		respSrvErr(w, err)
		return
	}
	respOk(w, ok)
}
