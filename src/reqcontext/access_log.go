package reqcontext

import (
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rs/zerolog/log"
	"github.com/whyun-com/go-micro-tpl/config"
)

/* 消费kafka日志打点要求
{
    req_type:Number,//请求类型，0: http, 1: kafka, 2: grpc
    is_ok: Boolean,//是否成功
    res_code: Number,//返回码
    begin: Number,//开始处理时间戳
    end: Number,//结束处理时间戳
    req_id: String,//请求ID
    pid: Number,//当前进程ID
    header: Object,//参见 protobuf 中的  MessageHeader message
    req_data: Object,//请求正文
    meta_headers: Object,// kafka 的头信息，具体驱动不同其字段可能不同，这里 node 使用的是 node-rdkafka，其 timestamp 字段代表写入 kafka 成功的时间戳
    host: String,// 当前服务IP
    duration: Number,// 处理时长
    end_minus_begin: Number,// end - time 的值
    begin_minus_request: Number,// begin - header.req_ms 的值
    begin_minus_kafka: Number,// begin - meta_headers.timestamp 的值
    request_time_str: String,// header.req_ms 转成的 ISO 时间字符串
    begin_time_str: String,// begin 转成的 ISO 时间字符串
    user: Object //参见protobuf的 User message
}
*/

type AccessLog struct {
	IsOk              bool                   `json:"is_ok,omitempty"`
	ReqType           string                 `json:"req_type,omitempty"`
	ResCode           int                    `json:"res_code,omitempty"`
	Begin             uint64                 `json:"begin,omitempty"`
	BeginTimeStr      string                 `json:"begin_time_str,omitempty"`
	End               uint64                 `json:"end,omitempty"`
	EndTimeStr        string                 `json:"end_time_str,omitempty"`
	RequestTimeStr    string                 `json:"request_time_str,omitempty"`
	RequestId         string                 `json:"req_id,omitempty"`
	Header            MessageHeader          `json:"header,omitempty"`
	User              User                   `json:"user,omitempty"`
	ReqData           map[string]interface{} `json:"req_data,omitempty"`
	MetaHeaders       interface{}            `json:"meta_headers,omitempty"`
	ClientIp          string                 `json:"client_ip,omitempty"`
	Host              string                 `json:"host"`
	Duration          int64                  `json:"duration"`
	EndMinusBegin     int64                  `json:"end_minus_begin,omitempty"`
	BeginMinusRequest int64                  `json:"begin_minus_request,omitempty"`
	BeginMinusKafka   int64                  `json:"begin_minus_kafka,omitempty"`
	ResData           string                 `json:"res_data,omitempty"`
}

var accessLogTopic = config.DEFAULT_KAFKA_ACCESS_LOG_TOPIC

func (aclog *AccessLog) DoSend() {
	bytes, err := json.Marshal(aclog)
	if err != nil {
		log.Error().Err(err).Msg("当前日志不能被正常序列化")
		return
	}
	config.LogProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &accessLogTopic, Partition: kafka.PartitionAny},
		Value:          bytes,
	}, nil)
}
