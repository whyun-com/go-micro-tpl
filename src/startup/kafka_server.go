package startup

import "github.com/whyun-com/go-micro-tpl/app"

type KafkaServer struct {
}

func (server KafkaServer) Init(config ServerConfig) {

	app4Server := app.KafkaApp{}
	go app4Server.Start()
}
