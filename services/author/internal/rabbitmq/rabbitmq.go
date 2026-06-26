package rabbitmq
import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"

	"author/internal/logger"
)

var mqConn *amqp.Connection
var mqChannel *amqp.Channel

func ConnectRabbitMQ(host, username, password string) error {
	protocol := "amqp"
	if strings.Contains(host, ".amazonaws.com") {
		protocol = "amqps"
	}

	hasPort := strings.Contains(host, ":")
	
	var uri string
	if hasPort {
		uri = fmt.Sprintf("%s://%s:%s@%s/",
			protocol,
			url.QueryEscape(username),
			url.QueryEscape(password),
			host,
		)
	} else {
		port := "5672"
		if protocol == "amqps" {
			port = "5671"
		}
		uri = fmt.Sprintf("%s://%s:%s@%s:%s/",
			protocol,
			url.QueryEscape(username),
			url.QueryEscape(password),
			host,
			port,
		)
	}

	var err error
	mqConn, err = amqp.Dial(uri)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	mqChannel, err = mqConn.Channel()
	if err != nil {
		mqConn.Close()
		return fmt.Errorf("failed to open a channel: %w", err)
	}

	logger.Logger.Info("Connected to RabbitMQ successfully")
	return nil
}

func CloseRabbitMQ() {
	if mqChannel != nil {
		_ = mqChannel.Close()
	}
	if mqConn != nil {
		_ = mqConn.Close()
	}
}

type CacheInvalidationMessage struct {
	Action string   `json:"action"`
	Keys   []string `json:"keys"`
}

func PublishToQueue(queueName string, message interface{}) error {
	if mqChannel == nil {
		return fmt.Errorf("RabbitMQ channel is not established")
	}

	// Assert the queue exists and is durable (to match JS durable: true)
	_, err := mqChannel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = mqChannel.PublishWithContext(ctx,
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // persistent: true
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	logger.Logger.Info("Message sent to RabbitMQ queue", zap.String("queue", queueName))
	return nil
}

func InvalidateCacheJob(keys []string) {
	msg := CacheInvalidationMessage{
		Action: "invalidateCache",
		Keys:   keys,
	}
	err := PublishToQueue("cache-invalidation", msg)
	if err != nil {
		logger.Logger.Error("Failed to publish cache invalidation job", zap.Error(err))
	} else {
		logger.Logger.Info("Cache invalidation job published to RabbitMQ successfully")
	}
}
