package routes

import (
	"github.com/baryogenesis2025/servu-go/auth"
	"github.com/baryogenesis2025/servu-go/handlers"
	"github.com/baryogenesis2025/servu-go/services"
	"github.com/gofiber/fiber/v2"

	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, authService *auth.AuthService, db *gorm.DB, rabbitMQService *services.RabbitMQService) {
	otpHandler := handlers.NewOTPHandler(authService)
	authGroup := app.Group("/auth")
	authGroup.Post("/send-verification", otpHandler.SendOTP)
	authGroup.Post("/check-verification", otpHandler.VerifyOTP)
	users := app.Group("/users")
	users.Get("/:phonenumber", handlers.HandleUser(db))
	users.Get("/service", handlers.GetAllServices(db))

	//addresses
	addresses := app.Group("/address")
	addresses.Post("/", handlers.CreateAddress(db))
	addresses.Get("/user/:userId", handlers.GetUserAddresses(db))
	addresses.Get("/:id", handlers.GetAddressByID(db))
	addresses.Put("/:id", handlers.UpdateAddress(db))
	addresses.Delete("/:id", handlers.DeleteAddress(db))

	//notifications
	notificationService := services.NewNotificationService(db)
	notifications := app.Group("/notifications")
	notifications.Get("/:user_id", handlers.GetNotifications(notificationService))

	//bookings
	bookingService := services.NewBookingService(db)
	bookings := app.Group("/bookings")
	bookings.Post("/", handlers.CreateBooking(bookingService, rabbitMQService))
	bookings.Post("/accept", handlers.AcceptBooking(bookingService))
	bookings.Get("/:id", handlers.GetBookingByID(bookingService))

	//providers
	providerService := services.NewProviderService(db)
	providers := app.Group("/providers")
	providers.Post("/", handlers.CreateProvider(providerService))
	providers.Get("/:id", handlers.GetProviderByID(providerService))

	//services
	servicesGroup := app.Group("/services")
	servicesGroup.Get("/", handlers.GetAllServices(db))

}
