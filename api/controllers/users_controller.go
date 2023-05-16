package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
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
	if !found {
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

func (server *Server) UpdateKeySettings(w http.ResponseWriter, r *http.Request, email string) {
	service := r.URL.Query().Get("service")

	if service == "" {
		response.ERROR(w, http.StatusBadRequest, errors.New("service is required"))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	key := models.Key{}
	err = json.Unmarshal(body, &key)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	prev, err := key.FindKeysByEmail(server.DB, email, service)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	log.Println(key.Compound)
	log.Println(prev.Compound)
	key.Validate(prev)
	log.Println(key.Compound)
	updatedkey, err := key.UpdateKeySettings(server.DB, email, service)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	response.JSON(w, http.StatusOK, updatedkey)
}

func (s *Server) DisconnectKey(w http.ResponseWriter, r *http.Request, email string) {
	service := r.URL.Query().Get("service")

	if service == "" {
		response.ERROR(w, http.StatusBadRequest, errors.New("service is required"))
		return
	}

	key := models.Key{}
	_, err := key.DeleteKey(s.DB, email, service)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, "disconnected key")
}
