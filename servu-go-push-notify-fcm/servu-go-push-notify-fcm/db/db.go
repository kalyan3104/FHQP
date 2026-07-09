package db

import (
	"fmt"
	"log"
	"os"

	"github.com/baryogenesis2025/servu-go/models"
	"github.com/baryogenesis2025/servu-go/seeders"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() *gorm.DB { // Added return type
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found or could not be loaded")
	}

	// Get DSN
	dsn := os.Getenv("SUPABASE_DB_URL")
	if dsn == "" {
		log.Fatal("SUPABASE_DB_URL not set in environment")
	}
	fmt.Println("Connecting to Supabase...", dsn)

	// Connect to Supabase
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	
	if err != nil {
		log.Fatalf("Failed to connect to Supabase: %v", err)
	}
	db.AutoMigrate(models.AllModels...)

	seed_err := seeders.SeedServices(db)

	if seed_err != nil {
		panic(seed_err)
	}

	DB = db
	fmt.Println("✅ Connected to Supabase successfully!")

	return db // Added return statement
}
