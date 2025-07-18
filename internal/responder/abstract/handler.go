package processing

import (
	mdl "message-responder/internal/model"
	mdlResp "message-responder/internal/responder/model"
)

type Handler interface {
	TryHandle(msg mdl.NormalizedRequest) (mdlResp.HandlerRequest, bool, error)
}
