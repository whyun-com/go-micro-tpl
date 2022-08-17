package app

import (
	"github.com/valyala/fasthttp"
	"github.com/whyun-com/go-micro-tpl/filter"
)

func FastHTTPHandler() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/healthz":
			ctx.Response.SetBody([]byte("OK"))
		case "/i/micro":
			filter.DoHttpFilter(ctx)
		default:
			ctx.Error("not found", fasthttp.StatusNotFound)
		}
	}
}
