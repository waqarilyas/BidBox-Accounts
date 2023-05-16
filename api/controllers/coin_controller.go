package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/kryptomind/bidboxapi/AccountsService/api/models/admin"
	"github.com/kryptomind/bidboxapi/AccountsService/api/response"
)

func (server *Server) CreateCoin(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	coin := admin.CoinPair{}
	err = json.Unmarshal(body, &coin)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	if coin.Coin == "" {
		response.ERROR(w, http.StatusUnprocessableEntity, errors.New("coin required"))
		return
	}

	coinCreated, err := coin.SaveCoinPair(server.DB)

	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusCreated, coinCreated)
}

func (server *Server) GetCoins(w http.ResponseWriter, r *http.Request) {

	coin := admin.CoinPair{}

	coins, err := coin.GetAllCoins(server.DB)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	response.JSON(w, http.StatusOK, coins)
}

func (server *Server) UpdateCoin(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	coin := admin.CoinPair{}
	err = json.Unmarshal(body, &coin)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	if coin.Coin == "" {
		response.ERROR(w, http.StatusUnprocessableEntity, errors.New("coin pair is required"))
		return
	}

	updatedcoin, err := coin.UpdateCoinPair(server.DB, coin.Coin)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	response.JSON(w, http.StatusOK, updatedcoin)
}

func (server *Server) DeleteCoin(w http.ResponseWriter, r *http.Request) {

	del := r.URL.Query().Get("coin")
	if del == "" {
		response.ERROR(w, http.StatusBadRequest, errors.New("coin required"))
		return
	}
	coin := admin.CoinPair{}

	_, err := coin.DeleteCoinPair(server.DB, del)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	response.JSON(w, http.StatusOK, "deleted row")
}
