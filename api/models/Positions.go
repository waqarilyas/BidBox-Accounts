package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Position struct {
	Id        int       `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	CreatedAt time.Time `gorm:"type:timestamptz;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamptz;default:now()" json:"updated_at"`
	Symbol    string    `json:"symbol"`
	Leverage  string    `json:"leverage"`
	// OpenPrice    float64   `json:"open_price"`
	// LiqPrice     float64   `json:"liq_price"`
	// TakeProfit   float64   `json:"take_profit"`
	// StopLoss     float64   `json:"stop_loss"`
	// UnrealizedPl float32   `json:"unrealized_pl"`
	// Markprice    float64   `json:"mark_price"`
	Side      string `json:"side"`
	Size      string `json:"size"`
	Margin    string `json:"margin"`
	UserEmail string `gorm:"not null" json:"user_email"`
	Status    string `gorm:"default:'opened'" json:"status"`
	Exchange  string
	Profit    string
	OrderId   string
}

func (position *Position) CreateNewPosition(db *gorm.DB) (*Position, error) {
	err := db.Create(&position).Error
	if err != nil {
		return &Position{}, err
	}
	return position, nil
}

func (u *Position) GetAllPositions(db *gorm.DB) (*[]Position, error) {
	positions := []Position{}
	err := db.Model(&Key{}).Limit(100).Find(&positions).Error
	if err != nil {
		return &[]Position{}, err
	}
	return &positions, nil
}

func (p *Position) GetOrderByEmail(db *gorm.DB, email string) (*[]Position, error) {
	pos := []Position{}
	err := db.Model(&Position{}).Where("user_email = ?", email).Find(&pos).Error
	if err != nil {
		return &[]Position{}, err
	}
	return &pos, nil
}

func (p *Position) UpdatePosition(db *gorm.DB, id int) error {
	db = db.Model(&Position{}).Where("id = ?", id).Take(&Position{}).UpdateColumns(
		map[string]interface{}{
			"status": "closed",
		},
	)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (p *Position) UpdateProfits(db *gorm.DB, id int, profit string) error {
	db = db.Model(&Position{}).Where("id = ?", id).Take(&Position{}).UpdateColumns(
		map[string]interface{}{
			"profit": profit,
		},
	)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (u *Position) GetActivePositions(db *gorm.DB) (int, error) {
	var count int

	err := db.Model(Position{}).Where("status = ?", "opened").Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (u *Position) GetSuccessfulPositions(db *gorm.DB) (int, error) {
	var count int

	err := db.Model(Position{}).Where("status = ?", "closed").Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (u *Position) GetClosedPositions(db *gorm.DB, email string, service string) (*[]Position, error) {
	pos := []Position{}

	err := db.Model(Position{}).Where("user_email = ? AND exchange = ? AND status = ?", email, service, "closed").Find(&pos).Error
	if err != nil {
		return &[]Position{}, err
	}

	return &pos, nil
}

func (u *Position) GetClosePositions(db *gorm.DB, email string, exchange string) (*[]Position, error) {
	pos := []Position{}

	err := db.Model(Position{}).Where("status = ? AND exchange = ? AND user_email = ? ", "closed", exchange, email).Find(&pos).Error
	if err != nil {
		return &[]Position{}, err
	}

	return &pos, nil
}