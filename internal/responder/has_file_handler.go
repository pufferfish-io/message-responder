package responder

import (
	"context"
	"encoding/json"
	"message-responder/internal/contract"
	"strings"
)

type Producer interface {
	Send(_ context.Context, topic string, data []byte) error
}

type HasFileHandler struct {
	ocrTopicName string
	producer     Producer
}

func NewHasFileHandler(t string, p Producer) *HasFileHandler {
	return &HasFileHandler{ocrTopicName: t, producer: p}
}

func (h HasFileHandler) TryHandle(ctx context.Context, msg contract.NormalizedRequest) (contract.HandlerRequest, bool, error) {
    if hasImage(msg) {
		ocrRequest := contract.OcrRequest{
			Source:         msg.Source,
			UserID:         msg.UserID,
			Username:       msg.Username,
			ChatID:         msg.ChatID,
			Timestamp:      msg.Timestamp,
			Media:          msg.Media,
			OriginalUpdate: msg.OriginalUpdate,
		}
        out, err := json.Marshal(ocrRequest)
        if err != nil {
            return contract.HandlerRequest{}, true, err
        }
        if err := h.producer.Send(ctx, h.ocrTopicName, out); err != nil {
            // Сообщаем процессору об ошибке, чтобы он отдал fallback-текст
            return contract.HandlerRequest{}, true, err
        }
        return contract.HandlerRequest{Text: "Ваше изображение отправлено на обработку, ожидайте."}, true, nil
    }
    return contract.HandlerRequest{}, false, nil
}

func hasImage(msg contract.NormalizedRequest) bool {
	if len(msg.Media) == 0 {
		return false
	}
	for _, m := range msg.Media {
		if m.MimeType != nil && strings.HasPrefix(strings.ToLower(*m.MimeType), "image/") {
			return true
		}
		t := strings.ToLower(m.Type)
		switch t {
		case "image", "photo", "picture", "img":
			return true
		}
	}
	return false
}
