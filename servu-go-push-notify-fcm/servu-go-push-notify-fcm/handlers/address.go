package handlers

import (
	"github.com/baryogenesis2025/servu-go/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)


// CreateAddress godoc
// @Summary Create a new address
// @Description Create a new address for a user. If marked as default, existing default will be unset.
// @Tags address
// @Accept json
// @Produce json
// @Param body body models.Address true "Address payload"
// @Success 201 {object} models.Address
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /address/ [post]
func CreateAddress(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var address models.Address

		if err := c.BodyParser(&address); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		if address.UserID == "" || address.Address1 == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "UserID and Address1 are required",
			})
		}

		// Handle default address
		if address.IsDefault {
			db.Model(&models.Address{}).
				Where("user_id = ?", address.UserID).
				Update("is_default", false)
		}

		if err := db.Create(&address).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create address",
			})
		}

		return c.Status(fiber.StatusCreated).JSON(address)
	}
}


// GetUserAddresses godoc
// @Summary Get all addresses of a user
// @Description Fetch all addresses belonging to a specific user
// @Tags address
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {array} models.Address
// @Failure 500 {object} map[string]string
// @Router /address/user/{userId} [get]
func GetUserAddresses(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Params("userId")

		var addresses []models.Address
		if err := db.
			Where("user_id = ?", userID).
			Order("is_default DESC, created_at DESC").
			Find(&addresses).Error; err != nil {

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch addresses",
			})
		}

		return c.JSON(addresses)
	}
}


// GetAddressByID godoc
// @Summary Get address by ID
// @Description Fetch a specific address using its ID
// @Tags address
// @Produce json
// @Param id path string true "Address ID"
// @Success 200 {object} models.Address
// @Failure 404 {object} map[string]string
// @Router /address/{id} [get]
func GetAddressByID(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		var address models.Address
		if err := db.First(&address, id).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Address not found",
			})
		}

		return c.JSON(address)
	}
}


// UpdateAddress godoc
// @Summary Update address
// @Description Update an existing address. If marked as default, other defaults will be unset.
// @Tags address
// @Accept json
// @Produce json
// @Param id path string true "Address ID"
// @Param body body models.Address true "Updated address payload"
// @Success 200 {object} models.Address
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /address/{id} [put]
func UpdateAddress(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		var address models.Address
		if err := db.First(&address, id).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Address not found",
			})
		}

		var input models.Address
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		// If setting default, unset others
		if input.IsDefault {
			db.Model(&models.Address{}).
				Where("user_id = ?", address.UserID).
				Update("is_default", false)
		}

		if err := db.Model(&address).Updates(input).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update address",
			})
		}

		return c.JSON(address)
	}
}


// DeleteAddress godoc
// @Summary Delete address
// @Description Delete an address by ID
// @Tags address
// @Produce json
// @Param id path string true "Address ID"
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /address/{id} [delete]
func DeleteAddress(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		if err := db.Delete(&models.Address{}, id).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to delete address",
			})
		}

		return c.JSON(fiber.Map{
			"message": "Address deleted successfully",
		})
	}
}
