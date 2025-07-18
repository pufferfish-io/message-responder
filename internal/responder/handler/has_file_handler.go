package handler

import (
	mdl "message-responder/internal/model"
	mdlResp "message-responder/internal/responder/model"
)

type HasFileHandler struct{}

func (h HasFileHandler) TryHandle(msg mdl.NormalizedRequest) (mdlResp.HandlerRequest, bool, error) {
	if len(msg.Media) > 0 {
		return mdlResp.HandlerRequest{Text: "BOBA"}, true, nil
	}
	return mdlResp.HandlerRequest{}, false, nil
}
