package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

type OrderRequest struct {
	Symbol     string `json:"symbol"`
	MarginCoin string `json:"marginCoin"`
	Size       string `json:"size"`
	Side       string `json:"side"`
	OrderType  string `json:"orderType"`
}

type CancelOrderRequest struct {
	Symbol     string `json:"symbol"`
	MarginCoin string `json:"marginCoin"`
	OrderId    string `json:"orderId"`
}

type OrderResponse struct {
	Code        string `json:"code"`
	Msg         string `json:"msg"`
	RequestTime int64  `json:"requestTime"`
	Data        struct {
		ClientOid string `json:"clientOid"`
		OrderID   string `json:"orderId"`
	} `json:"data"`
}

type Order struct {
	Email      string
	Symbol     string
	MarginCoin string
	Service    string
	Size       string
	Side       string
	OrderType  string
	CreatedAt  time.Time
	Profit     float64
}

func (o *Order) Initialize(order OrderRequest, email string, client_id string, order_id string) {
	o.MarginCoin = order.MarginCoin
	o.Side = order.Side
	o.Symbol = order.Symbol
	o.Size = order.Size
	o.OrderType = order.OrderType
	o.Email = email
}

func (o *OrderRequest) Validate() error {
	if o.MarginCoin == "" {
		return errors.New("margin coin is required")
	}
	if o.OrderType == "" {
		return errors.New("ordertype is required")
	}
	if o.Side == "" {
		return errors.New("side is required")
	}
	if o.Size == "" {
		return errors.New("size is required")
	}
	if o.Symbol == "" {
		return errors.New("symbol is required")
	}
	return nil
}

func (o *Order) SaveOrder(db *gorm.DB) (*Order, error) {
	err := db.Debug().Create(&o).Error
	if err != nil {
		return &Order{}, err
	}
	return o, nil
}

func (o *Order) GetActiveTrades(db *gorm.DB) (int, error) {
	var count int

	err := db.Model(Order{}).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (o *Order) GetOrderThisMonth(db *gorm.DB) (*[]Order, error) {
	orders := []Order{}

	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)

	err := db.Model(Order{}).Order("profit DESC").Where("created_at >= ?", startOfMonth).Find(&orders).Error
	if err != nil {
		return &[]Order{}, err
	}

	return &orders, nil
}

func (o *Order) GetOrderThisDay(db *gorm.DB) (*[]Order, error) {
	ords := []Order{}

	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.UTC)

	err := db.Model(Order{}).Order("profit DESC").Where("created_at >= ? AND created_at <= ?", startOfDay, endOfDay).Find(&ords).Error
	if err != nil {
		return &[]Order{}, err
	}

	return &ords, nil
}

func (o *Order) GetOrderAllTime(db *gorm.DB) (*[]Order, error) {
	ords := []Order{}

	err := db.Model(Order{}).Order("profit DESC").Find(&ords).Error
	if err != nil {
		return &[]Order{}, err
	}

	return &ords, nil
}
