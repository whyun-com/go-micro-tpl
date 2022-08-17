package filter

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/whyun-com/go-micro-tpl/helper"
	"github.com/whyun-com/go-micro-tpl/micro"
	"github.com/whyun-com/go-micro-tpl/reqcontext"
	"github.com/whyun-com/go-micro-tpl/route"
)

func genError(context *reqcontext.GrpcRequestContext, result *reqcontext.Result) *micro.MicroResponse {
	reqcontext.DoError(context, &reqcontext.MicroResponse{
		Result: *result,
	}, nil)
	return &micro.MicroResponse{
		Result: &micro.Result{
			Code: uint32(result.Code),
			Msg:  result.Msg,
		},
	}
}

func DoGrpcFilter(in *micro.MicroRequest) (*micro.MicroResponse, error) {
	var header *micro.MessageHeader = in.Header
	var user *micro.User = in.User

	var packet reqcontext.RequestPacket = reqcontext.RequestPacket{
		MicroRequest: reqcontext.MicroRequest{
			Header: reqcontext.MessageHeader{
				MsgType:        header.MsgType,
				ReqMs:          header.ReqMs,
				RequestId:      header.RequestId,
				ReqPid:         header.ReqPid,
				ReqIp:          header.ReqIp,
				ExtendedFields: header.ExtendedFields,
			},
			User: reqcontext.User{
				UserId: user.UserId,
			},
		},
	}
	context := reqcontext.GrpcRequestContext{}
	reqcontext.InitGrpcRequestContext(&context, &packet)
	reqStr := in.ReqDataJson
	err := json.Unmarshal([]byte(reqStr), &packet.ReqData)
	if err != nil {
		log.Error().Err(err).Msg("反序列化请求数据失败:" + reqStr)
		return genError(&context, &reqcontext.Result{ //TODO 增加错误码
			Code: 2000,
		}), nil
	}
	msgType := header.MsgType
	config, ok := route.RouteMap[msgType]
	if !ok {
		log.Error().Msg(msgType + "类型消息不被支持")
		return genError(&context, &reqcontext.Result{ //TODO 增加错误码
			Code: 2001,
		}), nil
	}
	resContent := helper.CallService(&context, &config)
	return &micro.MicroResponse{
		Header: &micro.MessageHeader{
			MsgType:        resContent.Header.MsgType,
			ReqMs:          header.ReqMs,
			RequestId:      header.RequestId,
			ReqPid:         header.ReqPid,
			ReqIp:          header.ReqIp,
			ExtendedFields: resContent.Header.ExtendedFields,
		},
		Result: &micro.Result{
			Code: uint32(resContent.Result.Code),
			Msg:  resContent.Result.Msg,
		},
		ResDataJson: resContent.ResData,
	}, nil
}
