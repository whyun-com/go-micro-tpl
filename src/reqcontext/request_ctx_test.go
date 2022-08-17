package reqcontext

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/whyun-com/go-micro-tpl/config"
)

const MICRO_REQUEST_JSON = config.MOCK_MICRO_REQUEST_JSON

func TestExtends(t *testing.T) {
	Convey("inner 继承测试", t, func() {
		kafkaContext := KafkaRequestContext{}
		var request MicroRequest
		err := json.Unmarshal([]byte(MICRO_REQUEST_JSON), &request)
		if err != nil {
			t.Fatal(err)
		}
		var packet RequestPacket = RequestPacket{MicroRequest: request}
		InitKafkaRequestContext(&kafkaContext, &packet)
		t.Logf("%v %v", kafkaContext.Begin, kafkaContext.EndMinusBegin)
		var context ContextInterface = &kafkaContext
		DoSuccess(context, &MicroResponse{}, &ServiceConfig{})
		// _doSucess(&kafkaContext)
	})
	Convey("inner json测试", t, func() {
		// kafkaContext := KafkaRequestContext{}
		var request MicroRequest
		err := json.Unmarshal([]byte(MICRO_REQUEST_JSON), &request)
		if err != nil {
			t.Fatal(err)
		}
		// t.Logf("%+v",request)
	})
}
