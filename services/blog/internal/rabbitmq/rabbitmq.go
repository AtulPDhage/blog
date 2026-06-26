package rabbitmq

import (
	"fmt"
	"net/url"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"

	"blog/internal/logger"
)

var Conn *amqp.Connection
var Channel *amqp.Channel

// ConnectRabbitMQ opens a connection and channel to RabbitMQ
func ConnectRabbitMQ(host, username, password string) error {
	protocol := "amqp"
	if strings.HasPrefix(host, "amqps://") ||
		strings.Contains(host, ".amazonaws.com") ||
		strings.Contains(host, ".on.aws") ||
		strings.HasSuffix(host, ":5671") {
		protocol = "amqps"
	}

	host = strings.TrimPrefix(host, "amqps://")
	host = strings.TrimPrefix(host, "amqp://")

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
