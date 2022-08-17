package startup

import (
	"github.com/rs/zerolog/log"
	"github.com/whyun-com/go-micro-tpl/app"
)

type HttpServer struct {
}

func (server HttpServer) Init(config ServerConfig) {

	app4Server := app.HttpApp{
		Port: config.HttpPort,
	}
	go func() {
		err := app4Server.Start()
		if err != nil {
			log.Error().Err(err).Msg("http服务出错")
		}
	}()
}
