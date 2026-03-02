package response

import (
	"net/http"

	"github.com/hyahm/xmux"
)

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func Success(r *http.Request, msg string, data any) {
	xmux.GetInstance(r).Response.(*Response).Msg = msg
	xmux.GetInstance(r).Response.(*Response).Data = data
}

func Failed(r *http.Request, code int, msg string) {
	xmux.GetInstance(r).Response.(*Response).Msg = msg
	xmux.GetInstance(r).Response.(*Response).Code = code
}
