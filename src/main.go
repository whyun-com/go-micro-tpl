package main

import (
	"os"
	"strconv"

	"github.com/rs/zerolog/log"
	_ "go.uber.org/automaxprocs"
	"github.com/whyun-com/go-micro-tpl/startup"
)

func main() {

	httpPort, _ := strconv.ParseUint(os.Getenv("HTTP_PORT"), 10, 32)
	grpcPort, _ := strconv.ParseUint(os.Getenv("GRPC_PORT"), 10, 32)
	kafkaDisabled := os.Getenv("KAFKA_STARTUP_DISABLED") == "true"
	serverConfig := startup.ServerConfig{
		HttpPort:      uint32(httpPort),
		GrpcPort:      uint32(grpcPort),
		KafkaDisabled: kafkaDisabled,
	}

	if !kafkaDisabled {
		log.Info().Msg("start kafka server now...")
		kafkaServer := startup.KafkaServer{}
		kafkaServer.Init(serverConfig)
	} else {
		log.Info().Msg("ignore kafka server")
	}

	if httpPort > 0 {
		log.Info().Msg("start http server now...")
		httpServer := startup.HttpServer{}
		httpServer.Init(serverConfig)
	} else {
		log.Info().Msg("ignore http server")
	}

	if grpcPort > 0 {
		log.Info().Msg("start grpc server now...")
		grpcServer := startup.GrpcServer{}
		grpcServer.Init(serverConfig)
	} else {
		log.Info().Msg("ignore grpc server")
	}

	log.Info().Msg("main begin")
	ch := make(chan int, 1)
	<-ch
}
