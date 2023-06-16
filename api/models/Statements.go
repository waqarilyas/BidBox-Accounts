package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Statements struct {
	StatementCreatedAt time.Time
	UserEmail          string
	OrderId            string
	Exchange           string
	OpenVal            string
	CloseVal           string
	Symbol             string
	CreatedTime        time.Time `gorm:"type:timestamptz;default:now()" json:"created_at"`
	UpdatedTime        time.Time `gorm:"type:timestamptz;default:now()" json:"updated_at"`
	Side               string
	ClosedPnl          string
	Quantity           string
}

func (st *Statements) FindStatements(db *gorm.DB, email string, service string) (*[]Statements, error) {
	statements := []Statements{}
	err := db.Debug().Model(Statements{}).Where("user_email = ? AND exchange = ?", email, service).Find(&statements).Error
	if err != nil {
		return &[]Statements{}, err
	}
	return &statements, nil
}

func (st *Statements) GetOrderThisMonth(db *gorm.DB) (*[]Statements, error) {
	orders := []Statements{}

	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)

	err := db.Model(Statements{}).Order("closed_pnl DESC").Where("created_time >= ?", startOfMonth).Find(&orders).Error
	if err != nil {
		return &[]Statements{}, err
	}

	return &orders, nil
}

func (st *Statements) GetOrderThisDay(db *gorm.DB) (*[]Statements, error) {
	ords := []Statements{}

	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.UTC)

	err := db.Model(Statements{}).Order("closed_pnl DESC").Where("created_time >= ? AND created_time <= ?", startOfDay, endOfDay).Find(&ords).Error
	if err != nil {
		return &[]Statements{}, err
	}

	return &ords, nil
}

func (st *Statements) GetOrderAllTime(db *gorm.DB) (*[]Statements, error) {
	ords := []Statements{}

	err := db.Model(Statements{}).Order("closed_pnl DESC").Find(&ords).Error
	if err != nil {
		return &[]Statements{}, err
	}

	return &ords, nil
}

func (st *Statements) GetStatementsToday(db *gorm.DB, email string, service string, limit int, offset int) (*[]Statements, error) {
	ords := []Statements{}

	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.UTC)

	err := db.Model(Statements{}).Where("user_email = ? AND exchange = ? AND created_time >= ? AND created_time <= ?", email, service, startOfDay, endOfDay).Limit(limit).Offset(offset).Find(&ords).Error
	if err != nil {
		return &[]Statements{}, err
	}

	return &ords, nil
}

func (st *Statements) GetStatementsAllTime(db *gorm.DB, email string, service string, limit int, offset int) (*[]Statements, error) {
	ords := []Statements{}

	err := db.Model(Statements{}).Where("user_email = ? AND exchange = ?", email, service).Limit(limit).Offset(offset).Find(&ords).Error
	if err != nil {
		return &[]Statements{}, err
	}

	return &ords, nil
}

func (st *Statements) GetCoinwiseToday(db *gorm.DB, email string, service string) (*[]Statements, error) {
	ords := []Statements{}

	err := db.Model(Statements{}).Select("SUM (closed_pnl)").Where("user_email = ? AND exchange = ?", email, service).Group("symbol").Find(&ords).Error
	if err != nil {
		return &[]Statements{}, err
	}

	return &ords, nil
}
