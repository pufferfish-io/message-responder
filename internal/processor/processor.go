package processor

import (
	"context"
	"encoding/json"
	"fmt"

	"message-responder/internal/contract"
	"message-responder/internal/logger"
)

type Producer interface {
	Send(_ context.Context, topic string, data []byte) error
}

type Handler interface {
	TryHandle(ctx context.Context, msg contract.NormalizedRequest) (contract.HandlerRequest, bool, error)
}

type MessageResponder struct {
	logger     logger.Logger
	kafkaTopic string
	producer   Producer
	handlers   []Handler
}

func NewMessageResponder(kafkaTopic string, producer Producer, h []Handler, lg logger.Logger) *MessageResponder {
	return &MessageResponder{
		kafkaTopic: kafkaTopic,
		producer:   producer,
		handlers:   h,
		logger:     lg,
	}
}

func (t *MessageResponder) Handle(ctx context.Context, raw []byte) error {
	var requestMessage contract.NormalizedRequest
	if err := json.Unmarshal(raw, &requestMessage); err != nil {
		t.logger.Error("unmarshal request error: %v", err)
		return err
	}

	textPresent := requestMessage.Text != nil && len(*requestMessage.Text) > 0
	t.logger.Debug(
		"handling message: chat_id=%d user_id=%d text_present=%t media_count=%d",
		requestMessage.ChatID, requestMessage.UserID, textPresent, len(requestMessage.Media),
	)

	response := contract.NormalizedResponse{
		ChatID:         requestMessage.ChatID,
		Source:         requestMessage.Source,
		UserID:         requestMessage.UserID,
		Username:       requestMessage.Username,
		Timestamp:      requestMessage.Timestamp,
		OriginalUpdate: requestMessage.OriginalUpdate,
	}

	matched := ""
	for _, h := range t.handlers {
		t.logger.Debug("trying handler: %T", h)
		result, ok, err := h.TryHandle(ctx, requestMessage)
		if ok {
			if err != nil {
				response.Text = "На сервере произошла ошибка, мы уже начали исправлять данную проблему."
				t.logger.Error("handler error (used fallback text): %T err=%v", h, err)
			} else {
				response.Text = result.Text
				matched := fmt.Sprintf("%T", h)
				t.logger.Info("matched handler: %s", matched)
			}
			break
		}
	}

	out, err := json.Marshal(response)
	if err != nil {
		t.logger.Error("marshal response error: %v", err)
		return err
	}
	topik := requestMessage.Source + t.kafkaTopic
	t.logger.Debug("sending response: topic=%s bytes=%d matched=%s", topik, len(out), matched)
	if err := t.producer.Send(ctx, topik, out); err != nil {
		t.logger.Error("producer send error: %v", err)
		return err
	}
	t.logger.Info("response sent: topic=%s", topik)
	return nil
}
