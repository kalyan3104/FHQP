package handlers

import (
	"github.com/baryogenesis2025/servu-go/models"
	"github.com/baryogenesis2025/servu-go/services"
	"github.com/gofiber/fiber/v2"
)

func GetNotifications(service *services.NotificationService) fiber.Handler {

	return func(c *fiber.Ctx) error {

		userID := c.Params("user_id")

		var notifications []models.Notification

		err := service.DB.
			Where("user_id = ?", userID).
			Order("created_at desc").
			Find(&notifications).Error

		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(notifications)
	}
}
