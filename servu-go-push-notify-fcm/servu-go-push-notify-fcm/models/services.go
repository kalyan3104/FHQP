package models

import (
	"time"

	"github.com/lucsky/cuid"
	"gorm.io/gorm"
)

type Service struct {
	Id          string         `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"unique;not null"`
	License     string         `json:"license"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func NewService(title, license, description string) *Service {
	return &Service{
		Id:          cuid.New(),
		Title:       title,
		License:     license,
		Description: description,
	}
}