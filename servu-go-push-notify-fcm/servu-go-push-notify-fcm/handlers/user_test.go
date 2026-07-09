package handlers

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/baryogenesis2025/servu-go/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestHandleUserCreatesUserWithGeneratedID(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		t.Fatalf("failed to migrate user model: %v", err)
	}

	app := fiber.New()
	app.Get("/users/:phonenumber", HandleUser(db))

	body := `{"name":"Jane Doe","email":"jane@example.com","phone":"1234567890","role":"customer","about":"test"}`
	req := httptest.NewRequest("GET", "/users/1234567890", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusCreated {
		t.Fatalf("expected status %d, got %d", fiber.StatusCreated, resp.StatusCode)
	}
}
