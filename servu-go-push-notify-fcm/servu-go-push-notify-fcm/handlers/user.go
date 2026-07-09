package handlers

import (
	"net/http"

	"github.com/baryogenesis2025/servu-go/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)


// HandleUser godoc
// @Summary Create or fetch user
// @Description Creates a new user if not exists, otherwise returns existing user based on email or phone
// @Tags users
// @Accept json
// @Produce json
// @Param body body models.User true "User payload"
// @Success 200 {object} map[string]interface{} "User already exists"
// @Success 201 {object} map[string]interface{} "New user created"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{phonenumber} [get]
func HandleUser(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var input models.User
		if err := c.BodyParser(&input); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}
		var existing models.User
		result := db.Where("email = ? OR phone = ?", input.Email, input.Phone).First(&existing)
		if result.Error == nil {
			return c.Status(http.StatusOK).JSON(fiber.Map{
				"message": "User already exists",
				"user":    existing,
			})
		}
		if result.Error == gorm.ErrRecordNotFound {
			newUser := models.NewUser("", input.Name, input.Email, input.Phone, input.Role, input.About)

			if err := db.Create(newUser).Error; err != nil {
				return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to create user",
				})
			}

			return c.Status(http.StatusCreated).JSON(fiber.Map{
				"message": "New user created successfully",
				"user":    newUser,
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error: " + result.Error.Error(),
		})
	}
}
