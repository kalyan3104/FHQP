package models

import (
	"time"

	"gorm.io/gorm"
)

type Booking struct {
	gorm.Model
	ID            string `gorm:"primaryKey"`
	CustomerID    string
	ScheduledTime time.Time
	Status        string
	AssignedToID  string
	CreatedAt     time.Time
}

type CreateBookingRequest struct {
	CustomerID string `json:"customer_id" example:"cust1"`
	ProviderID string `json:"provider_id" example:"prov1"`
	SlotTime   string `json:"slot_time" example:"2026-05-05T14:30:00Z"`
}