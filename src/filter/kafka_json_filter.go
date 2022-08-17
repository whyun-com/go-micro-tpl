package filter

import (
	"encoding/json"
	// "fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rs/zerolog/log"
	"github.com/whyun-com/go-micro-tpl/helper"
	"github.com/whyun-com/go-micro-tpl/reqcontext"
	"github.com/whyun-com/go-micro-tpl/route"
)

type KafkaFilter int

func (*KafkaFilter) DoFilter(data interface{}) {
	// fmt.Printf("data %+v\n", data)
	message, isOk := data.(kafka.Message)
	if !isOk {
		log.Error().Msg("非法kafka消息")
		return
	}
	var packet reqcontext.RequestPacket
	err := json.Unmarshal(message.Value, &packet)
	if err != nil {
		log.Error().Err(err).Msg("解析json失败" + message.String())
		return
	}
	msgType := packet.Header.MsgType
	config, ok := route.RouteMap[msgType]
	if !ok {
		log.Error().Msg(msgType + "类型消息不被支持")
		return
	}
	context := reqcontext.KafkaRequestContext{}
	reqcontext.InitKafkaRequestContext(&context, &packet)

	helper.CallService(&context, &config)
}
