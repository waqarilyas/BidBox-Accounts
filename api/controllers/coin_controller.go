package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/kryptomind/bidboxapi/AccountsService/api/models/admin"
	"github.com/kryptomind/bidboxapi/AccountsService/api/response"
)

// Save Coin godoc
// @Summary      Save Coin
// @Description  save coin
// @Tags         coins
// @Accept       json
// @Produce      json
// @Param        coin  body  string   true  "coin pair"
// @Success      200  {object}  admin.CoinPair
// @Failure      422  {string}  coin required
// @Failure      500  {string}  server error
// @Router       /admin/coins [post]
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

// GetCoins godoc
// @Summary      Get All Coins
// @Description  get coins
// @Tags         coins
// @Accept       json
// @Produce      json
// @Success      200  {object}  []admin.CoinPair
// @Failure      500  {string}  server error
// @Router       /admin/coins [get]
func (server *Server) GetCoins(w http.ResponseWriter, r *http.Request) {

	coin := admin.CoinPair{}

	coins, err := coin.GetAllCoins(server.DB)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	response.JSON(w, http.StatusOK, coins)
}

// Update Coin godoc
// @Summary      Update Coin
// @Description  update coin
// @Tags         coins
// @Accept       json
// @Produce      json
// @Param        coin  body  admin.CoinPair  false  "coin pair"
// @Success      200  {object}  admin.CoinPair
// @Failure      422  {string}  coin required
// @Failure      500  {string}  server error
// @Router       /admin/coins [put]
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

	prev, err := coin.GetCoinById(server.DB, coin.Coin)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	coin.Id = prev.Id
	coin.Validate(prev)
	updatedcoin, err := coin.UpdateCoinPair(server.DB, coin.Coin)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	response.JSON(w, http.StatusOK, updatedcoin)
}

// Delete Coin godoc
// @Summary      Delete Coin
// @Description  delete coin
// @Tags         coins
// @Accept       json
// @Produce      json
// @Param        coin  query  string true "coin pair"
// @Success      200  {string}  deleted row
// @Failure      400  {string}  coin required
// @Failure      500  {string}  server error
// @Router       /admin/coins [delete]
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
