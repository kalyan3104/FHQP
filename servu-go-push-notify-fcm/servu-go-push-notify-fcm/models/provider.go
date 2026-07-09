package models

import (
	"time"

	"gorm.io/gorm"
)

type Provider struct {
	gorm.Model
	ID            string `gorm:"primaryKey"`
	Name          string
	Email         string
	LicenseNumber string
	Phone         string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateProviderRequest struct {
	Name          string `json:"name"`
	Email         string `json:"email"`
	LicenseNumber string `json:"license_number"`
	Phone         string `json:"phone"`
}
