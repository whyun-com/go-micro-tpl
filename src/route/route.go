package route

import (
	"github.com/whyun-com/go-micro-tpl/reqcontext"
	"github.com/whyun-com/go-micro-tpl/service"
)

var RouteMap = map[string]reqcontext.ServiceConfig{
	"demo_get": {Service: service.DemoGet},
}
