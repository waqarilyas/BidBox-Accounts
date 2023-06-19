package admin

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type Admin struct {
	Id       uuid.UUID `gorm:"type:uuid;primary_key;"`
	Email    string    `gorm:"type:text;not null"`
	Password string    `gorm:"not null"`

	OtpVerified bool `gorm:"default:false;"`

	OtpSecret  string
	OtpUrl     string
	OtpEnabled string
}

type RegisterUserInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" bindinig:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginUserInput struct {
	UserId   string `json:"UserId"`
	Email    string `json:"email" bindinig:"required"`
	Password string `json:"password" binding:"required"`
}

type OTPInput struct {
	UserId string `json:"id"`
	Token  string `json:"token"`
}

type OTPResponse struct {
	OPTEnabled string `json:"otp_enabled"`
}

type ChangePasswordInput struct {
	Password    string `json:"password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type ChangeTimeframeInput struct {
	Timeframe string `json:"timeframe" binding:"required"`
}

type Resp struct {
	Id          uuid.UUID `json:"id"`
	Email       string    `json:"email"`
	OtpVerified bool      `json:"otp_verified"`
}

func (a *Admin) UpdateOtp(db *gorm.DB, val string) (*Resp, error) {
	cp := Admin{}
	db = db.Model(&Admin{}).
		UpdateColumns(
			map[string]interface{}{
				"otp_enabled": val,
			},
		).Take(&cp)
	if db.Error != nil {
		return &Resp{}, db.Error
	}
	res := Resp{
		Id:          cp.Id,
		Email:       cp.Email,
		OtpVerified: cp.OtpVerified,
	}
	return &res, nil
}
