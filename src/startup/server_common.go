package startup

type ServerInterface interface {
    Init(config ServerConfig)
    Destory()
}

type ServerConfig struct {
    HttpPort uint32
    GrpcPort uint32
    KafkaDisabled bool
}
