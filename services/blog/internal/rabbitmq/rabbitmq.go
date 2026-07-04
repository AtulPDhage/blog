package rabbitmq

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"

	"blog/internal/logger"
)

var Conn *amqp.Connection
var Channel *amqp.Channel

// ConnectRabbitMQ opens a connection and channel to RabbitMQ
func ConnectRabbitMQ(host, username, password string) error {
	var uri string

	// If host is already a full URI, use it directly
	if strings.Contains(host, "://") {
		uri = host
	} else {
		protocol := "amqp"
		if strings.HasPrefix(host, "amqps://") ||
			strings.Contains(host, ".amazonaws.com") ||
			strings.Contains(host, ".on.aws") ||
			strings.Contains(host, "cloudamqp.com") ||
			strings.Contains(host, "lavemq.com") ||
			strings.HasSuffix(host, ":5671") {
			protocol = "amqps"
		}

		host = strings.TrimPrefix(host, "amqps://")
		host = strings.TrimPrefix(host, "amqp://")

		hasPort := strings.Contains(host, ":")
		
		// Detect vhost (CloudAMQP uses the username as the vhost)
		vhost := ""
		if strings.Contains(host, "cloudamqp.com") || strings.Contains(host, "lavemq.com") {
			vhost = username
		}
		
		// Support manual override
		if customVhost := os.Getenv("Rabbitmq_Vhost"); customVhost != "" {
			vhost = customVhost
		}

		if hasPort {
			uri = fmt.Sprintf("%s://%s:%s@%s/%s",
				protocol,
				url.QueryEscape(username),
				url.QueryEscape(password),
				host,
				url.PathEscape(vhost),
			)
		} else {
			port := "5672"
			if protocol == "amqps" {
				port = "5671"
			}
			uri = fmt.Sprintf("%s://%s:%s@%s:%s/%s",
				protocol,
				url.QueryEscape(username),
				url.QueryEscape(password),
				host,
				port,
				url.PathEscape(vhost),
			)
		}
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
