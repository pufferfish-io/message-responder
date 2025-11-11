package main

import (
	"context"
	"errors"
	"message-responder/internal/config"
	"message-responder/internal/logger"
	"message-responder/internal/messaging"
	"message-responder/internal/processor"
	"message-responder/internal/responder"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	lg, cleanup := logger.NewZapLogger()
	defer cleanup()

	lg.Info("🚀 Starting message-responder…")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.Load()
	if err != nil {
		lg.Error("❌ Failed to load config: %v", err)
		os.Exit(1)
	}

	producer, err := messaging.NewKafkaProducer(messaging.Option{
		Logger:       lg,
		Broker:       cfg.Kafka.BootstrapServersValue,
		SaslUsername: cfg.Kafka.SaslUsername,
		SaslPassword: cfg.Kafka.SaslPassword,
		ClientID:     cfg.Kafka.ClientID,
		Context:      ctx,
	})
	if err != nil {
		lg.Error("❌ Failed to create producer: %v", err)
		os.Exit(1)
	}
	defer producer.Close()

	var handlers = []processor.Handler{
		responder.NewHasFileHandler(cfg.Kafka.OcrTopicName, producer),
		responder.NewDefaultHandler(),
	}
	responder := processor.NewMessageResponder(cfg.Kafka.ResponseTopicName, producer, handlers, lg)

	consumer, err := messaging.NewKafkaConsumer(messaging.ConsumerOption{
		Logger:       lg,
		Broker:       cfg.Kafka.BootstrapServersValue,
		GroupID:      cfg.Kafka.GroupID,
		Topics:       []string{cfg.Kafka.RequestTopicName},
		Handler:      responder,
		SaslUsername: cfg.Kafka.SaslUsername,
		SaslPassword: cfg.Kafka.SaslPassword,
		ClientID:     cfg.Kafka.ClientID,
		Context:      ctx,
	})
	if err != nil {
		lg.Error("❌ Failed to create consumer: %v", err)
		os.Exit(1)
	}

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		cancel()
	}()

	runConsumerSupervisor(ctx, consumer, lg)
}

func runConsumerSupervisor(ctx context.Context, consumer *messaging.KafkaConsumer, lg logger.Logger) {
	wait := 1 * time.Second
	maxWait := 30 * time.Second
	for {
		err := consumer.Start(ctx)
		if err == nil {
			lg.Info("consumer finished without error")
			return
		}
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			lg.Info("consumer stopped: %v", err)
			return
		}
		lg.Error("consumer error: %v — retry in %s", err, wait)
		select {
		case <-time.After(wait):
			wait *= 2
			if wait > maxWait {
				wait = maxWait
			}
		case <-ctx.Done():
			lg.Info("context canceled while waiting to restart consumer: %v", ctx.Err())
			return
		}
	}
}
