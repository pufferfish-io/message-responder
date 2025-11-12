package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/go-playground/validator/v10"
)

type Kafka struct {
	BootstrapServersValue string `validate:"required" env:"BOOTSTRAP_SERVERS_VALUE"`
	GroupID               string `validate:"required" env:"GROUP_ID_MESSAGE_RESPONDER"`
	RequestTopicName      string `validate:"required" env:"TOPIC_NAME_NORMALIZED_MSG"`
	ResponseTopicName     string `validate:"required" env:"TOPIC_NAME_TG_RESPONSE_PREPARER"`
	SaslUsername          string `env:"SASL_USERNAME"`
	SaslPassword          string `env:"SASL_PASSWORD"`
	ClientID              string `env:"CLIENT_ID_MESSAGE_RESPONDER"`
	OcrTopicName          string `validate:"required" env:"TOPIC_NAME_OCR_REQUEST"`
}
type Config struct {
	Kafka Kafka `envPrefix:"KAFKA_"`
}

func Load() (*Config, error) {
	var c Config
	if err := env.Parse(&c); err != nil {
		return nil, fmt.Errorf("env parse: %w", err)
	}
	v := validator.New()
	if err := v.Struct(c); err != nil {
		return nil, fmt.Errorf("config validate: %w", err)
	}
	return &c, nil
}
