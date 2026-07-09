package main

import (
	"log"
	"os"

	"github.com/baryogenesis2025/servu-go/auth"
	"github.com/baryogenesis2025/servu-go/db"
	"github.com/baryogenesis2025/servu-go/firebase"
	"github.com/baryogenesis2025/servu-go/routes"
	"github.com/baryogenesis2025/servu-go/services"
	myws "github.com/baryogenesis2025/servu-go/websocket"

	// "github.com/baryogenesis2025/servu-go/handlers"
	_ "github.com/baryogenesis2025/servu-go/docs"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	fiberws "github.com/gofiber/websocket/v2"
	"gorm.io/gorm"
)

// @title Fiber API
// @version 1.0
// @description This is a sample Fiber server.
// @host localhost:8000
// @BasePath /
func main() {
	database := db.ConnectDB()

	// Auto-migrate auth models
	database.AutoMigrate(&auth.User{}, &auth.RefreshToken{})

	//Initialize Firebase
	firebase.InitFirebase()

	// Same binary can run as API or as notification worker.
	// API mode: creates bookings and publishes RabbitMQ events.
	// Worker mode: consumes RabbitMQ events and owns /ws.
	if os.Getenv("APP_MODE") == "notification-worker" {
		startNotificationWorker(database)
		return
	}

	startAPI(database)
}

func startAPI(database *gorm.DB) {
	// Initialize Twilio service
	twilioService := auth.NewTwilioService()

	// Initialize Auth service
	authService := auth.NewAuthService(
		database,
		twilioService,
		os.Getenv("JWT_SECRET"),
	)

	rabbitMQURL := getEnv("RABBITMQ_URL", "amqp://guest:guest@127.0.0.1:5672/")
	log.Printf("Connecting to RabbitMQ at %s", rabbitMQURL)
	rabbitMQService, err := services.NewRabbitMQService(rabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitMQService.Close()

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowMethods: "GET,POST,PUT,DELETE",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))
	app.Get("/swagger/*", swagger.HandlerDefault)
	routes.SetupRoutes(app, authService, database, rabbitMQService)
	err = app.Listen(":" + getEnv("PORT", "8001"))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	} else {
		log.Println("Server started")
	}
}

func startNotificationWorker(database *gorm.DB) {
	rabbitMQService, err := services.NewRabbitMQService(getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"))
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitMQService.Close()

	worker := services.NewNotificationWorker(database, rabbitMQService)
	go func() {
		if err := worker.Start(); err != nil {
			log.Fatalf("Notification worker stopped: %v", err)
		}
	}()

	// Websocket lives on the notification worker, not the booking API.
	// Local URL: ws://localhost:8002/ws?user_id=<provider_id>
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))
	app.Get("/ws", fiberws.New(myws.WebSocketHandler))

	if err := app.Listen(":" + getEnv("NOTIFICATION_PORT", "8002")); err != nil {
		log.Fatalf("Failed to start notification worker server: %v", err)
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
