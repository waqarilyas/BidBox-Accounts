package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type History struct {
	Id          uuid.UUID `json:"Id"`
	Symbol      string    `json:"symbol"`
	Email       string    `json:"email"`
	Side        string    `json:"side"`
	Leverage    int       `json:"leverage"`
	Service     string    `json:"service"`
	MarginMode  string    `json:"margin_mod"`
	MarketPrice string    `json:"market_price"`
}

func (h *History) GetHistory(db *gorm.DB, email string, service string) (*[]History, error) {
	hist := []History{}
	err := db.Debug().Model(&History{}).Where("email = ? AND service = ?", email, service).Find(&hist).Error
	if err != nil {
		return &[]History{}, err
	}
	return &hist, nil

}
