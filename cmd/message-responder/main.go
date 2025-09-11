package main

import (
	"context"
	"log"
	cfg "message-responder/internal/config"
	cfgModel "message-responder/internal/config/model"
	"message-responder/internal/messaging"
	"message-responder/internal/processor"
	"message-responder/internal/responder"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Println("🚀 Starting message-responder...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	kafkaConf, err := cfg.LoadSection[cfgModel.KafkaConfig]()
	if err != nil {
		log.Fatalf("❌ Failed to load Kafka config: %v", err)
	}

	if err != nil {
		log.Fatalf("❌ Failed to create S3 uploader: %v", err)
	}

	respondProcessor := responder.NewRespondProcessor()
	normalizer := processor.NewMessageResponder(kafkaConf.ResponseMessageTopicName, respondProcessor)

	handler := func(msg []byte) error {
		return normalizer.Handle(ctx, msg)
	}

	messaging.Init(kafkaConf.BootstrapServersValue)

	consumer, err := messaging.NewConsumer(kafkaConf.BootstrapServersValue, kafkaConf.RequestMessageGroupId, kafkaConf.RequestMessageTopicName, handler)
	if err != nil {
		log.Fatalf("❌ Failed to create consumer: %v", err)
	}

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		cancel()
	}()

	if err := consumer.Start(ctx); err != nil {
		log.Fatalf("❌ Consumer error: %v", err)
	}
}
