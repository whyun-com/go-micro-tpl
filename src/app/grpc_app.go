package app

import (
	"context"
	"net"
	"strconv"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"github.com/whyun-com/go-micro-tpl/filter"
	"github.com/whyun-com/go-micro-tpl/grpc_health_v1"
	"github.com/whyun-com/go-micro-tpl/micro"
)

const (
	grpcPort = 9100
)

type GrpcApp struct {
	Port uint32
}

type server struct {
	micro.UnimplementedMicroServer
	grpc_health_v1.UnimplementedHealthServer
}

func (s *server) DoRequest(ctx context.Context, in *micro.MicroRequest) (*micro.MicroResponse, error) {
	return filter.DoGrpcFilter(in)
}

func (s *server) Check(context.Context, *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}
func (s *server) Watch(*grpc_health_v1.HealthCheckRequest, grpc_health_v1.Health_WatchServer) error {
	return nil
}

func (app GrpcApp) Start() error {
	var address string
	if app.Port == 0 {
		address = ":" + strconv.Itoa(grpcPort)
	} else {
		address = ":" + strconv.Itoa(int(app.Port))
	}
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Error().Err(err).Msg("监听grpc端口失败")
		return err
	}
	s := grpc.NewServer()
	micro.RegisterMicroServer(s, &server{})
	grpc_health_v1.RegisterHealthServer(s, &server{})
	log.Info().Msg("grpc服务端口监听完成:" + lis.Addr().String())
	return s.Serve(lis)
}
