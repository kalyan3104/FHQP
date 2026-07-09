package models

import (
	"time"

	"github.com/lucsky/cuid"
)

type UserService struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	UserID    string    `json:"user_id" gorm:"not null;index"`    // Foreign key to User
	ServiceID string    `json:"service_id" gorm:"not null;index"` // Foreign key to Service
	CreatedAt time.Time `json:"created_at"`

	// GORM relations
	User    User    `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Service Service `gorm:"foreignKey:ServiceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func NewUserService(userID, serviceID string) *UserService {
	return &UserService{
		ID:        cuid.New(),
		UserID:    userID,
		ServiceID: serviceID,
		CreatedAt: time.Now(),
	}
}
