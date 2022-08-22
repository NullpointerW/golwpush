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

func respOk(w http.ResponseWriter, any interface{}) {
	resp(OK, w, any)
}
func respSrvErr(w http.ResponseWriter, any interface{}) {
	resp(ServerError, w, any)
}
func respBadReq(w http.ResponseWriter, any interface{}) {
	resp(BadRequest, w, any)
}
func respMethodNA(w http.ResponseWriter, any interface{}) {
	resp(MethodNotAllowed, w, any)
}
func respUnauth(w http.ResponseWriter, any interface{}) {
	resp(Unauthorized, w, any)
}
func respForbid(w http.ResponseWriter, any interface{}) {
	resp(Forbidden, w, any)
}

func resp(code statusCode, w http.ResponseWriter, any interface{}) {
	w.WriteHeader(int(code))
	switch v := any.(type) {
	case uint:
		fmt.Fprintf(w, "%d", v)
	case uint64:
		fmt.Fprintf(w, "%d", v)
	case int64:
		fmt.Fprintf(w, "%d", v)
	case int:
		fmt.Fprintf(w, "%d", v)
	case bool:
		fmt.Fprintf(w, "%t", v)
	case error:
		fmt.Fprintf(w, "%s", v)

	case string:
		fmt.Fprintf(w, "%s", v)
	}
}
