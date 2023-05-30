package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/kryptomind/bidboxapi/AccountsService/api/models"
	"github.com/kryptomind/bidboxapi/AccountsService/api/response"
)

func (server *Server) GetSettings(w http.ResponseWriter, r *http.Request, email string) {
	key := models.Key{}
	keys, err := key.GetSettings(server.DB, email)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	res := make(map[string]interface{})
	for _, v := range *keys {
		res["strategy"] = v.Strategy
		res["mode"] = v.Mode
		res["auto_compound"] = v.Compound
	}

	response.JSON(w, http.StatusOK, res)
}

func (server *Server) GetHistory(w http.ResponseWriter, r *http.Request, email string) {
	service := r.URL.Query().Get("service")

	if service == "" {
		response.ERROR(w, http.StatusBadRequest, errors.New("service is required"))
		return
	}

	hist := models.History{}
	history, err := hist.GetHistory(server.DB, email, service)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, history)
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

	if err := key.Validate(prev); err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

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
