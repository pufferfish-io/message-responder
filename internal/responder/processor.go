package responder

import (
	mdl "message-responder/internal/model"
	absrct "message-responder/internal/responder/abstract"
	hdlr "message-responder/internal/responder/handler"
)

type RespondProcessor struct {
}

func NewRespondProcessor() *RespondProcessor {
	return &RespondProcessor{}
}

func (RespondProcessor) ProcessMessage(msg mdl.NormalizedRequest) (mdl.NormalizedResponse, error) {
	handlers := []absrct.Handler{
		hdlr.HasFileHandler{},
		hdlr.ContainsBibaHandler{},
		hdlr.DefaultHandler{},
	}

	for _, handler := range handlers {
		result, ok, err := handler.TryHandle(msg)
		if err != nil {
			return mdl.NormalizedResponse{}, err
		}
		if ok {
			return mdl.NormalizedResponse{
				ChatID: msg.ChatID,
				Text:   &result.Text,
			}, nil
		}
	}

	return mdl.NormalizedResponse{}, nil
}
