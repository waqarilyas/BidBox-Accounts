package models

import (
	"errors"

	"github.com/jinzhu/gorm"
)

type OrderRequest struct {
	Symbol           string `json:"symbol"`
	MarginCoin       string `json:"marginCoin"`
	Size             string `json:"size"`
	Side             string `json:"side"`
	OrderType        string `json:"orderType"`
	TimeInForceValue string `json:"timeInForceValue"`
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
	Size       string
	Side       string
	OrderType  string
	ClientID   string
	OrderID    string
}

func (o *Order) Initialize(order OrderRequest, email string, client_id string, order_id string) {
	o.ClientID = client_id
	o.OrderID = order_id
	o.MarginCoin = order.MarginCoin
	o.Side = order.Side
	o.Symbol = order.Symbol
	o.Size = order.Size
	o.OrderType = order.OrderType
	o.Email = email
}

func (o *Order) Validate() error {
	if o.Email == "" {
		errors.New("email is required")
	}
	if o.MarginCoin == "" {
		errors.New("margin coin is required")
	}
	if o.OrderType == "" {
		errors.New("ordertype is required")
	}
	if o.Side == "" {
		errors.New("side is required")
	}
	if o.Size == "" {
		errors.New("size is required")
	}
	if o.Symbol == "" {
		errors.New("symbol is required")
	}
	if o.OrderID == "" {
		errors.New("bitget error, could not get order ID")
	}
	if o.ClientID == "" {
		errors.New("bitget error, could not get client ID")
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
