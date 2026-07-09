package handlers

import (
	"github.com/baryogenesis2025/servu-go/models"
	"github.com/baryogenesis2025/servu-go/services"
	"github.com/gofiber/fiber/v2"
)

type ProviderHandler struct {
	Service *services.ProviderService
}

func NewProviderHandler(s *services.ProviderService) *ProviderHandler {
	return &ProviderHandler{Service: s}
}

func CreateProvider(service *services.ProviderService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req models.CreateProviderRequest

		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
		}

		provider, err := service.CreateProvider(req.Name, req.Email, req.LicenseNumber, req.Phone)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(201).JSON(provider)
	}
}

func GetProviderByID(service *services.ProviderService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		provider, err := service.GetProviderByID(id)

		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "provider not found"})
		}
		return c.JSON(provider)
	}
}
