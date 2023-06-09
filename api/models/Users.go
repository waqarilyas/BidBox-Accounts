package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type User struct {
	Id        uuid.UUID `gorm:"primary_key;auto_increment" json:"id"`
	Email     string    `gorm:"size:100;not null;unique" json:"email"`
	Password  string    `gorm:"size:100;not null;" json:"password"`
	Confirmed bool      `gorm:"default:false" json:"confirmed"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Strategy  string    `json:"strategy"`
	Mode      string    `json:"mode"`
	Name      string    `gorm:"default:''" json:"name"`
	Country   string    `gorm:"default:''" json:"country"`
	UserName  string    `gorm:"default:''" json:"user_name"`
	TimeZone  string    `gorm:"default:''" json:"time_zone"`
	Phone     string    `gorm:"default:''" json:"phone"`
}

func (u *User) FindUserById(db *gorm.DB, uid uint32) (*User, error) {
	err := db.Debug().Model(User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &User{}, errors.New("User not found")
	}
	return u, nil
}

func (u *User) FindUserByEmail(db *gorm.DB, email string) (*User, error) {
	err := db.Debug().Model(User{}).Where("email = ?", email).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &User{}, errors.New("User not found")
	}
	return u, nil
}

func (u *User) ChangeStrategy(db *gorm.DB, strategy string) (*User, error) {
	db = db.Debug().Model(&User{}).Where("email = ?", u.Email).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"strategy":   strategy,
			"updated_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &User{}, db.Error
	}
	// This is the display the updated user
	err := db.Debug().Model(&User{}).Where("email = ?", u.Email).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) ChangeMode(db *gorm.DB, mode string) (*User, error) {
	db = db.Debug().Model(&User{}).Where("email = ?", u.Email).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"mode":       mode,
			"updated_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &User{}, db.Error
	}
	// This is the display the updated user
	err := db.Debug().Model(&User{}).Where("email = ?", u.Email).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (k *User) ValidateUser(prev *User) error {
	if k.Name == "" {
		k.Name = prev.Name
	}
	if k.Country == "" {
		k.Country = prev.Country
	}
	if k.UserName == "" {
		k.UserName = prev.UserName
	}
	if k.TimeZone == "" {
		k.TimeZone = prev.TimeZone
	}

	if k.Phone == prev.Phone {
		k.Phone = prev.Phone
	}

	return errors.New("strategy or mode is incorrect")
}
