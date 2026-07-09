package services

import (
	"log"

	"github.com/baryogenesis2025/servu-go/models"
	"github.com/baryogenesis2025/servu-go/websocket"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type NotificationWorker struct {
	RabbitMQ                *RabbitMQService
	NotificationService     *NotificationService
	PushNotificationService *PushNotificationService
}

func NewNotificationWorker(db *gorm.DB, rabbitMQ *RabbitMQService) *NotificationWorker {
	return &NotificationWorker{
		RabbitMQ:                rabbitMQ,
		NotificationService:     NewNotificationService(db),
		PushNotificationService: NewPushNotificationService(db),
	}
}

func (w *NotificationWorker) Start() error {
	return w.RabbitMQ.ConsumeBookingCreated(w.handleBookingCreated)
}

func (w *NotificationWorker) handleBookingCreated(event models.BookingCreatedEvent) error {
	// Step 1: save notification in DB, so the provider can see it even if they are offline.
	if err := w.NotificationService.CreateNotification(
		event.ProviderID,
		event.Title,
		event.Message,
	); err != nil {
		return err
	}

	// Step 2: send Firebase push notification if the provider has an FCM token.
	token, err := w.PushNotificationService.GetUserFCMToken(event.ProviderID)
	if err == nil {
		if err := w.PushNotificationService.SendPushNotification(
			token,
			event.Title,
			event.Message,
		); err != nil {
			log.Println("Error sending push notification:", err)
		}
	} else {
		log.Println("FCM token not found:", err)
	}

	// Step 3: send live websocket notification only if the provider is connected to this worker.
	if err := websocket.SendToUser(event.ProviderID, fiber.Map{
		"type":       "NEW_BOOKING",
		"booking_id": event.BookingID,
		"title":      event.Title,
		"message":    event.Message,
	}); err != nil {
		log.Println("Websocket client not available:", err)
	}

	return nil
}
