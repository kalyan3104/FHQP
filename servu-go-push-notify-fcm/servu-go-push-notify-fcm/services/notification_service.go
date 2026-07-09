package services

import (
	"github.com/baryogenesis2025/servu-go/models"
	"github.com/lucsky/cuid"
	"gorm.io/gorm"
)

type NotificationService struct {
	DB *gorm.DB
}

func NewNotificationService(db *gorm.DB) *NotificationService {
	return &NotificationService{
		DB: db,
	}
}

func (s *NotificationService) CreateNotification(
	userID string,
	title string,
	message string,
) error {

	notification := models.Notification{
		ID:      cuid.New(),
		UserID:  userID,
		Title:   title,
		Message: message,
		IsRead:  false,
	}

	return s.DB.Create(&notification).Error
}