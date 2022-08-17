package reqcontext
import (
	"github.com/rs/zerolog/log"
)
type GrpcRequestContext struct {
	BaseContext
}
func InitGrpcRequestContext(req *GrpcRequestContext, packet *RequestPacket) {
	initBaseContext(&req.BaseContext, packet)
	req.RequestType = REQUEST_TYPE_GRPC
}

func (req *GrpcRequestContext) doSuccess(data string, conf *ServiceConfig) {
	log.Debug().Msg("DoSuccess in GrpcRequestContext" + data)
}

func (req *GrpcRequestContext) doError(result Result, conf *ServiceConfig) {
	log.Debug().Int("code", int(result.Code)).Msg("DoError in GrpcRequestContext")
}