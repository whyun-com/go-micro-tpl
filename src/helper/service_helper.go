package helper

import (
	"github.com/whyun-com/go-micro-tpl/reqcontext"
)

func CallService(req reqcontext.ContextInterface, conf *reqcontext.ServiceConfig) *reqcontext.MicroResponse {
	res := make(chan reqcontext.MicroResponse, 1)

	conf.Service(req, res)

	resContent := <-res

	if resContent.Result.Code == 0 {
		reqcontext.DoSuccess(req, &resContent, conf)
	} else {
		reqcontext.DoError(req, &resContent, conf)
	}
	close(res)
	return &resContent
}
