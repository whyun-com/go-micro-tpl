package client

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	// "github.com/hashicorp/consul/api"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"github.com/whyun-com/go-micro-tpl/consul"
	"github.com/whyun-com/go-micro-tpl/micro"
	"github.com/whyun-com/go-micro-tpl/reqcontext"
)

type ServcieConfig struct {
	ServiceName           string
	ServicePort           uint
	InitTimeoutSeconds    uint
	RequestTimeoutSeconds uint
}
type GrpcClient struct {
	conn                  *grpc.ClientConn
	requestTimeoutSeconds uint
}

const DEFAULT_INIT_TIMEOUT_SECONDS = 30
const DEFAULT_REQUEST_TIMEOUT_SECONDS = 3

func init() {
	fmt.Println("init grpc client")
	consul.Init()
}

func NewGrpcClient(conf *ServcieConfig) (*GrpcClient, error) {
	grpcClient := &GrpcClient{
		requestTimeoutSeconds: conf.RequestTimeoutSeconds,
	}
	if grpcClient.requestTimeoutSeconds == 0 {
		grpcClient.requestTimeoutSeconds = DEFAULT_REQUEST_TIMEOUT_SECONDS
	}
	initTimeoutSeconds := conf.InitTimeoutSeconds
	if initTimeoutSeconds == 0 {
		initTimeoutSeconds = DEFAULT_INIT_TIMEOUT_SECONDS
	}


	target := fmt.Sprintf("%s:%d", conf.ServiceName, conf.ServicePort)
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(initTimeoutSeconds)*time.Second) //TODO CANCEL
	conn, err := grpc.DialContext(ctx, target, grpc.WithBlock(), grpc.WithInsecure(), grpc.WithBalancerName("round_robin"))
	if err != nil {
		log.Error().Err(err).Msg("初始化grpc客户端错误")
		return nil, err
	}
	grpcClient.conn = conn

	return grpcClient, nil

}

func (client *GrpcClient) Send(req *reqcontext.RequestPacket, res chan<- reqcontext.MicroResponse) {
	response := reqcontext.GetDefaultMicroResponse()
	defer func() {
		res <- response
	}()
	instance := micro.NewMicroClient(client.conn)
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(client.requestTimeoutSeconds)*time.Second) //TODO CANCEL
	originalReq := req.MicroRequest
	bytes, err := json.Marshal(originalReq.ReqData)
	if err != nil {
		response.Result.Code = 2001 //TODO defined error code
		log.Warn().Err(err).Msg("序列化请求数据失败")
		return
	}
	header := originalReq.Header
	user := originalReq.User

	body, err := instance.DoRequest(ctx, &micro.MicroRequest{
		Header: &micro.MessageHeader{
			MsgType: header.MsgType,
			ReqMs:   header.ReqMs,
			RequestId:   header.RequestId, //TODO 更多字段

		},
		User: &micro.User{
			UserId: user.UserId,
		},
		ReqDataJson: string(bytes),
	})
	if err != nil {
		response.Result.Code = 2000 //TODO 定义错误码
	} else {
		response.Result.Code = 0
		response.ResData = body.ResDataJson
	}
}
