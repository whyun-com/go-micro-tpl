package client

import (
	"github.com/whyun-com/go-micro-tpl/reqcontext"
)

type ClientInterface interface {
	Send(req *reqcontext.RequestPacket, res chan<- reqcontext.MicroResponse)
}
