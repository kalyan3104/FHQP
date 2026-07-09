package services

import (
	"time"

	"github.com/baryogenesis2025/servu-go/models"
	"github.com/lucsky/cuid"
	"gorm.io/gorm"
)

type ProviderService struct {
	DB *gorm.DB
}

func NewProviderService(db *gorm.DB) *ProviderService {
	return &ProviderService{DB: db}
}

func (s *ProviderService) CreateProvider(name, email, licenseNumber, phone string) (*models.Provider, error) {

	provider := models.Provider{
		ID:            cuid.New(), // your cuid/uuid
		Name:          name,
		Email:         email,
		LicenseNumber: licenseNumber,
		Phone:         phone,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.DB.Create(&provider).Error; err != nil {
		return nil, err
	}
	return &provider, nil
}

func (s *ProviderService) GetProviderByID(id string) (*models.Provider, error) {
	var provider models.Provider
	if err := s.DB.First(&provider, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &provider, nil
}
