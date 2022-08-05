package httpHandler

import (
	"GoPush"
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
