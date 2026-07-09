package handlers

import (
	"log"
	"time"

	"github.com/baryogenesis2025/servu-go/models"
	"github.com/baryogenesis2025/servu-go/services"
	"github.com/gofiber/fiber/v2"
	"github.com/lucsky/cuid"
)

type BookingHandler struct {
	Service *services.BookingService
}

func NewBookingHandler(s *services.BookingService) *BookingHandler {
	return &BookingHandler{Service: s}
}

func parseTime(value string) time.Time {
	t, _ := time.Parse(time.RFC3339, value)
	return t
}

// CreateBooking godoc
// @Summary Create a booking
// @Description Create a new service booking for a customer
// @Tags bookings
// @Accept json
// @Produce json
// @Param body body models.CreateBookingRequest true "Booking payload"
// @Success 201 {object} models.Booking
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /bookings/ [post]
func CreateBooking(service *services.BookingService, rabbitMQService *services.RabbitMQService) fiber.Handler {
	log.Println("Initializing CreateBooking handler")
	return func(c *fiber.Ctx) error {
		var req models.CreateBookingRequest

		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
		}

		slotTime, err := time.Parse(time.RFC3339, req.SlotTime)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid time format"})
		}

		booking, err := service.CreateBooking(req.CustomerID, slotTime, req.ProviderID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		// Booking API does not send notifications directly anymore.
		// It publishes one event to RabbitMQ, then the notification worker handles DB notification, FCM, and websocket.
		event := models.BookingCreatedEvent{
			EventID:    cuid.New(),
			EventType:  "BOOKING_CREATED",
			BookingID:  booking.ID,
			CustomerID: req.CustomerID,
			ProviderID: req.ProviderID,
			Title:      "New Booking",
			Message:    "You received a booking request",
			CreatedAt:  time.Now(),
		}

		if err := rabbitMQService.PublishBookingCreated(event); err != nil {
			log.Println("Error publishing booking notification event:", err)
			return c.Status(500).JSON(fiber.Map{"error": "booking created but notification event failed"})
		}

		return c.Status(fiber.StatusCreated).JSON(booking)
	}
}

// AcceptBooking godoc
// @Summary Accept a booking
// @Description Accept a booking request
// @Tags bookings
// @Accept json
// @Produce json
// @Param id path string true "Booking ID"
// @Param body body models.Booking true "Accept booking payload"
// @Success 200 {object} models.Booking
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /bookings/accept [post]
func AcceptBooking(service *services.BookingService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req struct {
			BookingID  string `json:"booking_id"`
			ProviderID string `json:"provider_id"`
			Decision   string `json:"decision"`
		}

		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
		}

		log.Println("Accept request:", req)

		err := service.AcceptBooking(req.BookingID, req.ProviderID, req.Decision)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"status": "success", "message": "Booking " + req.Decision})
	}
}

func GetBookingByID(service *services.BookingService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		booking_id := c.Params("id")

		booking, err := service.GetBookingByID(booking_id)

		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Booking not found",
			})
		}

		return c.JSON(booking)
	}
}
