package models

import (
	"errors"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type Accounts struct {
	Id                 uuid.UUID `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	CreatedAt          string    `gorm:"size:255;not null" json:"created_at"`
	UpdatedAt          string    `gorm:"not null" json:"updated_at"`
	MarginCoin         string    `gorm:"not null" json:"margin_coin"`
	AvailableBalance   string    `gorm:"not null" json:"available_balance"`
	TotalMarginBalance string    `gorm:"not null" json:"total_margin_balance"`
	MarginValueUSDT    string    `gorm:"not null" json:"margin_value_usdt"`
	MarginValueBTC     string    `gorm:"not null" json:"margin_value_btc"`
	FloatingPnl        string    `gorm:"not null" json:"floating_pnl"`
	ApiKeyId           uuid.UUID `gorm:"not null;type:uuid" json:"api_key_id"`
}

func (e *Accounts) ValidateAccountDetails() error {
	if e.MarginCoin == "" {
		return errors.New("margin_coin is required")
	}
	if e.AvailableBalance == "" {
		return errors.New("available_balance is required")
	}
	if e.TotalMarginBalance == "" {
		return errors.New("total_margin_balance is required")
	}
	if e.MarginValueUSDT == "" {
		return errors.New("margin_value_usdt is required")
	}
	if e.MarginValueBTC == "" {
		return errors.New("margin_value_btc is required")
	}
	if e.FloatingPnl == "" {
		return errors.New("FloatingPnl is required")
	}
	// if e.ApiKeyId == "" {
	// 	return errors.New("api_key_id is required")
	// }
	return nil
}

func (e *Accounts) SaveAccount(db *gorm.DB) (*Accounts, error) {
	account := Accounts{}
	err := db.Debug().Where("api_key_id = ? AND margin_coin = ?", e.ApiKeyId, e.MarginCoin).Assign(&e).FirstOrCreate(&account).Error
	if err != nil {
		return &Accounts{}, err
	}
	return e, nil
}

func (e *Accounts) GetAccountByApiKeyId(db *gorm.DB, apiKeyId uuid.UUID) (*Accounts, error) {
	err := db.Debug().Model(Accounts{}).Where("api_key_id = ?", apiKeyId).Take(&e).Error
	if err != nil {
		return &Accounts{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Accounts{}, errors.New("account not found")
	}
	return e, nil
}
