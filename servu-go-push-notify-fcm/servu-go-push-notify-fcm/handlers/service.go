package handlers

import (
	"github.com/baryogenesis2025/servu-go/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)


func GetAllServices(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var services []models.Service
		if err := db.Find(&services).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(services)
	}
}