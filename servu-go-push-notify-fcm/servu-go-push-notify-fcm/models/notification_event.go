package models

import "time"

type BookingCreatedEvent struct {
	EventID    string    `json:"event_id"`
	EventType  string    `json:"event_type"`
	BookingID  string    `json:"booking_id"`
	CustomerID string    `json:"customer_id"`
	ProviderID string    `json:"provider_id"`
	Title      string    `json:"title"`
	Message    string    `json:"message"`
	CreatedAt  time.Time `json:"created_at"`
}
