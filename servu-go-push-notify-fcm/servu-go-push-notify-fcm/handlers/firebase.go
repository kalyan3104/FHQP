package handlers

import (
	"github.com/baryogenesis2025/servu-go/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SaveFCMToken(db *gorm.DB) fiber.Handler {

	return func(c *fiber.Ctx) error {

		var req models.DeviceToken

		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "invalid body",
			})
		}

		err := db.Create(&req).Error

		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"message": "token saved",
		})
	}
}
