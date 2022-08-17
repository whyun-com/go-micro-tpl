package client

import (
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rs/zerolog/log"
	"github.com/whyun-com/go-micro-tpl/reqcontext"
)

type KafkaClient struct {
	Producer *kafka.Producer
	Topic    string
}

type KafkaOption struct {
	KafkaConf *kafka.ConfigMap
	Topic     string
}

func NewKafkaClient(kafkaOption *KafkaOption) (*KafkaClient, error) {
	client := &KafkaClient{Topic: kafkaOption.Topic}
	producer, err := kafka.NewProducer(kafkaOption.KafkaConf)
	if err != nil {
		return nil, err
	}
	client.Producer = producer
	return client, nil
}

func (client *KafkaClient) Send(req *reqcontext.RequestPacket, res chan<- reqcontext.MicroResponse) {

	deliveryChan := make(chan kafka.Event)
	response := reqcontext.GetDefaultMicroResponse()
	defer func() {
		res <- response
		close(deliveryChan)
	}()
	bytes, err := json.Marshal(req.MicroRequest)
	if err != nil {
		response.Result.Code = 1000
		log.Error().Err(err).Msg("序列化json失败")
		return
	}
	client.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &client.Topic, Partition: kafka.PartitionAny},
		Value:          bytes,
		// Headers:        []kafka.Header{{Key: "myTestHeader", Value: []byte("header values are binary")}},
	}, deliveryChan)

	event := <-deliveryChan
	message := event.(*kafka.Message)

	if message.TopicPartition.Error != nil {
		// fmt.Printf("Delivery failed: %v\n", message.TopicPartition.Error)
		response.Result.Code = 1001
		log.Error().Err(message.TopicPartition.Error).Msg("发送数据到kafka失败" + string(bytes))
	} else {
		// fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
		// 	*message.TopicPartition.Topic, message.TopicPartition.Partition, message.TopicPartition.Offset)
		response.Result.Code = 0
	}

}
