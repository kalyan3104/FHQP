package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type TwilioService struct {
	accountSID string
	authToken  string
	serviceID  string
	client     *http.Client
}

func NewTwilioService() *TwilioService {
	return &TwilioService{
		accountSID: os.Getenv("TWILIO_ACCOUNT_SID"),
		authToken:  os.Getenv("TWILIO_AUTH_TOKEN"),
		serviceID:  os.Getenv("TWILIO_VERIFY_SERVICE_SID"),
		client:     &http.Client{},
	}
}

type VerificationResponse struct {
	SID     string `json:"sid"`
	Status  string `json:"status"`
	Valid   bool   `json:"valid"`
	Message string `json:"message,omitempty"`
}

// SendOTP sends an OTP to the provided phone number
func (t *TwilioService) SendOTP(phoneNumber string) error {
	if t.accountSID == "" || t.authToken == "" || t.serviceID == "" {
		return errors.New("twilio credentials not configured")
	}

	apiURL := fmt.Sprintf(
		"https://verify.twilio.com/v2/Services/%s/Verifications",
		t.serviceID,
	)

	data := url.Values{}
	data.Set("To", phoneNumber)
	data.Set("Channel", "sms")

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(t.accountSID, t.authToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := t.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send OTP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		var errResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errResp)
		return fmt.Errorf("twilio error: %v", errResp)
	}

	return nil
}

// VerifyOTP verifies the OTP code for a phone number
func (t *TwilioService) VerifyOTP(phoneNumber, code string) (bool, error) {
	if t.accountSID == "" || t.authToken == "" || t.serviceID == "" {
		return false, errors.New("twilio credentials not configured")
	}

	apiURL := fmt.Sprintf(
		"https://verify.twilio.com/v2/Services/%s/VerificationCheck",
		t.serviceID,
	)

	data := url.Values{}
	data.Set("To", phoneNumber)
	data.Set("Code", code)

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(t.accountSID, t.authToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := t.client.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to verify OTP: %w", err)
	}
	defer resp.Body.Close()

	var result VerificationResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, fmt.Errorf("failed to decode response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("verification failed: %s", result.Status)
	}

	return result.Valid && result.Status == "approved", nil
}
