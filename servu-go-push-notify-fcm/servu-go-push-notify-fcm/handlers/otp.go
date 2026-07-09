package handlers

import (
	"github.com/baryogenesis2025/servu-go/auth"
	"github.com/baryogenesis2025/servu-go/constants"
	"github.com/baryogenesis2025/servu-go/middleware"
	"github.com/gofiber/fiber/v2"
)

type OTPHandler struct {
	authService *auth.AuthService // Changed from twilioService
}

func NewOTPHandler(authService *auth.AuthService) *OTPHandler { // Changed parameter
	return &OTPHandler{
		authService: authService,
	}
}

type SendOTPRequest struct {
	PhoneNumber string `json:"phone_number"`
}

type VerifyOTPRequest struct {
	PhoneNumber string `json:"phone_number"`
	Code        string `json:"code"`
	DeviceID    string `json:"device_id"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type UpdateProfileRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}



// SendOTP sends verification code to phone number
// SendOTP godoc
// @Summary Send OTP
// @Description Sends a login OTP to the given phone number
// @Tags auth
// @Accept json
// @Produce json
// @Param body body SendOTPRequest true "Phone number payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/send-verification [post]
func (h *OTPHandler) SendOTP(c *fiber.Ctx) error {
	var req struct {
		PhoneNumber string `json:"phone_number"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": constants.InvalidRequestBody,
		})
	}

	if req.PhoneNumber == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Phone number is required",
		})
	}

	err := h.authService.SendLoginOTP(req.PhoneNumber)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to send OTP",
		})
	}

	return c.JSON(fiber.Map{
		"message":      "OTP sent successfully",
		"phone_number": req.PhoneNumber,
	})
}

// VerifyOTP verifies the OTP and returns authentication tokens
// VerifyOTP godoc
// @Summary Verify OTP and login
// @Description Verifies OTP and returns authentication tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param body body VerifyOTPRequest true "OTP verification payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/check-verification [post]
func (h *OTPHandler) VerifyOTP(c *fiber.Ctx) error {
	var req struct {
		PhoneNumber string `json:"phone_number"`
		Code        string `json:"code"`
		DeviceID    string `json:"device_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": constants.InvalidRequestBody,
		})
	}

	if req.PhoneNumber == "" || req.Code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Phone number and code are required",
		})
	}

	if req.DeviceID == "" {
		req.DeviceID = "web-browser"
	}

	response, err := h.authService.VerifyOTPAndLogin(req.PhoneNumber, req.Code, req.DeviceID)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Login successful",
		"data":    response,
	})
}

// RefreshToken generates new access token
// RefreshToken godoc
// @Summary Refresh access token
// @Description Generates a new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param body body RefreshTokenRequest true "Refresh token payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/refresh-token [post]
func (h *OTPHandler) RefreshToken(c *fiber.Ctx) error {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": constants.InvalidRequestBody,
		})
	}

	if req.RefreshToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Refresh token is required",
		})
	}

	response, err := h.authService.RefreshAccessToken(req.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Token refreshed successfully",
		"data":    response,
	})
}

// Logout revokes the refresh token
// Logout godoc
// @Summary Logout user
// @Description Revokes a refresh token and logs out the user
// @Tags auth
// @Accept json
// @Produce json
// @Param body body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/logout [post]
func (h *OTPHandler) Logout(c *fiber.Ctx) error {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": constants.InvalidRequestBody,
		})
	}

	if req.RefreshToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Refresh token is required",
		})
	}

	err := h.authService.Logout(req.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to logout",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Logged out successfully",
	})
}

// LogoutAllDevices revokes all refresh tokens
// LogoutAllDevices godoc
// @Summary Logout from all devices
// @Description Revokes all refresh tokens for the authenticated user
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/logout-all [post]
func (h *OTPHandler) LogoutAllDevices(c *fiber.Ctx) error {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	err := h.authService.LogoutAllDevices(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to logout from all devices",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Logged out from all devices successfully",
	})
}

// GetProfile returns user profile
// GetProfile godoc
// @Summary Get user profile
// @Description Retrieves authenticated user's profile
// @Tags profile
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Router /profile [get]
func GetProfile(c *fiber.Ctx) error {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Profile retrieved",
		"user_id": userID,
	})
}

// UpdateProfile updates user profile
// UpdateProfile godoc
// @Summary Update profile
// @Description Updates authenticated user's profile
// @Tags profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body UpdateProfileRequest true "Profile update payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /profile [put]
func UpdateProfile(c *fiber.Ctx) error {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": constants.InvalidRequestBody,
		})
	}

	return c.JSON(fiber.Map{
		"message": "Profile updated",
		"user_id": userID,
	})
}
