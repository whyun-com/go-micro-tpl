package reqcontext

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)
type HTTPRequestContext struct {
	BaseContext
	ctx *fasthttp.RequestCtx
}

func NewHTTPRequestContext(packet *RequestPacket, ctx *fasthttp.RequestCtx) (req *HTTPRequestContext){
	req = &HTTPRequestContext{}
	initBaseContext(&req.BaseContext, packet)
	req.RequestType = REQUEST_TYPE_HTTP
	req.ctx = ctx;
	return req
}
func (req *HTTPRequestContext) setJsonResHeader() {
	req.ctx.Response.Header.Add("content-type","application/json;charset=utf-8")
}

func (req *HTTPRequestContext) doSuccess(data string, conf *ServiceConfig) {
	log.Debug().Msg("DoSuccess in HTTPRequestContext" + data)
	req.setJsonResHeader()
	req.ctx.Response.SetBody([]byte("{\"code\":0,\"data\":" + data + "}"))
}

func (req *HTTPRequestContext) doError(result Result, conf *ServiceConfig) {
	log.Debug().Int("code", int(result.Code)).Msg("DoError in HTTPRequestContext")
	req.setJsonResHeader()
	bytes,err := json.Marshal(result)
	if err != nil {
		log.Error().Msg("序列化错误信息失败")
		return
	}
	req.ctx.Response.SetBody(bytes)
}