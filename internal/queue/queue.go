package queue

import (
	"fmt"
	"log"
	"messages/internal/config"

	"github.com/streadway/amqp"
)

type Queue struct {
	config *config.QueueConfig
	conn   *amqp.Connection
	ch     *amqp.Channel
}

func NewQueue(config *config.QueueConfig) (*Queue, error) {
	conn, err := amqp.Dial(config.AMQPURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %v", err)
	}

	if _, err = ch.QueueDeclare(
		config.QueueName,
		config.Durable,
		config.AutoDelete,
		config.Exclusive,
		config.NoWait,
		nil,
	); err != nil {
		return nil, fmt.Errorf("failed to declare a queue: %v", err)
	}

	return &Queue{
		config: config,
		conn:   conn,
		ch:     ch,
	}, nil
}

// Close closes the RabbitMQ channel.
func (q *Queue) Close() {
	if err := q.ch.Close(); err != nil {
		log.Printf("Failed to close channel: %v", err)
	}
}

// Publish sends a message to the RabbitMQ queue.
func (q *Queue) Publish(body []byte) error {
	err := q.ch.Publish(
		"",
		q.config.QueueName,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish a message: %v", err)
	}
	return nil
}

// Consume starts consuming messages from the RabbitMQ queue.
func (q *Queue) Consume() (<-chan amqp.Delivery, error) {
	msgs, err := q.ch.Consume(
		q.config.QueueName,
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register a consumer: %v", err)
	}
	return msgs, nil
}

// func main() {
// 	// Example usage of the Queue
// 	config, err := NewQueueConfig()
// 	if err != nil {
// 		log.Fatalf("Failed to load queue config: %v", err)
// 	}

// 	queue, err := NewQueue(config)
// 	if err != nil {
// 		log.Fatalf("Failed to create queue: %v", err)
// 	}
// 	defer queue.Close() // Ensure the queue is closed when done

// 	// Publish a sample message
// 	msg := []byte(`{"id":"1","name":"Test Message","user":"user@example.com","state":"Running","startTime":"2024-10-26T15:24:40Z","finishedTime":"0001-01-01T00:00:00Z","duration":"10 secs","errorMessage":"","percentageComplete":0.5}`)
// 	if err := queue.Publish(msg); err != nil {
// 		log.Printf("Error publishing message: %v", err)
// 	}

// 	// Start consuming messages
// 	msgs, err := queue.Consume()
// 	if err != nil {
// 		log.Fatalf("Failed to consume messages: %v", err)
// 	}

// 	// Process incoming messages
// 	go func() {
// 		for msg := range msgs {
// 			fmt.Printf("Received message: %s\n", msg.Body)
// 		}
// 	}()

// 	// Keep the main goroutine running
// 	select {}
// }
