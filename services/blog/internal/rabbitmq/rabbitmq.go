package rabbitmq

import (
	"fmt"
	"net/url"

	amqp "github.com/rabbitmq/amqp091-go"

	"blog/internal/logger"
)

var Conn *amqp.Connection
var Channel *amqp.Channel

// ConnectRabbitMQ opens a connection and channel to RabbitMQ
func ConnectRabbitMQ(host, username, password string) error {
	uri := fmt.Sprintf("amqp://%s:%s@%s:5672/",
		url.QueryEscape(username),
		url.QueryEscape(password),
		host,
	)

	var err error
	Conn, err = amqp.Dial(uri)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	Channel, err = Conn.Channel()
	if err != nil {
		Conn.Close()
		return fmt.Errorf("failed to open a channel: %w", err)
	}

	logger.Logger.Info("Connected to RabbitMQ successfully")
	return nil
}

// CloseRabbitMQ cleans up RabbitMQ resources safely
func CloseRabbitMQ() {
	if Channel != nil {
		_ = Channel.Close()
	}
	if Conn != nil {
		_ = Conn.Close()
	}
}

// CacheInvalidationMessage represents the invalidation event structure
type CacheInvalidationMessage struct {
	Action string   `json:"action"`
	Keys   []string `json:"keys"`
}
