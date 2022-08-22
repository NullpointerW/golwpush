package httphandler

import (
	"fmt"
	"net/http"
)

type statusCode uint16

const (
	OK               statusCode = 200
	BadRequest       statusCode = 400
	Unauthorized     statusCode = 401
	Forbidden        statusCode = 403
	MethodNotAllowed statusCode = 405
	ServerError      statusCode = 500
)

func RespOk(w http.ResponseWriter, any interface{}) {
	resp(OK, w, any)
}
func RespSrvErr(w http.ResponseWriter, any interface{}) {
	resp(ServerError, w, any)
}
func RespBadReq(w http.ResponseWriter, any interface{}) {
	resp(BadRequest, w, any)
}
func RespMethodNA(w http.ResponseWriter, any interface{}) {
	resp(MethodNotAllowed, w, any)
}
func RespUauth(w http.ResponseWriter, any interface{}) {
	resp(Unauthorized, w, any)
}
func RespForbid(w http.ResponseWriter, any interface{}) {
	resp(Forbidden, w, any)
}

func resp(code statusCode, w http.ResponseWriter, any interface{}) {
	w.WriteHeader(int(code))
	switch v := any.(type) {
	case uint:
	case uint64:
	case int64:
	case int:
		fmt.Fprintf(w, "%d", v)
	case bool:
		fmt.Fprintf(w, "%t", v)
	case error:
	case string:
		fmt.Fprintf(w, "%s", v)
	}
}
