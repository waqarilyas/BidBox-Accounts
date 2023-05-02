package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/kryptomind/bidboxapi/AccountsService/api/models"
	"github.com/kryptomind/bidboxapi/AccountsService/api/response"
)

var strategy = []string{"Cycle", "Single", "Stop Make", "Stop Long", "Stop Short"}
var modes = []string{"Conservative", "Aggressive"}

type stratBody struct {
	Strategy string `json:"strategy"`
}

type modeBody struct {
	Mode string `json:"mode"`
}

func (server *Server) UpdateStrategy(w http.ResponseWriter, r *http.Request, email string) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := models.User{}
	user.Email = email

	strat := stratBody{}
	err = json.Unmarshal(body, &strat)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	found := false
	for _, v := range strategy {
		if v == strat.Strategy {
			found = true
			break
		}
	}
	if found == false {
		response.ERROR(w, http.StatusUnprocessableEntity, errors.New("no such strategy exists"))
		return
	}
	updatedUser, err := user.ChangeStrategy(server.DB, strat.Strategy)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	response.JSON(w, http.StatusOK, updatedUser)
}

func (server *Server) UpdateMode(w http.ResponseWriter, r *http.Request, email string) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := models.User{}
	user.Email = email

	mode := modeBody{}
	err = json.Unmarshal(body, &mode)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	found := false
	for _, v := range modes {
		if v == mode.Mode {
			found = true
			break
		}
	}
	if found == false {
		response.ERROR(w, http.StatusUnprocessableEntity, errors.New("no such mode exists"))
		return
	}
	updatedUser, err := user.ChangeMode(server.DB, mode.Mode)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	response.JSON(w, http.StatusOK, updatedUser)
}
