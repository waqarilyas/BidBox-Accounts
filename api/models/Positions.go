package models

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

type Position struct {
	Id              int       `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	CreatedAt       time.Time `gorm:"type:timestamptz;default:now()" json:"created_at"`
	UpdatedAt       time.Time `gorm:"type:timestamptz;default:now()" json:"updated_at"`
	Symbol          string    `json:"symbol"`
	Leverage        string    `json:"leverage"`
	OpenPrice       string    `json:"open_price"`
	LiqPrice        string    `json:"liq_price"`
	TakeProfit      string    `json:"take_profit"`
	MarkPrice       string    `json:"mark_price"`
	StopLoss        string    `json:"stop_loss"`
	UnrealizedPl    string    `json:"unrealized_pl"`
	Side            string    `json:"side"`
	Size            string    `json:"size"`
	Margin          string    `json:"margin"`
	UserEmail       string    `gorm:"not null" json:"user_email"`
	Status          string    `gorm:"default:'opened'" json:"status"`
	Exchange        string    `json:"exchange"`
	LastUpdatePrice string    `json:"last_update_price"`
	OrderId         string    `json:"order_id"`
	Layer           int       `json:"layer"`
	TotalProfit     float64   `json:"total_profit"`
	FirstBuyAmount  string    `json:"first_buy_amount"`
	HedgeId         string    `json:"hedge_id"`
	Fee             float64   `json:"fee"`
}

type GroupedPosition struct {
	HedgeID   string
	Positions []Position
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

func (u *Position) GetOpenPositions(db *gorm.DB, email string, service string) (*[]Position, error) {
	pos := []Position{}

	err := db.Model(Position{}).Where("user_email = ? AND exchange = ? AND status = ?", email, service, "opened").Order("created_at DESC").Find(&pos).Error
	if err != nil {
		return &[]Position{}, err
	}

	return &pos, nil
}

func (u *Position) GetClosedGroupedPositions(db *gorm.DB, email string, service string) ([]GroupedPosition, error) {
	var positions []Position
	err := db.Model(Position{}).
		Where("user_email = ? AND exchange = ? AND status = ?", email, service, "closed").Order("created_at DESC").
		Find(&positions).Error
	if err != nil {
		fmt.Println("🚀 ~ file: Positions.go:131 ~ func ~ err:", err)
		return nil, err
	}

	groupedPositions := make([]GroupedPosition, 0)
	positionMap := make(map[string][]Position) // Map to temporarily store positions by hedge_id

	// Group positions by hedge_id
	for _, position := range positions {
		positionMap[string(position.HedgeId)] = append(positionMap[position.HedgeId], position)
	}

	// Convert the map to an array of GroupedPosition
	for hedgeID, pos := range positionMap {
		groupedPositions = append(groupedPositions, GroupedPosition{
			HedgeID:   hedgeID,
			Positions: pos,
		})
	}

	return groupedPositions, nil
}
