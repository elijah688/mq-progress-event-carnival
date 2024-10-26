package config

import (
	"fmt"
	"os"
	"strings"
)

type QueueConfig struct {
	AMQPURL    string
	QueueName  string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
}

func parseBool(value string) (bool, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "true", "1":
		return true, nil
	case "false", "0":
		return false, nil
	default:
		return false, fmt.Errorf("invalid boolean value: %s", value)
	}
}

func NewQueueConfig() (*QueueConfig, error) {
	amqpURL := os.Getenv("RABBITMQ_URL")
	queueName := os.Getenv("RABBITMQ_QUEUE_NAME")
	durable, err := parseBool(os.Getenv("RABBITMQ_DURABLE"))
	if err != nil {
		return nil, err
	}
	autoDelete, err := parseBool(os.Getenv("RABBITMQ_AUTO_DELETE"))
	if err != nil {
		return nil, err
	}
	exclusive, err := parseBool(os.Getenv("RABBITMQ_EXCLUSIVE"))
	if err != nil {
		return nil, err
	}
	noWait, err := parseBool(os.Getenv("RABBITMQ_NO_WAIT"))
	if err != nil {
		return nil, err
	}

	// Validate required fields
	if amqpURL == "" {
		return nil, fmt.Errorf("RABBITMQ_URL environment variable is required")
	}
	if queueName == "" {
		return nil, fmt.Errorf("RABBITMQ_QUEUE_NAME environment variable is required")
	}

	return &QueueConfig{
		AMQPURL:    amqpURL,
		QueueName:  queueName,
		Durable:    durable,
		AutoDelete: autoDelete,
		Exclusive:  exclusive,
		NoWait:     noWait,
	}, nil
}
