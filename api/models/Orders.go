package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
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
	Id          uuid.UUID `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	Email       string    `json:"email"`
	Symbol      string    `json:"symbol"`
	MarginCoin  string    `json:"margin_coin"`
	Service     string    `json:"service"`
	Size        string    `json:"size"`
	Side        string    `json:"side"`
	OrderType   string    `json:"order_type"`
	CreatedAt   time.Time `json:"created_at"`
	Profit      float64   `json:"profit"`
	PositionId  int       `json:"position_id"`
	OrderPrice  string    `json:"order_price"`
	QuoteAmount string    `json:"quote_amount"`
	Fee         float64   `json:"fee"`
	OrderId     string    `json:"order_id"`
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

func GetOrdersByUserEmailAndExchange(db *gorm.DB, email string, exchange string) ([]*Order, error) {
	var dbOrders []*Order
	err := db.Where("email = ? AND service = ?", email, exchange).Order("created_at DESC").Find(&dbOrders).Error
	if err != nil {
		return nil, err
	}
	return dbOrders, nil
}
