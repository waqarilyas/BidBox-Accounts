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

type OrderWithPosition struct {
	Id              uuid.UUID `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	Email           string    `json:"email"`
	Symbol          string    `json:"symbol"`
	MarginCoin      string    `json:"margin_coin"`
	Service         string    `json:"service"`
	Size            string    `json:"size"`
	Side            string    `json:"side"`
	OrderType       string    `json:"order_type"`
	CreatedAt       time.Time `json:"created_at"`
	Profit          float64   `json:"profit"`
	PositionId      int       `json:"position_id"`
	OrderPrice      string    `json:"order_price"`
	QuoteAmount     string    `json:"quote_amount"`
	Fee             float64   `json:"fee"`
	OrderId         string    `json:"order_id"`
	UpdatedAt       time.Time `gorm:"type:timestamptz;default:now()" json:"updated_at"`
	Leverage        string    `json:"leverage"`
	OpenPrice       string    `json:"open_price"`
	LiqPrice        string    `json:"liq_price"`
	TakeProfit      string    `json:"take_profit"`
	MarkPrice       string    `json:"mark_price"`
	StopLoss        string    `json:"stop_loss"`
	UnrealizedPl    string    `json:"unrealized_pl"`
	Margin          string    `json:"margin"`
	UserEmail       string    `gorm:"not null" json:"user_email"`
	Status          string    `gorm:"default:'opened'" json:"status"`
	Exchange        string    `json:"exchange"`
	LastUpdatePrice string    `json:"last_update_price"`
	Layer           int       `json:"layer"`
	TotalProfit     float64   `json:"total_profit"`
	FirstBuyAmount  string    `json:"first_buy_amount"`
	HedgeId         string    `json:"hedge_id"`
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

// func GetOrdersWithPositionDetails(db *gorm.DB, email string, exchange string) ([]*OrderWithPosition, error) {
// 	var orders []*OrderWithPosition

// 	err := db.Where("email = ? AND service = ?", email, exchange).
// 		Order("orders.created_at DESC").
// 		Joins("JOIN positions ON orders.position_id = positions.id").
// 		Preload("Order").
// 		Preload("Position").
// 		Find(&orders).Error

// 	if err != nil {
// 		return nil, err
// 	}

// 	return orders, nil
// }

func GetOrdersWithPositionDetails(db *gorm.DB, email string, exchange string) ([]OrderWithPosition, error) {
	var ordersWithPosition []OrderWithPosition

	err := db.
		Select("orders.*, positions.*").
		Table("orders").
		Joins("JOIN positions ON orders.position_id = positions.id").
		Where("orders.email = ? AND orders.service = ? And (orders.side = ? OR orders.side = ?)", email, exchange, "close_long", "close_short").
		Order("orders.created_at DESC").
		Scan(&ordersWithPosition).
		Error

	if err != nil {
		return nil, err
	}

	return ordersWithPosition, nil
}
