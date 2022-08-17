package app

import (
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/panjf2000/ants/v2"
	"github.com/whyun-com/go-micro-tpl/config"
	"github.com/whyun-com/go-micro-tpl/filter"
)

type KafkaApp struct {
}

func (app KafkaApp) Start() error {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": config.Config.Kafka.Communication.Brokers,
		// Avoid connecting to IPv6 brokers:
		// This is needed for the ErrAllBrokersDown show-case below
		// when using localhost brokers on OSX, since the OSX resolver
		// will return the IPv6 addresses first.
		// You typically don't need to specify this configuration property.
		"broker.address.family": "v4",
		"group.id":              "go-micro-tpl", //修改组ID
		"session.timeout.ms":    16000,
		"auto.offset.reset":     "earliest",
	})
	if err != nil {
		return err
	}
	filterKafka := filter.KafkaFilter(filter.FILTER_TYPE_KAFKA)
	pool, _ := ants.NewPoolWithFunc(1024*100, filterKafka.DoFilter)
	defer pool.Release()

	err = consumer.SubscribeTopics([]string{config.Config.Kafka.Communication.Topic}, nil)
	if err != nil {
		return err
	}
	run := true

	for run {
		event := consumer.Poll(100)
		switch data := event.(type) {
		case *kafka.Message:
			err := pool.Invoke(data)
			if err != nil {
				fmt.Fprintf(os.Stderr, "提交数据到池失败: %v : %v", err, data)
			}
		case kafka.Error:
			// Errors should generally be considered
			// informational, the client will try to
			// automatically recover.
			// But in this example we choose to terminate
			// the application if all brokers are down.
			fmt.Fprintf(os.Stderr, "%% cunsumer Error: %v: %v\n", data.Code(), data)
			//即使redis宕机也要重试
			// if data.Code() == kafka.ErrAllBrokersDown {
			// 	run = false
			// }
		default:
			// fmt.Printf("Ignored %v\n", data)
		}
	}
	return nil
}
