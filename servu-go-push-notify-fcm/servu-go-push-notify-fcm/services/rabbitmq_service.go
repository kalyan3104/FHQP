package services

import (
	"context"
	"encoding/json"
	"time"

	"github.com/baryogenesis2025/servu-go/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	// Booking API publishes events to this exchange.
	// Notification worker listens through the queue below.
	notificationExchange  = "servu.notifications"
	bookingCreatedQueue   = "servu.booking.notifications"
	bookingCreatedRouting = "booking.created"
)

type RabbitMQService struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

func NewRabbitMQService(url string) (*RabbitMQService, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	service := &RabbitMQService{
		Conn:    conn,
		Channel: ch,
	}

	if err := service.setupBookingCreatedQueue(); err != nil {
		service.Close()
		return nil, err
	}

	return service, nil
}

func (s *RabbitMQService) setupBookingCreatedQueue() error {
	// Durable exchange and queue survive RabbitMQ restarts.
	if err := s.Channel.ExchangeDeclare(
		notificationExchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	if _, err := s.Channel.QueueDeclare(
		bookingCreatedQueue,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	return s.Channel.QueueBind(
		bookingCreatedQueue,
		bookingCreatedRouting,
		notificationExchange,
		false,
		nil,
	)
}

func (s *RabbitMQService) PublishBookingCreated(event models.BookingCreatedEvent) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.Channel.PublishWithContext(
		ctx,
		notificationExchange,
		bookingCreatedRouting,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			// Persistent messages stay in RabbitMQ if the broker restarts before delivery.
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
}

func (s *RabbitMQService) ConsumeBookingCreated(handler func(models.BookingCreatedEvent) error) error {
	// Process one message at a time so one worker does not take many unprocessed jobs.
	if err := s.Channel.Qos(1, 0, false); err != nil {
		return err
	}

	messages, err := s.Channel.Consume(
		bookingCreatedQueue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	for message := range messages {
		var event models.BookingCreatedEvent
		if err := json.Unmarshal(message.Body, &event); err != nil {
			// Bad JSON cannot be processed, so reject it without requeueing forever.
			message.Nack(false, false)
			continue
		}

		if err := handler(event); err != nil {
			// If processing fails, requeue so RabbitMQ can retry later.
			message.Nack(false, true)
			continue
		}

		// Ack only after notification persistence and delivery attempts finish.
		message.Ack(false)
	}

	return nil
}

func (s *RabbitMQService) Close() {
	if s.Channel != nil {
		s.Channel.Close()
	}

	if s.Conn != nil {
		s.Conn.Close()
	}
}
