package handler

import (
	mdl "message-responder/internal/model"
	mdlResp "message-responder/internal/responder/model"
	"strings"
)

type ContainsBibaHandler struct{}

func (h ContainsBibaHandler) TryHandle(msg mdl.NormalizedRequest) (mdlResp.HandlerRequest, bool, error) {
	if msg.Text != nil && strings.Contains(*msg.Text, "BIBA") {
		return mdlResp.HandlerRequest{Text: "BOBA"}, true, nil
	}
	return mdlResp.HandlerRequest{}, false, nil
}
