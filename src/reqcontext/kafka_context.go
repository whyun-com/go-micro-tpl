package reqcontext
import (
	"github.com/rs/zerolog/log"
)
type KafkaRequestContext struct {
	BaseContext
}
func InitKafkaRequestContext(req *KafkaRequestContext, packet *RequestPacket) {
	initBaseContext(&req.BaseContext, packet)
	req.RequestType = REQUEST_TYPE_KAFKA
}

func (req *KafkaRequestContext) doSuccess(data string, conf *ServiceConfig) {
	log.Debug().Msg("DoSuccess in KafkaRequestContext" + data)
}

func (req *KafkaRequestContext) doError(result Result, conf *ServiceConfig) {
	log.Debug().Int("code", int(result.Code)).Msg("DoError in KafkaRequestContext")
}