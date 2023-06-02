package admin

import (
	uuid "github.com/satori/go.uuid"
)

type Admin struct {
	Id       uuid.UUID `gorm:"type:uuid;primary_key;"`
	Email    string    `gorm:"type:text;not null"`
	Password string    `gorm:"not null"`

	OtpVerified bool `gorm:"default:false;"`

	OtpSecret string
	OtpUrl    string
	OtpEnabled bool
}

type RegisterUserInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" bindinig:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginUserInput struct {
	UserId string `json:"UserId"`
	Email    string `json:"email" bindinig:"required"`
	Password string `json:"password" binding:"required"`
}

type OTPInput struct {
	UserId string `json:"id"`
	Token  string `json:"token"`
}

type OTPResponse struct {
	UserId string `json:"UserId"`
	OPTEnabled bool
}

type ChangePasswordInput struct {
	Password string `json:"password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}