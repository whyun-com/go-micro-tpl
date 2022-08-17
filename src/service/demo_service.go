package service

import (
	"github.com/whyun-com/go-micro-tpl/reqcontext"
)

// func(req *reqcontext.ContextInterface)
func DemoGet(req reqcontext.ContextInterface, res chan<- reqcontext.MicroResponse) {
	content := reqcontext.MicroResponse{
		Result:  reqcontext.Result{Code: 0},
		ResData: "{\"aaa\":111}",
	}
	res <- content
}
