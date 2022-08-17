package startup

import (
	"github.com/rs/zerolog/log"
	"github.com/whyun-com/go-micro-tpl/app"
)

type GrpcServer struct {
}

func (server GrpcServer) Init(config ServerConfig) {

	app4Server := app.GrpcApp{
		Port: config.GrpcPort,
	}
	go func() {
		err := app4Server.Start()
		if err != nil {
			log.Error().Err(err).Msg("grpc服务出错")
		}
	}()
}
