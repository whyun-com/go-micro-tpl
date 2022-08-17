package filter

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
	"github.com/whyun-com/go-micro-tpl/helper"
	"github.com/whyun-com/go-micro-tpl/reqcontext"
	"github.com/whyun-com/go-micro-tpl/route"
)

const HEADER_MSG_TYPE = "message-msg-type"

func DoHttpFilter(ctx *fasthttp.RequestCtx) {
	packet := &reqcontext.RequestPacket{}
	req := reqcontext.NewHTTPRequestContext(packet, ctx)

	msgType := string(ctx.Request.Header.Peek(HEADER_MSG_TYPE))
	packet.Header.MsgType = msgType
	configRoute, ok := route.RouteMap[msgType]
	responseContent := reqcontext.GetDefaultMicroResponse()
	if !ok {
		log.Error().Msg(msgType + "类型消息不被支持")
		responseContent.Result.Code = 1001             //TODO给出具体code定义
		reqcontext.DoError(req, &responseContent, nil) //TODO 给出具体数据结构
		return
	}

	if string(ctx.Method()) != "POST" {
		log.Error().Msg("只支持post方法")
		responseContent.Result.Code = 1002           //TODO给出具体code定义
		reqcontext.DoError(req, &responseContent, nil) //TODO 给出具体数据结构
		return
	}
	var data reqcontext.MicroRequest
	err := json.Unmarshal(ctx.PostBody(), &data)
	if err != nil {
		log.Error().Err(err).Msg("反序列化http请求数据失败")
		responseContent.Result.Code = 1003             //TODO给出具体code定义
		reqcontext.DoError(req, &responseContent, nil) //TODO 给出具体数据结构
		return
	}


	packet.MicroRequest = data
	helper.CallService(req, &configRoute)
}
