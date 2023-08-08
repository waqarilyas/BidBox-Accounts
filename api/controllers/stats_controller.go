package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/kryptomind/bidboxapi/AccountsService/api/models"
	"github.com/kryptomind/bidboxapi/AccountsService/api/response"
)

type StatsResponse struct {
	OrdersProfit               float64 `json:"ordersProfit"`
	OrdersProfitCount          int     `json:"ordersProfitCount"`
	TodayOrdersProfit          float64 `json:"todayOrdersProfit"`
	TodayOrdersProfitCount     int     `json:"todayOrdersProfitCount"`
	StatementsProfit           float64 `json:"statementsProfit"`
	StatementsProfitCount      int     `json:"statementsProfitCount"`
	TodayStatementsProfit      float64 `json:"todayStatementsProfit"`
	TodayStatementsProfitCount int     `json:"todayStatementsProfitCount"`
}

func (s *Server) GetAccountStats(w http.ResponseWriter, r *http.Request, email string) {
	service := r.URL.Query().Get("service")
	if service == "" {
		response.ERROR(w, http.StatusBadRequest, errors.New("service required"))
		return
	}

	if service != "bitget" && service != "binance" && service != "bybit" {
		response.ERROR(w, http.StatusBadRequest, errors.New("service not supported"))
		return
	}

	sumProfit, profitOrdersCount, err := models.SumProfitByEmailAndService(s.DB, email, service)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, errors.New("unable to get stats at the moment"))
		return
	}

	todayProfitSum, todayTotalProfitOrders, err := models.SumAndCountProfitForToday(s.DB, email, service)
	if err != nil {
		fmt.Println("🚀 ~ file: stats_controller.go:40 ~ func ~ err:", err)
		response.ERROR(w, http.StatusBadRequest, errors.New("unable to get stats at the moment"))
		return
	}

	statementSum, totalProfitStatement, err := models.SumStatementByEmailAndService(s.DB, email, service)
	if err != nil {
		fmt.Println("🚀 ~ file: stats_controller.go:40 ~ func ~ err:", err)
		response.ERROR(w, http.StatusBadRequest, errors.New("unable to get stats at the moment"))
		return
	}

	todayStatementsSum, todayTotalStatementsOrders, err := models.SumAndCountStatementForToday(s.DB, email, service)
	if err != nil {
		fmt.Println("🚀 ~ file: stats_controller.go:40 ~ func ~ err:", err)
		response.ERROR(w, http.StatusBadRequest, errors.New("unable to get stats at the moment"))
		return
	}

	payload := StatsResponse{
		OrdersProfit:               sumProfit,
		OrdersProfitCount:          profitOrdersCount,
		TodayOrdersProfit:          todayProfitSum,
		TodayOrdersProfitCount:     todayTotalProfitOrders,
		TodayStatementsProfit:      todayStatementsSum,
		TodayStatementsProfitCount: todayTotalStatementsOrders,
		StatementsProfit:           statementSum,
		StatementsProfitCount:      totalProfitStatement,
	}

	response.JSON(w, http.StatusOK, payload)

}
