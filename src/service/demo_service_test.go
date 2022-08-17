package service

import (
	// "encoding/json"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/whyun-com/go-micro-tpl/reqcontext"
	"testing"
)

func TestDemo(t *testing.T) {
	Convey("DemoGet函数测试", t, func(c C) {
		kafkaContext := reqcontext.KafkaRequestContext{}

		var request reqcontext.MicroRequest = reqcontext.MicroRequest{
			Header: reqcontext.MessageHeader{
				MsgType: "",
			},
		}

		var packet reqcontext.RequestPacket = reqcontext.RequestPacket{MicroRequest: request}
		reqcontext.InitKafkaRequestContext(&kafkaContext, &packet)
		res := make(chan reqcontext.MicroResponse, 1)
		DemoGet(&kafkaContext, res)
		resContent := <-res
		So(resContent.Result.Code, ShouldEqual, 0)
		So(resContent.ResData, ShouldEqual, "{\"aaa\":111}")
		defer close(res)
	})
}
