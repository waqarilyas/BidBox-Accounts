package admin

import (
	uuid "github.com/satori/go.uuid"
)

type Admin struct {
	Id       uuid.UUID `gorm:"type:uuid;primary_key;"`
	Email    string    `gorm:"type:text;not null"`
	Password string    `gorm:"not null"`

	OtpVerified bool `gorm:"default:false;"`

	Otp_secret   string `gorm:"-"`
	Otp_auth_url string `gorm:"-"`
}

type RegisterUserInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" bindinig:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginUserInput struct {
	Email    string `json:"email" bindinig:"required"`
	Password string `json:"password" binding:"required"`
}

type OTPInput struct {
	UserId string `json:"user_id"`
	Token  string `json:"token"`
}
