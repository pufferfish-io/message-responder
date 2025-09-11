package processor

import (
	"context"
	"encoding/json"

	"message-responder/internal/messaging"
	mdl "message-responder/internal/model"
	"message-responder/internal/responder"
)

type MessageResponder struct {
	kafkaTopic       string
	respondProcessor *responder.RespondProcessor
}

func NewMessageResponder(kafkaTopic string, respondProcessor *responder.RespondProcessor) *MessageResponder {
	return &MessageResponder{
		kafkaTopic:       kafkaTopic,
		respondProcessor: respondProcessor,
	}
}

func (t *MessageResponder) Handle(ctx context.Context, raw []byte) error {
	var requestMessage mdl.NormalizedRequest
	if err := json.Unmarshal(raw, &requestMessage); err != nil {
		return err
	}

	var response, err = t.respondProcessor.ProcessMessage(requestMessage)
	if err != nil {
		return err
	}

	out, err := json.Marshal(response)
	if err != nil {
		return err
	}

	return messaging.Send(t.kafkaTopic, out)
}
