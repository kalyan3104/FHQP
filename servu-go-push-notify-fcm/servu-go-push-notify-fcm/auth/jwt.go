package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lucsky/cuid"
	"gorm.io/gorm"
)

// Token configuration
const (
	AccessTokenDuration  = 30 * time.Minute
	RefreshTokenDuration = 90 * 24 * time.Hour
)

// User model
type User struct {
	ID            string    `gorm:"primaryKey" json:"id"`
	PhoneNumber   string    `gorm:"uniqueIndex;not null" json:"phone_number"`
	PhoneVerified bool      `gorm:"default:false" json:"phone_verified"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Role          string    `gorm:"default:'customer'" json:"role"` // customer, provider, admin
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// RefreshToken model for tracking refresh tokens
type RefreshToken struct {
	ID        string    `gorm:"primaryKey"`
	UserID    string    `gorm:"index;not null"`
	Token     string    `gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"index"`
	DeviceID  string    `json:"device_id"` // Track devices
	CreatedAt time.Time
	Revoked   bool `gorm:"default:false;index"`
}

// JWT Claims
type Claims struct {
	UserID      string `json:"user_id"`
	PhoneNumber string `json:"phone_number"`
	Role        string `json:"role"`
	jwt.RegisteredClaims
}

// Auth responses
type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	User         *User  `json:"user"`
}

type AuthService struct {
	db            *gorm.DB
	twilioService *TwilioService
	jwtSecret     []byte
}

func NewAuthService(db *gorm.DB, twilioService *TwilioService, jwtSecret string) *AuthService {
	return &AuthService{
		db:            db,
		twilioService: twilioService,
		jwtSecret:     []byte(jwtSecret),
	}
}

// =============================================================================
// AUTHENTICATION FLOW
// =============================================================================

// Step 1: Send OTP to phone number
func (s *AuthService) SendLoginOTP(phoneNumber string) error {
	return s.twilioService.SendOTP(phoneNumber)
}

// Step 2: Verify OTP and return tokens
func (s *AuthService) VerifyOTPAndLogin(phoneNumber, code, deviceID string) (*AuthResponse, error) {
	// Verify OTP with Twilio
	valid, err := s.twilioService.VerifyOTP(phoneNumber, code)
	if err != nil || !valid {
		return nil, errors.New("invalid or expired OTP")
	}

	// Find or create user
	var user User
	result := s.db.Where("phone_number = ?", phoneNumber).First(&user)

	if result.Error == gorm.ErrRecordNotFound {
		// New user - create account
		user = User{
			ID:            cuid.New(),
			PhoneNumber:   phoneNumber,
			PhoneVerified: true,
			Role:          "customer",
		}
		if err := s.db.Create(&user).Error; err != nil {
			return nil, err
		}
	} else if result.Error != nil {
		return nil, result.Error
	} else {
		// Existing user - mark phone as verified
		user.PhoneVerified = true
		s.db.Save(&user)
	}

	// Generate and return tokens
	return s.generateTokens(&user, deviceID)
}

// =============================================================================
// TOKEN GENERATION
// =============================================================================

// Generate access token (JWT) and refresh token
func (s *AuthService) generateTokens(user *User, deviceID string) (*AuthResponse, error) {
	// 1. Generate JWT Access Token
	accessToken, err := s.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}

	// 2. Generate opaque Refresh Token
	refreshToken, err := generateSecureToken(32)
	if err != nil {
		return nil, err
	}

	// 3. Store refresh token in database
	rt := RefreshToken{
		ID:        cuid.New(),
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(RefreshTokenDuration),
		DeviceID:  deviceID,
		Revoked:   false,
	}
	if err := s.db.Create(&rt).Error; err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(AccessTokenDuration.Seconds()),
		User:         user,
	}, nil
}

// Generate JWT access token
func (s *AuthService) GenerateAccessToken(user *User) (string, error) {
	claims := Claims{
		UserID:      user.ID,
		PhoneNumber: user.PhoneNumber,
		Role:        user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

// =============================================================================
// TOKEN VALIDATION & VERIFICATION
// =============================================================================

// Verify and parse JWT access token
func (s *AuthService) VerifyAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// Get user from token
func (s *AuthService) GetUserFromToken(tokenString string) (*User, error) {
	claims, err := s.VerifyAccessToken(tokenString)
	if err != nil {
		return nil, err
	}

	var user User
	if err := s.db.First(&user, claims.UserID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	return &user, nil
}

// =============================================================================
// TOKEN REFRESH
// =============================================================================

// Refresh access token using refresh token
func (s *AuthService) RefreshAccessToken(refreshToken string) (*AuthResponse, error) {
	// Find valid refresh token
	var rt RefreshToken
	err := s.db.Where("token = ? AND revoked = false AND expires_at > ?",
		refreshToken, time.Now()).First(&rt).Error

	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	// Get user
	var user User
	if err := s.db.First(&user, rt.UserID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// Generate new access token
	accessToken, err := s.GenerateAccessToken(&user)
	if err != nil {
		return nil, err
	}

	// Return tokens (refresh token stays the same)
	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(AccessTokenDuration.Seconds()),
		User:         &user,
	}, nil
}

// =============================================================================
// TOKEN REVOCATION (Logout)
// =============================================================================

// Logout - revoke refresh token
func (s *AuthService) Logout(refreshToken string) error {
	return s.db.Model(&RefreshToken{}).
		Where("token = ?", refreshToken).
		Update("revoked", true).Error
}

// Logout all devices for a user
func (s *AuthService) LogoutAllDevices(userID string) error {
	return s.db.Model(&RefreshToken{}).
		Where("user_id = ? AND revoked = false", userID).
		Update("revoked", true).Error
}

// =============================================================================
// HELPER FUNCTIONS
// =============================================================================

// Generate cryptographically secure random token
func generateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// Clean up expired tokens (run as cron job)
func (s *AuthService) CleanupExpiredTokens() error {
	return s.db.Where("expires_at < ? OR revoked = true", time.Now()).
		Delete(&RefreshToken{}).Error
}
