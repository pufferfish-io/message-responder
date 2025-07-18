package handler

import (
	mdl "message-responder/internal/model"
	mdlResp "message-responder/internal/responder/model"
)

type DefaultHandler struct{}

func (h DefaultHandler) TryHandle(msg mdl.NormalizedRequest) (mdlResp.HandlerRequest, bool, error) {
	return mdlResp.HandlerRequest{Text: "НЕ ПОДОШЛО"}, true, nil
}
