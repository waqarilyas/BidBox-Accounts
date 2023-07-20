package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/kryptomind/bidboxapi/AccountsService/api/models"
	"github.com/kryptomind/bidboxapi/AccountsService/api/models/admin"
	"github.com/kryptomind/bidboxapi/AccountsService/api/response"
)

func (server *Server) GetConditions(w http.ResponseWriter, r *http.Request) {

	// Parse query parameters
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 10
	}

	// Calculate offset based on page and limit
	offset := (page - 1) * limit

	cond := admin.Conditions{}

	conds, count, err := cond.GetConditions(server.DB, limit, offset)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	resp := make(map[string]interface{})

	resp["total"] = count
	resp["conditions"] = conds

	response.JSON(w, http.StatusOK, resp)
}

func (server *Server) UpdateConditions(w http.ResponseWriter, r *http.Request) {

	capital, err := strconv.Atoi(r.URL.Query().Get("capital"))
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, errors.New("capital required"))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	cond := admin.Conditions{}
	err = json.Unmarshal(body, &cond)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	cond.Capital = capital
	prev, err := cond.FindConditionById(server.DB, capital)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	cond.Validate(prev)
	updatedcond, err := cond.UpdateConditions(server.DB, capital)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	response.JSON(w, http.StatusOK, updatedcond)
}

func (server *Server) changeSettings(w http.ResponseWriter, r *http.Request) {
	res := make(map[string]interface{}, 0)

	s := admin.Settings{}
	curr, err := s.GetSettings(server.DB)

	if err != nil {
		res["status"] = "error"
		res["message"] = err.Error()
		res["data"] = nil

		response.JSON(w, http.StatusBadRequest, res)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		res["status"] = "error"
		res["message"] = err.Error()
		res["data"] = nil

		response.JSON(w, http.StatusBadRequest, res)
		return
	}

	if s.Layers != 0 {
		_, err := s.UpdateLayers(server.DB, s.Layers)
		if err != nil {
			res["status"] = "error"
			res["message"] = err.Error()
			res["data"] = nil

			response.JSON(w, http.StatusInternalServerError, err)
			return
		}
		curr.Layers = s.Layers
	}

	if s.ProfitPercentage != 0 {
		_, err := s.UpdateProfitPercentage(server.DB, s.ProfitPercentage)
		if err != nil {
			res["status"] = "error"
			res["message"] = err.Error()
			res["data"] = nil

			response.JSON(w, http.StatusInternalServerError, err)
			return
		}
		curr.ProfitPercentage = s.ProfitPercentage
	}

	if s.Leverage != 0 {
		_, err := s.UpdateLeverage(server.DB, s.Leverage)
		if err != nil {
			res["status"] = "error"
			res["message"] = err.Error()
			res["data"] = nil
			response.JSON(w, http.StatusInternalServerError, err)
			return
		}
		curr.Leverage = s.Leverage
	}

	res["status"] = "ok"
	res["message"] = "success"
	res["data"] = curr
	response.JSON(w, http.StatusOK, res)
}

func (server *Server) getSettings(w http.ResponseWriter, r *http.Request) {

	s := admin.Settings{}

	setting, err := s.GetSettings(server.DB)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, errors.New("error getting settings"))
		return
	}

	response.JSON(w, http.StatusOK, setting)
}

// check this
func (server *Server) GetLeaderBoard(w http.ResponseWriter, r *http.Request) {
	time := r.URL.Query().Get("time")
	var ords []models.Statements
	if time == "month" {
		order := models.Statements{}
		orders, err := order.GetOrderThisMonth(server.DB)
		if err != nil {
			response.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		ords = *orders
	} else if time == "day" {
		order := models.Statements{}
		orders, err := order.GetOrderThisDay(server.DB)
		if err != nil {
			response.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		ords = *orders

	} else if time == "week" {
		order := models.Statements{}
		orders, err := order.GetOrderThisDay(server.DB)
		if err != nil {
			response.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		ords = *orders
	} else if time == "alltime" || time == "" {
		order := models.Statements{}
		orders, err := order.GetOrderAllTime(server.DB)
		if err != nil {
			response.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		ords = *orders
	} else {
		response.ERROR(w, http.StatusBadRequest, errors.New("timeframe param is incorrect"))
		return
	}
	response.JSON(w, http.StatusOK, ords)
}

func (server *Server) GetLeaderBoardv2(w http.ResponseWriter, r *http.Request) {
	time := r.URL.Query().Get("time")
	var ords []models.LeaderboardUser
	if time == "month" {
		api := models.NewLeaderboardAPI(server.DB)
		data, error := api.GetLeaderboardThisMonth()
		if error != nil {
			fmt.Println("------- error handling ------", error)
			response.ERROR(w, http.StatusBadRequest, errors.New("something went wrong"))
			return
		}
		ords = data

	} else if time == "day" {
		api := models.NewLeaderboardAPI(server.DB)
		data, error := api.GetLeaderboardToday()
		if error != nil {

			fmt.Println("------- error handling ------", error)
		}
		ords = data

	} else if time == "week" {
		api := models.NewLeaderboardAPI(server.DB)
		data, error := api.GetLeaderboardThisWeek()
		if error != nil {
			response.ERROR(w, http.StatusBadRequest, errors.New("something went wrong"))
			fmt.Println("------- error handling ------", error)
			return
		}
		ords = data

	} else if time == "alltime" || time == "" {
		api := models.NewLeaderboardAPI(server.DB)
		data, error := api.GetLeaderboardAllTime()
		if error != nil {
			fmt.Println("------- error handling ------", error)
			response.ERROR(w, http.StatusBadRequest, errors.New("something went wrong"))
			return
		}
		ords = data

	} else {
		response.ERROR(w, http.StatusBadRequest, errors.New("timeframe param is incorrect"))
		return
	}
	response.JSON(w, http.StatusOK, ords)
}
