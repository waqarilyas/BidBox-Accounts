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
	Size        string    `json:"size"`
	Profit      float64   `json:"profit"`
	QuoteAmount float64   `json:"quote_amount"`
}

func (h *History) GetHistory(db *gorm.DB, email string, service string) (*[]History, error) {
	hist := []History{}
	err := db.Debug().Model(&History{}).Where("email = ? AND service = ?", email, service).Find(&hist).Error
	if err != nil {
		return &[]History{}, err
	}
	return &hist, nil

}

func (h *History) GetSuccessfulTrades(db *gorm.DB) (int, error) {

	var count int
	err := db.Debug().Model(&History{}).Count(&count).Error

	if err != nil {
		return 0, err
	}
	return count, nil

}
