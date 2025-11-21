package messaging

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQClient wraps RabbitMQ connection and channel
type RabbitMQClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// Config holds RabbitMQ configuration
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	VHost    string
}

// NewRabbitMQ creates a new RabbitMQ client
func NewRabbitMQ(cfg Config) (*RabbitMQClient, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.VHost,
	)

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	return &RabbitMQClient{
		conn:    conn,
		channel: channel,
	}, nil
}

// DeclareQueue declares a queue
func (r *RabbitMQClient) DeclareQueue(name string, durable, autoDelete, exclusive, noWait bool) (amqp.Queue, error) {
	return r.channel.QueueDeclare(
		name,      // name
		durable,   // durable
		autoDelete, // auto-delete
		exclusive, // exclusive
		noWait,    // no-wait
		nil,       // arguments
	)
}

// Publish publishes a message to a queue
func (r *RabbitMQClient) Publish(exchange, routingKey string, mandatory, immediate bool, msg amqp.Publishing) error {
	return r.channel.Publish(
		exchange,   // exchange
		routingKey, // routing key
		mandatory,  // mandatory
		immediate,  // immediate
		msg,
	)
}

// Consume consumes messages from a queue
func (r *RabbitMQClient) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool) (<-chan amqp.Delivery, error) {
	return r.channel.Consume(
		queue,     // queue
		consumer,  // consumer
		autoAck,   // auto-ack
		exclusive, // exclusive
		noLocal,   // no-local
		noWait,    // no-wait
		nil,       // args
	)
}

// DeclareExchange declares an exchange
func (r *RabbitMQClient) DeclareExchange(name, kind string, durable, autoDelete, internal, noWait bool) error {
	return r.channel.ExchangeDeclare(
		name,       // name
		kind,       // kind
		durable,    // durable
		autoDelete, // auto-delete
		internal,   // internal
		noWait,     // no-wait
		nil,        // arguments
	)
}

// QueueBind binds a queue to an exchange
func (r *RabbitMQClient) QueueBind(queue, key, exchange string, noWait bool) error {
	return r.channel.QueueBind(
		queue,    // queue name
		key,      // routing key
		exchange, // exchange
		noWait,   // no-wait
		nil,      // arguments
	)
}

// Close closes the connection
func (r *RabbitMQClient) Close() error {
	if r.channel != nil {
		if err := r.channel.Close(); err != nil {
			log.Printf("Error closing channel: %v", err)
		}
	}
	if r.conn != nil {
		if err := r.conn.Close(); err != nil {
			return fmt.Errorf("error closing connection: %w", err)
		}
	}
	return nil
}

// Channel returns the underlying channel
func (r *RabbitMQClient) Channel() *amqp.Channel {
	return r.channel
}

// Connection returns the underlying connection
func (r *RabbitMQClient) Connection() *amqp.Connection {
	return r.conn
}

