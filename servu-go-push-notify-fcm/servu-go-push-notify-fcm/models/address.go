package models

import "time"

type Address struct {
	ID        uint     `json:"id" gorm:"primaryKey"`
	UserID    string   `json:"user_id" gorm:"index"`
	Address1  string   `json:"address_line_1"`
	Address2  string   `json:"address_line_2,omitempty"`
	Locality  string   `json:"locality,omitempty"`
	City      string   `json:"city,omitempty"`
	State     string   `json:"state,omitempty"`
	Pincode   string   `json:"pincode,omitempty"`
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
	IsDefault bool     `json:"is_default"`
	CreatedAt time.Time
}

func NewAddress(userID, address1, address2, locality, city, state, pincode string, latitude, longitude *float64, isDefault bool) *Address {
	return &Address{
		UserID:    userID,
		Address1:  address1,
		Address2:  address2,
		Locality:  locality,
		City:      city,
		State:     state,
		Pincode:   pincode,
		Latitude:  latitude,
		Longitude: longitude,
		IsDefault: isDefault,
	}
}
