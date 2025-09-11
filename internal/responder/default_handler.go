package responder

import (
	"context"
	"message-responder/internal/contract"
)

type DefaultHandler struct{}

func NewDefaultHandler() *DefaultHandler {
	return &DefaultHandler{}
}

func (h DefaultHandler) TryHandle(ctx context.Context, msg contract.NormalizedRequest) (contract.HandlerRequest, bool, error) {
	return contract.HandlerRequest{Text: "Отправьте одно изображение. Мы извлечём из него текст и вернём результат."}, true, nil
}
