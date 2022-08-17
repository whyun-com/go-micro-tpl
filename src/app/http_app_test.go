package app

import (
	"encoding/json"

	"github.com/whyun-com/go-micro-tpl/config"
	"github.com/whyun-com/go-micro-tpl/reqcontext"
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestHttp(t *testing.T) {
	Convey("healthz test", t, func(c C) {
		e := ConveyHTTPTester(c)

		e.GET("/healthz").Expect().
			Status(200).
			Text().Equal("OK")
	})
	SkipConvey("not exist test", t, func(c C) { //这个地方故意让其不成功
		e := ConveyHTTPTester(c)

		e.GET("/not-exist").Expect().
			Status(200).
			Text().Equal("OK")
	})
}

func TestMicro(t *testing.T) {
	Convey("请求demo_get消息", t, func(c C) {
		var request reqcontext.MicroRequest
		err := json.Unmarshal([]byte(config.MOCK_MICRO_REQUEST_JSON), &request)
		if err != nil {
			t.Fatal(err)
			return
		}

		e := ConveyHTTPTester(c)
		e.POST("/i/micro").WithJSON(request).
			WithHeader("message-msg-type", "demo_get").
			Expect().
			Status(200).
			JSON().Object().ContainsKey("code")
	})
}
