package service

import (
	"encoding/json"
	"fmt"
	"time"

	"unsri-backend/internal/api-gateway/handler"
	"unsri-backend/internal/shared/logger"
	"unsri-backend/internal/shared/messaging"

	amqp "github.com/rabbitmq/amqp091-go"
)

// MessageBrokerService handles message broker operations
type MessageBrokerService struct {
	client *messaging.RabbitMQClient
	logger logger.Logger
}

// NewMessageBrokerService creates a new message broker service
func NewMessageBrokerService(client *messaging.RabbitMQClient, logger logger.Logger) *MessageBrokerService {
	return &MessageBrokerService{
		client: client,
		logger: logger,
	}
}

// Initialize sets up exchanges and queues
func (s *MessageBrokerService) Initialize() error {
	// Declare exchanges
	exchanges := []struct {
		name       string
		kind       string
		durable    bool
		autoDelete bool
	}{
		{"api_gateway_events", "topic", true, false},
		{"service_events", "topic", true, false},
		{"notifications", "topic", true, false},
		{"audit_logs", "direct", true, false},
	}

	for _, ex := range exchanges {
		if err := s.client.DeclareExchange(ex.name, ex.kind, ex.durable, ex.autoDelete, false, false); err != nil {
			return fmt.Errorf("failed to declare exchange %s: %w", ex.name, err)
		}
		s.logger.Infof("Exchange '%s' declared", ex.name)
	}

	// Declare queues
	queues := []struct {
		name       string
		durable    bool
		autoDelete bool
		exclusive  bool
	}{
		{"api_gateway_audit_logs", true, false, false},
		{"api_gateway_request_logs", true, false, false},
		{"notification_queue", true, false, false},
	}

	for _, q := range queues {
		_, err := s.client.DeclareQueue(q.name, q.durable, q.autoDelete, q.exclusive, false)
		if err != nil {
			return fmt.Errorf("failed to declare queue %s: %w", q.name, err)
		}
		s.logger.Infof("Queue '%s' declared", q.name)
	}

	// Bind queues to exchanges
	bindings := []struct {
		queue      string
		key        string
		exchange   string
	}{
		{"api_gateway_audit_logs", "audit.#", "audit_logs"},
		{"api_gateway_request_logs", "request.#", "api_gateway_events"},
		{"notification_queue", "notification.#", "notifications"},
	}

	for _, b := range bindings {
		if err := s.client.QueueBind(b.queue, b.key, b.exchange, false); err != nil {
			return fmt.Errorf("failed to bind queue %s: %w", b.queue, err)
		}
		s.logger.Infof("Queue '%s' bound to exchange '%s' with key '%s'", b.queue, b.exchange, b.key)
	}

	return nil
}

// Note: RequestLog and AuditLog types are defined in handler package
// This service implements the MessageBrokerService interface

// PublishRequestLog publishes a request log to message broker
func (s *MessageBrokerService) PublishRequestLog(log *handler.RequestLog) error {
	body, err := json.Marshal(log)
	if err != nil {
		return fmt.Errorf("failed to marshal request log: %w", err)
	}

	routingKey := fmt.Sprintf("request.%s.%s", log.Service, log.Method)
	if err := s.client.Publish(
		"api_gateway_events",
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			Body:         body,
		},
	); err != nil {
		return fmt.Errorf("failed to publish request log: %w", err)
	}

	s.logger.Debugf("Published request log: %s %s", log.Method, log.Path)
	return nil
}

// PublishAuditLog publishes an audit log to message broker
func (s *MessageBrokerService) PublishAuditLog(log *handler.AuditLog) error {
	body, err := json.Marshal(log)
	if err != nil {
		return fmt.Errorf("failed to marshal audit log: %w", err)
	}

	routingKey := fmt.Sprintf("audit.%s.%s", log.Action, log.Resource)
	if err := s.client.Publish(
		"audit_logs",
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			Body:         body,
		},
	); err != nil {
		return fmt.Errorf("failed to publish audit log: %w", err)
	}

	s.logger.Debugf("Published audit log: %s %s", log.Action, log.Resource)
	return nil
}

// PublishServiceEvent publishes a service event to message broker
func (s *MessageBrokerService) PublishServiceEvent(eventType, service, routingKey string, payload interface{}) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal event payload: %w", err)
	}

	fullRoutingKey := fmt.Sprintf("%s.%s.%s", eventType, service, routingKey)
	if err := s.client.Publish(
		"service_events",
		fullRoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			Body:         body,
			Headers: amqp.Table{
				"event_type": eventType,
				"service":    service,
			},
		},
	); err != nil {
		return fmt.Errorf("failed to publish service event: %w", err)
	}

	s.logger.Debugf("Published service event: %s", fullRoutingKey)
	return nil
}

// StartConsumer starts consuming messages from a queue
func (s *MessageBrokerService) StartConsumer(queue, consumerTag string, handler func(amqp.Delivery) error) error {
	msgs, err := s.client.Consume(queue, consumerTag, false, false, false, false)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	go func() {
		for msg := range msgs {
			if err := handler(msg); err != nil {
				s.logger.Errorf("Error processing message: %v", err)
				msg.Nack(false, true) // Requeue on error
			} else {
				msg.Ack(false)
			}
		}
	}()

	s.logger.Infof("Started consumer '%s' on queue '%s'", consumerTag, queue)
	return nil
}

// Close closes the message broker connection
func (s *MessageBrokerService) Close() error {
	return s.client.Close()
}

