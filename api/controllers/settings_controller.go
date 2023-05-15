package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

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

	conds, err := cond.GetConditions(server.DB, limit, offset)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	response.JSON(w, http.StatusOK, conds)
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
