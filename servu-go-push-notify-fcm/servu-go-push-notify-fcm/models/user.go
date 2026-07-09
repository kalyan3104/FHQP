package models

import (
	"time"

	"github.com/lucsky/cuid"
)

type User struct {
	Id        string `json:"id" gorm:"primaryKey"`
	Name      string `json:"name"`
	Email     string `json:"email" gorm:"uniqueIndex"`
	Phone     string `json:"phone" gorm:"uniqueIndex"`
	About     string `json:"about"`
	Role      string `json:"role"`
	CreatedAt time.Time
	Addresses []Address `json:"addresses,omitempty" gorm:"foreignKey:UserID"`
}

func NewUser(id, name, email, phone, role, about string) *User {
	userID := id
	if userID == "" {
		userID = cuid.New()
	}

	return &User{
		Id:    userID,
		Name:  name,
		Email: email,
		Phone: phone,
		Role:  role,
		About: about,
	}
}
