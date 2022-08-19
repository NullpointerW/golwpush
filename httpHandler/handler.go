package httpHandler

import (
	"GoPush"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

var PushHandler = Handler{
	Adapter: GoPush.Default,
}

type Handler struct {
	GoPush.Adapter
}

func (httpPush Handler) Push(w http.ResponseWriter, req *http.Request) {
	_msg, _ := ioutil.ReadAll(req.Body)
	idStr := req.URL.Query().Get("id")
	idInt, _ := strconv.ParseUint(idStr, 10, 64)
	err := httpPush.Adapter.Push(idInt, string(_msg))
	if err != nil {
		fmt.Fprintf(w, "%s\n", err)
	}
	fmt.Fprintf(w, "ok")
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
		w.Header().Set("Status-Code", "405")
		fmt.Fprintf(w, "method not allowed")
		return
	}
	jsBody, _ := ioutil.ReadAll(req.Body)
	cts := &GoPush.Content{}
	err := json.Unmarshal(jsBody, cts)
	if err != nil {
		w.Header().Set("Status-Code", "500")
		fmt.Fprintf(w, "json unmarshal error")
		return
	}
	//err := httpPush.Adapter.Push(idInt, string(_msg))
	if err != nil {
		fmt.Fprintf(w, "%s\n", err)
	}
	fmt.Fprintf(w, "ok")
}
