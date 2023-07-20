package models

import (
	"sort"
	"time"

	"github.com/jinzhu/gorm"
)

type Statements struct {
	UserEmail   string    `json:"user_email"`
	Exchange    string    `json:"exchange"`
	Symbol      string    `json:"symbol"`
	CreatedTime time.Time `json:"created_time"`
	UpdatedTime time.Time `json:"updated_time"`
	Side        string    `json:"side"`
	ClosedPnl   float64   `json:"closed_pnl"`
	Size        float64   `json:"size"`
	PositionId  int       `json:"position_id"`
	QuoteAmount float64   `json:"quote_amount"`
	ProfitUSD   float64   `json:"profit_usd"`
}

type LeaderboardUser struct {
	UserEmail string
	Name      string
	ClosedPnl float64
	Country   string
	Username  string
}

type LeaderboardAPI struct {
	db *gorm.DB
}

func NewLeaderboardAPI(db *gorm.DB) *LeaderboardAPI {
	return &LeaderboardAPI{db: db}
}

func (api *LeaderboardAPI) GetLeaderboardToday() ([]LeaderboardUser, error) {
	users := []User{}
	err := api.db.Find(&users).Error
	if err != nil {
		return nil, err
	}

	leaderboard := make([]LeaderboardUser, len(users))
	for i, user := range users {
		orders, err := api.getOrdersToday(user.Email)
		if err != nil {
			return nil, err
		}

		closedPnl := api.calculateCumulativePnl(orders)
		leaderboard[i] = LeaderboardUser{
			UserEmail: user.Email,
			Name:      user.Name,
			ClosedPnl: closedPnl,
			Country:   user.Country,
			Username:  user.UserName,
		}
	}

	sort.Slice(leaderboard, func(i, j int) bool {
		return leaderboard[i].ClosedPnl > leaderboard[j].ClosedPnl
	})

	return leaderboard, nil
}

func (api *LeaderboardAPI) getOrdersToday(email string) ([]Statements, error) {
	orders := []Statements{}

	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.UTC)

	err := api.db.Model(Statements{}).Where("user_email = ? AND created_time >= ? AND created_time <= ?", email, startOfDay, endOfDay).Find(&orders).Error
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (api *LeaderboardAPI) calculateCumulativePnl(orders []Statements) float64 {
	var pnl float64
	for _, order := range orders {
		// closedPnl, err := (order.Size, 64)
		// if err != nil {
		// 	// Handle parsing error if needed
		// 	fmt.Println("----error close dpnl float format ---", err)
		// 	continue
		// }
		pnl += order.Size
	}
	return pnl
}

func (api *LeaderboardAPI) GetLeaderboardThisWeek() ([]LeaderboardUser, error) {
	users := []User{}
	err := api.db.Find(&users).Error
	if err != nil {
		return nil, err
	}

	leaderboard := make([]LeaderboardUser, len(users))
	for i, user := range users {
		orders, err := api.getOrdersThisWeek(user.Email)
		if err != nil {
			return nil, err
		}

		closedPnl := api.calculateCumulativePnl(orders)
		leaderboard[i] = LeaderboardUser{
			UserEmail: user.Email,
			Name:      user.Name,
			ClosedPnl: closedPnl,
			Country:   user.Country,
			Username:  user.UserName,
		}
	}

	sort.Slice(leaderboard, func(i, j int) bool {
		return leaderboard[i].ClosedPnl > leaderboard[j].ClosedPnl
	})

	return leaderboard, nil
}

func (api *LeaderboardAPI) GetLeaderboardThisMonth() ([]LeaderboardUser, error) {
	users := []User{}
	err := api.db.Find(&users).Error
	if err != nil {
		return nil, err
	}

	leaderboard := make([]LeaderboardUser, len(users))
	for i, user := range users {
		orders, err := api.getOrdersThisMonth(user.Email)
		if err != nil {
			return nil, err
		}

		closedPnl := api.calculateCumulativePnl(orders)
		leaderboard[i] = LeaderboardUser{
			UserEmail: user.Email,
			Name:      user.Name,
			ClosedPnl: closedPnl,
			Country:   user.Country,
			Username:  user.UserName,
		}
	}

	sort.Slice(leaderboard, func(i, j int) bool {
		return leaderboard[i].ClosedPnl > leaderboard[j].ClosedPnl
	})

	return leaderboard, nil
}

func (api *LeaderboardAPI) getOrdersThisWeek(email string) ([]Statements, error) {
	orders := []Statements{}

	now := time.Now()
	startOfWeek := now.AddDate(0, 0, -int(now.Weekday())).Truncate(24 * time.Hour)
	endOfWeek := startOfWeek.AddDate(0, 0, 7).Add(-time.Nanosecond)

	err := api.db.Model(Statements{}).Where("user_email = ? AND created_time >= ? AND created_time <= ?", email, startOfWeek, endOfWeek).Find(&orders).Error
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (api *LeaderboardAPI) getOrdersThisMonth(email string) ([]Statements, error) {
	orders := []Statements{}

	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	err := api.db.Model(Statements{}).Where("user_email = ? AND created_time >= ? AND created_time <= ?", email, startOfMonth, endOfMonth).Find(&orders).Error
	if err != nil {
		return nil, err
	}

	return orders, nil
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

func (api *LeaderboardAPI) GetLeaderboardAllTime() ([]LeaderboardUser, error) {
	users := []User{}
	err := api.db.Find(&users).Error
	if err != nil {
		return nil, err
	}

	leaderboard := make([]LeaderboardUser, len(users))
	for i, user := range users {
		orders, err := api.getAllTimeOrders(user.Email)
		if err != nil {
			return nil, err
		}

		closedPnl := api.calculateCumulativePnl(orders)
		leaderboard[i] = LeaderboardUser{
			UserEmail: user.Email,
			Name:      user.Name,
			ClosedPnl: closedPnl,
			Country:   user.Country,
			Username:  user.UserName,
		}
	}

	sort.Slice(leaderboard, func(i, j int) bool {
		return leaderboard[i].ClosedPnl > leaderboard[j].ClosedPnl
	})

	return leaderboard, nil
}

func (api *LeaderboardAPI) getAllTimeOrders(email string) ([]Statements, error) {
	orders := []Statements{}

	err := api.db.Model(Statements{}).Where("user_email = ?", email).Find(&orders).Error
	if err != nil {
		return nil, err
	}

	return orders, nil
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

// func (api *LeaderboardAPI) GetLeaderboardThisWeek() ([]LeaderboardUser, error) {
// 	users := []User{}
// 	err := api.db.Find(&users).Error
// 	if err != nil {
// 		return nil, err
// 	}

// 	leaderboard := make([]LeaderboardUser, len(users))
// 	for i, user := range users {
// 		orders, err := api.getOrdersThisWeek(user.Email)
// 		if err != nil {
// 			return nil, err
// 		}

// 		closedPnl := api.calculateCumulativePnl(orders)
// 		leaderboard[i] = LeaderboardUser{
// 			UserEmail: user.Email,
// 			Name:      user.Name,
// 			ClosedPnl: closedPnl,
// 		}
// 	}

// 	sort.Slice(leaderboard, func(i, j int) bool {
// 		return leaderboard[i].ClosedPnl > leaderboard[j].ClosedPnl
// 	})

// 	return leaderboard, nil
// }

// func (api *LeaderboardAPI) getOrdersThisWeek(email string) ([]Statements, error) {
// 	orders := []Statements{}

// 	now := time.Now()
// 	startOfWeek := now.AddDate(0, 0, -int(now.Weekday())).Truncate(24 * time.Hour)
// 	endOfWeek := startOfWeek.AddDate(0, 0, 7).Add(-time.Nanosecond)

// 	err := api.db.Model(Statements{}).Where("user_email = ? AND created_time >= ? AND created_time <= ?", email, startOfWeek, endOfWeek).Find(&orders).Error
// 	if err != nil {
// 		return nil, err
// 	}

// 	return orders, nil
// }

type Result struct {
	Symbol string  `json:"symbol"`
	Profit float64 `json:"profit"`
}

func (st *Statements) GetCoinwiseToday(db *gorm.DB, email string, service string) (*[]Result, error) {
	ords := []Result{}

	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.UTC)

	err := db.Model(Statements{}).Select("SUM(closed_pnl::float) AS profit, symbol").Where("user_email = ? AND exchange = ? AND created_time >= ? AND created_time <= ?", email, service, startOfDay, endOfDay).Group("symbol").Scan(&ords).Error

	if err != nil {
		return &[]Result{}, err
	}

	return &ords, nil
}

func (st *Statements) GetCoinwiseMonth(db *gorm.DB, email string, service string) (*[]Result, error) {
	ords := []Result{}

	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)

	err := db.Model(Statements{}).Select("SUM(closed_pnl::float) AS profit, symbol").Where("user_email = ? AND exchange = ? AND created_time >= ?", email, service, startOfMonth).Group("symbol").Scan(&ords).Error

	if err != nil {
		return &[]Result{}, err
	}

	return &ords, nil
}

func (st *Statements) GetCoinwiseAllTime(db *gorm.DB, email string, service string) (*[]Result, error) {
	ords := []Result{}

	err := db.Model(Statements{}).Select("SUM(closed_pnl::float) AS profit, symbol").Where("user_email = ? AND exchange = ?", email, service).Group("symbol").Scan(&ords).Error

	if err != nil {
		return &[]Result{}, err
	}

	return &ords, nil
}
