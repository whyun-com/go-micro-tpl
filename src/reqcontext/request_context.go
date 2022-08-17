package reqcontext

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

const REQUEST_TYPE_HTTP string = "http"
const REQUEST_TYPE_KAFKA string = "kafka"
const REQUEST_TYPE_GRPC string = "grpc"

type User struct {
	UserId string `json:"user_id"`
}

type MessageHeader struct {
	MsgType        string            `json:"msg_type"`
	ReqMs          uint64            `json:"req_ms"`
	RequestId      string            `json:"req_id"`
	ReqPid         uint32            `json:"req_pid"`
	ReqIp          string            `json:"req_ip"`
	UserpoolId     string            `json:"userpool_id"`
	AppId          string            `json:"app_id,omitempty"`
	ExtendedFields map[string]string `json:"extended_fields"`
}

type MicroRequest struct {
	Header  MessageHeader          `json:"header"`
	User    User                   `json:"user"`
	ReqData map[string]interface{} `json:"req_data"`
}

type RequestPacket struct {
	MicroRequest
	ProtocolHeader interface{} //协议相关头信息
}

type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type MicroResponse struct {
	Result  Result
	Header  MessageHeader
	ResData string //json字符串
}

type ContextInterface interface {
	doSuccess(data string, conf *ServiceConfig)
	doError(result Result, conf *ServiceConfig)
	doEnd(content *MicroResponse)
	getReqData() map[string]interface{}
	getMessageHeader() MessageHeader
}

type BaseContext struct {
	Begin          uint64
	BeginTimeStr   string
	End            uint64
	EndTimeStr     string
	EndMinusBegin  uint32
	RequestType    string
	WaitGroupNow   sync.WaitGroup
	User           User
	ReqData        map[string]interface{}
	Header         MessageHeader
	IsOk           bool
	ReqId          uint64
	ProtocolHeader interface{} //协议相关头信息
}

func (result Result) String() string {
	return fmt.Sprintf("[%d] %s", result.Code, result.Msg)
}

type ServiceFun func(req ContextInterface, res chan<- MicroResponse)

type ServiceOption struct {
	MaxAge uint32
}
type ServiceConfig struct {
	Service ServiceFun
	Option  ServiceOption
}

func initBaseContext(req *BaseContext, packet *RequestPacket) {
	now := time.Now()
	req.Begin = uint64(now.UnixNano() / 1e6)
	req.BeginTimeStr = now.Format(time.RFC3339)
	req.ReqData = packet.ReqData
	req.Header = packet.Header
	req.User = packet.User
	req.ProtocolHeader = packet.ProtocolHeader
	req.ReqId++
}

func (req *BaseContext) doEnd(content *MicroResponse) {
	now := time.Now()
	req.End = uint64(now.UnixNano() / 1e6)
	req.EndTimeStr = now.Format(time.RFC3339)
	req.EndMinusBegin = uint32(req.End - req.Begin)
	result := content.Result
	header := req.Header
	requestMs := header.ReqMs
	requestTimeStr := time.Unix(
		0, int64(requestMs*uint64(time.Microsecond)),
	).Format(time.RFC3339)

	aclog := &AccessLog{
		ReqType:           req.RequestType,
		IsOk:              result.Code == 0,
		ResCode:           result.Code,
		Begin:             req.Begin,
		BeginTimeStr:      req.BeginTimeStr,
		End:               req.End,
		EndTimeStr:        req.EndTimeStr,
		RequestTimeStr:    requestTimeStr,
		RequestId:         strconv.FormatUint(req.ReqId, 10),
		Header:            header,
		User:              req.User,
		ReqData:           req.ReqData,
		MetaHeaders:       req.ProtocolHeader,
		EndMinusBegin:     int64(req.End - req.Begin),
		BeginMinusRequest: int64(req.Begin - requestMs),
		ResData:           content.ResData,
	}
	aclog.DoSend()
}
func (req *BaseContext) doSuccess(data string, conf *ServiceConfig) {
	panic("not support dosuccess in base")
}

func DoSuccess(req ContextInterface, content *MicroResponse, conf *ServiceConfig) {
	req.doSuccess(content.ResData, conf)
	req.doEnd(content)
}

func DoError(req ContextInterface, content *MicroResponse, conf *ServiceConfig) {
	req.doError(content.Result, conf)
	req.doEnd(content)
}

func (req *BaseContext) doError(result Result, conf *ServiceConfig) {
	panic("not support doerror in base")
}

func (req *BaseContext) getReqData() map[string]interface{} {
	return req.ReqData
}

func GetReqData(req ContextInterface) map[string]interface{} {
	return req.getReqData()
}

func (req *BaseContext) getMessageHeader() MessageHeader {
	return req.Header
}

func GetMessageHeader(req ContextInterface) MessageHeader {
	return req.getMessageHeader()
}

func GetDefaultMicroResponse() MicroResponse {
	return MicroResponse{
		Result: Result{Code: 1}, //TODO 定义全局变量，承载未知错误
	}
}
