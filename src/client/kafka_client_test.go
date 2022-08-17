package client

import (
	"encoding/json"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/whyun-com/go-micro-tpl/config"
	"github.com/whyun-com/go-micro-tpl/reqcontext"
)

func TestKafkaSend(t *testing.T) {
	Convey("kafka测试", t, func() {
		var request reqcontext.MicroRequest
		err := json.Unmarshal([]byte(config.MOCK_MICRO_REQUEST_JSON), &request)
		if err != nil {
			t.Fatal(err)
		}
		reqPacket := reqcontext.RequestPacket{
			MicroRequest: request,
		}

		option := KafkaOption{
			KafkaConf: &kafka.ConfigMap{
				"bootstrap.servers": config.Config.Kafka.Communication.Brokers,
			},
			Topic: "test",
		}

		kafkaClient, err := NewKafkaClient(&option)

		if err != nil {
			t.Fatal(err)
			return
		}

		resChan := make(chan reqcontext.MicroResponse, 1)
		kafkaClient.Send(&reqPacket, resChan)

		response := <-resChan

		if response.Result.Code != 0 {
			t.Fatalf("kakfa send出错了：%d", response.Result.Code)
			return
		}
		t.Log("kafka send finished")
	})
}
