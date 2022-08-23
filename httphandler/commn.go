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

func respOk(w http.ResponseWriter, t any) {
	resp(OK, w, t)
}
func respSrvErr(w http.ResponseWriter, t any) {
	resp(ServerError, w, t)
}
func respBadReq(w http.ResponseWriter, t any) {
	resp(BadRequest, w, t)
}
func respMethodNA(w http.ResponseWriter, t any) {
	resp(MethodNotAllowed, w, t)
}
func respUnauth(w http.ResponseWriter, t any) {
	resp(Unauthorized, w, t)
}
func respForbid(w http.ResponseWriter, t any) {
	resp(Forbidden, w, t)
}

func resp(code statusCode, w http.ResponseWriter, t any) {
	w.WriteHeader(int(code))
	switch v := t.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		fmt.Fprintf(w, "%d", v)
	case bool:
		fmt.Fprintf(w, "%t", v)
	case error:
		fmt.Fprintf(w, "error:%s", v)
	case string:
		fmt.Fprintf(w, "%s", v)
	}
}
