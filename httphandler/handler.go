package httphandler

import (
	"encoding/json"
	"fmt"
	"gopush"
	"gopush/utils"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var PushHandler = Handler{
	Adapter: gopush.Default,
}

type Handler struct {
	gopush.Adapter
}

func (h Handler) Push(w http.ResponseWriter, req *http.Request) {
	_msg, _ := ioutil.ReadAll(req.Body)
	idStr := req.URL.Query().Get("id")
	uid, _ := strconv.ParseUint(idStr, 10, 64)
	err := h.Adapter.Push(uid, utils.Bcs(_msg))
	if err != nil {
		respSrvErr(w, err)
		return
	}
	respOk(w, "ok")
}

func (h Handler) Broadcast(w http.ResponseWriter, req *http.Request) {
	_msg, _ := ioutil.ReadAll(req.Body)
	err := h.Adapter.Broadcast(utils.Bcs(_msg))
	if err != nil {
		respSrvErr(w, err)
		return
	}
	respOk(w, "ok")
}

func (h Handler) MultiPush(w http.ResponseWriter, req *http.Request) {
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
	err, ok := h.Adapter.MultiPush(cts)
	if err != nil {
		respSrvErr(w, err)
		return
	}
	respOk(w, ok)
}

func (h Handler) Count(w http.ResponseWriter, req *http.Request) {
	c := h.Adapter.Count()
	respOk(w, c)
}

func (h Handler) Info(w http.ResponseWriter, req *http.Request) {
	idStr := req.URL.Query().Get("id")
	uid, _ := strconv.ParseUint(idStr, 10, 64)
	res := make(chan any, 1)
	i, err := h.Adapter.Info(gopush.BizReq{Res: res, Uid: uid, Typ: gopush.Info})
	if err != nil {
		respSrvErr(w, err)
		return
	}
	if i == nil {
		respOk(w, fmt.Sprintf("uid :%d offline [%s]", uid, time.Now().Format(utils.TimeParseLayout)))
		return
	}
	w.Header().Add("Content-Type", "application/json")
	j, _ := json.Marshal(i)
	respOk(w, utils.Bcs(j))
}
