package config

import (
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rs/zerolog/log"
)

type KafkaConfig struct {
}

var (
	LogProducer *kafka.Producer
)

const (
	ENV_KAFKA_LOG_BROKERS             = "KAFKA_LOG_BROKERS"
	ENV_KAFKA_COMMUNICATION_BROKERS   = "KAFKA_COMMUNICATION_BROKERS"
	ENV_KAFKA_COMMUNICATION_TOPIC     = "KAFKA_COMMUNICATION_TOPIC"
	ENV_KAFKA_LOG_TOPICS_ACCESS_LOG   = "KAFKA_LOG_TOPICS_ACCESS_LOG"
	DEFAULT_KAFKA_COMMUNICATION_TOPIC = "go-micro-tpl-queue"
	DEFAULT_KAFKA_ACCESS_LOG_TOPIC    = "go-micro-tpl-access-log"
)

func getProducerConf(bootstrap string) *kafka.ConfigMap {
	return &kafka.ConfigMap{
		"bootstrap.servers": bootstrap,
	}
}
func alarmKafkaError(p *kafka.Producer) {
	for e := range p.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				// fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				log.Error().Err(ev.TopicPartition.Error).Msg(
					*ev.TopicPartition.Topic + "发送数据失败" + string(ev.Value),
				)
				//TODO 发送报警信息
			}
		}
	}
}
func (config *KafkaConfig) Init() {
	logHostsFromEnv := os.Getenv(ENV_KAFKA_LOG_BROKERS)

	if logHostsFromEnv != "" {
		Config.Kafka.Log.Brokers = logHostsFromEnv
	}
	communicationHostsFromEnv := os.Getenv(ENV_KAFKA_COMMUNICATION_BROKERS)
	if communicationHostsFromEnv != "" {
		Config.Kafka.Communication.Brokers = communicationHostsFromEnv
	}
	if Config.Kafka.Log.Brokers == "" {
		log.Fatal().Msg("日志用kafka地址为空")
		return
	}
	logProducer, err := kafka.NewProducer(getProducerConf(Config.Kafka.Log.Brokers))
	if err != nil {
		log.Fatal().Err(err).Msg("初始化 log kafka失败")
		return
	}
	LogProducer = logProducer
	if len(Config.Kafka.Communication.Brokers) == 0 &&
		os.Getenv("KAFKA_STARTUP_DISABLED") != "true" {
		log.Fatal().Msg("通信用kafka地址为空")
		return
	}
	accessToic := os.Getenv(ENV_KAFKA_LOG_TOPICS_ACCESS_LOG)
	if accessToic != "" {
		Config.Kafka.Log.Topics.AccessLog = accessToic
	} else if Config.Kafka.Log.Topics.AccessLog == "" {
		Config.Kafka.Log.Topics.AccessLog = DEFAULT_KAFKA_ACCESS_LOG_TOPIC
	}
	commTopic := os.Getenv(ENV_KAFKA_COMMUNICATION_TOPIC)
	if commTopic != "" {
		Config.Kafka.Communication.Topic = commTopic
	} else if Config.Kafka.Communication.Topic == "" {
		Config.Kafka.Communication.Topic = DEFAULT_KAFKA_COMMUNICATION_TOPIC
	}
}
