package models

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
)

type Statements struct {
	UserEmail   string
	OrderId     string
	Exchange    string
	OpenVal     string
	CloseVal    string
	Symbol      string
	CreatedTime time.Time `gorm:"type:timestamptz;default:now()" json:"created_at"`
	UpdatedTime time.Time `gorm:"type:timestamptz;default:now()" json:"updated_at"`
	Side        string
	ClosedPnl   string
	Quantity    string
}

func (st *Statements) FindStatements(db *gorm.DB, email string, service string) (*[]Statements, error) {
	statements := []Statements{}
	err := db.Debug().Model(Statements{}).Where("user_email = ? AND exchange = ?", email, service).Find(&statements).Error
	if err != nil {
		return &[]Statements{}, err
	}
	return &statements, nil
}

func (st *Statements) CalculateTodayProfit(db *gorm.DB, email string, exchange string) string {
	todayProfitQuery := `
		SELECT SUM(quantity) AS today_profit
		FROM statements
		WHERE DATE(created_at) = CURDATE()
		AND user_email = ? AND exchange = ?
	`

	var todayProfit sql.NullFloat64
	err := db.Debug().Raw(todayProfitQuery, email, exchange).Scan(&todayProfit).Error
	if err != nil {
		// Handle the error
	}

	if todayProfit.Valid {
		return strconv.FormatFloat(todayProfit.Float64, 'f', 2, 64)
	}
	return "0.00"
}

// Calculate total profit value
func (st *Statements) CalculateTotalProfit(db *gorm.DB, email string, exchange string) string {
	totalProfitQuery := `
		SELECT SUM(quantity) AS total_profit
		FROM statements
		WHERE user_email = ? AND exchange = ?
	`

	var totalProfit sql.NullFloat64
	err := db.Debug().Raw(totalProfitQuery, email, exchange).Scan(&totalProfit).Error
	if err != nil {
		// Handle the error
	}

	if totalProfit.Valid {
		return strconv.FormatFloat(totalProfit.Float64, 'f', 2, 64)
	}
	return "0.00"
}

// Filter statements for today's profit
func (st *Statements) FilterTodayProfitStatements(db *gorm.DB, email string, exchange string) []Statements {
	todayProfitStatementsQuery := `
		SELECT *
		FROM statements
		WHERE DATE(created_at) = CURDATE()
		AND user_email = ? AND exchange = ?
	`

	filteredStatements := []Statements{}
	err := db.Debug().Raw(todayProfitStatementsQuery, email, exchange).Scan(&filteredStatements).Error
	if err != nil {
		// Handle the error
	}

	return filteredStatements
}

// Calculate coinwise or symbol wise profit
func (st *Statements) CalculateCoinwiseProfit(db *gorm.DB, email string, exchange string) map[string]string {
	coinwiseProfitQuery := `
		SELECT symbol, SUM(quantity) AS coinwise_profit
		FROM statements
		WHERE user_email = ? AND exchange = ?
		GROUP BY symbol
	`

	type CoinwiseProfitResult struct {
		Symbol         string
		CoinwiseProfit sql.NullFloat64
	}

	coinwiseProfits := []CoinwiseProfitResult{}
	err := db.Debug().Raw(coinwiseProfitQuery, email, exchange).Scan(&coinwiseProfits).Error
	if err != nil {
		// Handle the error
	}

	coinwiseProfitMap := make(map[string]string)
	for _, result := range coinwiseProfits {
		if result.CoinwiseProfit.Valid {
			coinwiseProfitMap[result.Symbol] = strconv.FormatFloat(result.CoinwiseProfit.Float64, 'f', 2, 64)
		} else {
			coinwiseProfitMap[result.Symbol] = "0.00"
		}
	}

	return coinwiseProfitMap
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
