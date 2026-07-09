package models

import "gorm.io/gorm"

type DeviceToken struct {
	gorm.Model
	ID     string `gorm:"primaryKey"`
	UserID string
	Token  string
}
