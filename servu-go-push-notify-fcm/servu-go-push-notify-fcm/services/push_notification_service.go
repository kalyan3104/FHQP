package services

import (
	"context"

	"firebase.google.com/go/messaging"
	"gorm.io/gorm"

	"github.com/baryogenesis2025/servu-go/firebase"
	"github.com/baryogenesis2025/servu-go/models"
)

type PushNotificationService struct {
	DB *gorm.DB
}

func NewPushNotificationService(db *gorm.DB) *PushNotificationService {
	return &PushNotificationService{DB: db}
}

func (s *PushNotificationService) SendPushNotification(
	token string,
	title string,
	body string,
) error {

	client, err := firebase.App.Messaging(
		context.Background(),
	)

	if err != nil {
		return err
	}

	message := &messaging.Message{
		Token: token,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
	}

	_, err = client.Send(
		context.Background(),
		message,
	)

	return err
}

func (s *PushNotificationService) GetUserFCMToken(
	userID string,
) (string, error) {

	var token models.DeviceToken

	err := s.DB.
		Where("user_id = ?", userID).
		First(&token).Error

	if err != nil {
		return "", err
	}

	return token.Token, nil
}
