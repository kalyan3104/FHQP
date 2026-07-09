package services

import (
	"errors"
	"time"

	"github.com/baryogenesis2025/servu-go/models"
	"github.com/lucsky/cuid"
	"gorm.io/gorm"
)

type BookingService struct {
	DB *gorm.DB
}

func NewBookingService(db *gorm.DB) *BookingService {
	return &BookingService{DB: db}
}

func (s *BookingService) CreateBooking(customerID string, slot time.Time, providerID string) (*models.Booking, error) {

	booking := models.Booking{
		ID:            cuid.New(), // your cuid/uuid
		CustomerID:    customerID,
		ScheduledTime: slot,
		AssignedToID:  providerID,
		Status:        "PENDING",
		CreatedAt:     time.Now(),
	}

	if err := s.DB.Create(&booking).Error; err != nil {
		return nil, err
	}

	return &booking, nil
}

func (s *BookingService) AcceptBooking(bookingID, providerID string, decision string) error {


	var status string

	switch decision {
	case "ACCEPT":
		status = "CONFIRMED"

	case "REJECT":
		status = "REJECTED"

	default:
		return errors.New("invalid decision")
	}

	result := s.DB.Model(&models.Booking{}).
		Where("id = ? AND status = ?", bookingID, "PENDING").
		Update("status", status)

		if result.RowsAffected == 0 {
			return errors.New("already taken")
		}

	return nil
}

func (s * BookingService) GetBookingByID(bookingID string) (*models.Booking, error) {
	var booking models.Booking
	if err := s.DB.First(&booking, "id = ?", bookingID).Error; err != nil {
		return nil, err
	}
	return &booking, nil
}
